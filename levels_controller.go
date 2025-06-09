package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto/v2/z"
	otrace "go.opencensus.io/trace"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"wiscdb/pb"
	"wiscdb/table"
	"wiscdb/y"
)

type LevelsController struct {
	nextFileID    atomic.Uint64
	l0stallMs     atomic.Int64
	levels        []*levelHandler
	kv            *DB
	compactStatus compactStatus
}

func newLevelController(db *DB, mf *Manifest) (*LevelsController, error) {
	return &LevelsController{}, nil
}

// return file ID
func (s *LevelsController) ReserveFileID() uint64 {
	id := s.nextFileID.Add(1)
	return id - 1
}

func (s *LevelsController) close() error {
	return nil
}

func (s *LevelsController) cleanupLevels() error {
	var firstErr error
	for _, l := range s.levels {
		if err := l.close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (s *LevelsController) AddLevel0Table(t *table.Table) error {
	if !t.IsInMemory {
		err := s.kv.Manifest.AddChanges([]*pb.ManifestChange{
			NewCreateChange(t.ID(), 0, t.KeyID(), t.CompressionType()),
		})
		if err != nil {
			return err
		}
	}

	// 这是for  循环,如果 table数量超过15则一直会循环 然后retry
	for !s.levels[0].tryAddLevel0Table(t) {
		timeStart := time.Now()
		// 如果当前Level0的table数量超过 L0 的15个stall table的数量则休息10毫秒
		for s.levels[0].numTables() >= s.kv.Opt.NumLevelZeroTablesStall {
			time.Sleep(10 * time.Millisecond)
		}
		dur := time.Since(timeStart)
		if dur > time.Second {
			s.kv.Opt.Infof("L0 was stalled for %s\n", dur.Round(time.Millisecond))
		}
		s.l0stallMs.Add(int64(dur.Round(time.Millisecond)))
	}
	return nil
}

func (s *LevelsController) startCompact(lc *z.Closer) {
	n := s.kv.Opt.NumCompactors
	lc.AddRunning(n - 1)
	//N指的是执行压缩协程的数量.
	for i := 0; i < n; i++ {
		go s.runCompactor(i, lc)
	}
}

// 执行压缩操作.
func (s *LevelsController) runCompactor(id int, lc *z.Closer) {
	defer lc.Done()
	//创建一个计时器,随机设置计时器的时间为1秒内
	randomDelay := time.NewTimer(time.Duration(rand.Int31n(1000)) * time.Millisecond)
	select {
	//如果时间倒后或者完成了那么久停止计时器同时返回函数
	case <-randomDelay.C:
	case <-lc.HasBeenClosed():
		randomDelay.Stop()
		return
	}
	//输入多个level的优先级，返回的还是多个level的优先级
	//将L0移动到前方
	moveL0toFront := func(prios []compactionPriority) []compactionPriority {
		idx := -1
		for i, p := range prios {
			//找到Level0 所在的索引位置.
			if p.level == 0 {
				idx = i
				break
			}
		}
		//如果idx大于0表示L0不在最前方，那么需要重新拼接prios
		//使用空值append level0，然后不改变其他相对顺序.
		if idx > 0 {
			out := append([]compactionPriority{}, prios[idx])
			out = append(out, prios[:idx]...)
			out = append(out, prios[idx+1:]...)
			return out
		}
		return prios
	}

	//执行压缩操作
	run := func(p compactionPriority) bool {
		//按照id和压缩优先级进行压缩
		err := s.doCompact(id, p)
		switch err {
		case nil:
			return true
		case errFillTables:
		default:
			s.kv.Opt.Warningf("While runing doCompact: %v\n", err)
		}
		return false
	}

	var priosBuffer []compactionPriority
	//用压缩优先级数组来存储每一层的层高,分数等信息
	runOnce := func() bool {
		prios := s.pickCompactLevels(priosBuffer)
		defer func() {
			priosBuffer = prios
		}()

		//id为0表示 执行压缩操作的第一个协程
		if id == 0 {
			prios = moveL0toFront(prios)
		}
		//runOnce会优先将L0放在最上方值i相关压缩
		for _, p := range prios {
			if id == 0 && p.level == 0 {
				//	对于Level0 什么事情都不做
			} else if p.adjusted < 1.0 {
				break
			}
			//执行具体的压缩工作
			if run(p) {
				return true
			}
		}
		return false
	}

	//只压缩最高层Ln
	tryLmaxToLmaxCompaction := func() {
		p := compactionPriority{
			level: s.lastLevel().level,
			//计算每一层的目标尺寸和file尺寸，按照10倍递增
			t: s.levelTargets(),
		}
		run(p)
	}
	count := 0
	//设置计时器50毫秒,每50毫秒返回一次ticker.C
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	//无限循环
	for {
		select {
		case <-ticker.C:
			count++
			//如果count超过200表示压缩时间超过了10秒钟
			if s.kv.Opt.LMaxCompaction && id == 2 && count >= 200 {
				//调用tryLmaxToLmaxCompaction
				tryLmaxToLmaxCompaction()
				count = 0
			} else {
				//如果没有超过10秒则调用runOnce函数
				runOnce()
			}
		case <-lc.HasBeenClosed():
			return
		}
	}
}

var errFillTables = errors.New("Unable to fill tables")

// pickCompactLevel determines which level to compact.
// Based on: https://github.com/facebook/rocksdb/wiki/Leveled-Compaction
// 计算优先级
func (s *LevelsController) pickCompactLevels(priosBuffer []compactionPriority) (prios []compactionPriority) {
	t := s.levelTargets()
	addPriority := func(level int, score float64) {
		pri := compactionPriority{
			level:    level,
			score:    score,
			adjusted: score,
			t:        t,
		}
		prios = append(prios, pri)
	}
	if cap(priosBuffer) < len(s.levels) {
		priosBuffer = make([]compactionPriority, 0, len(s.levels))
	}
	prios = priosBuffer[:0]
	addPriority(0, float64(s.levels[0].numTables())/float64(s.kv.Opt.NumLevelZeroTables))

	for i := 1; i < len(s.levels); i++ {
		delSize := s.compactStatus.delSize(i)
		l := s.levels[i]
		sz := l.getTotalSize() - delSize
		addPriority(i, float64(sz)/float64(t.targetSz[i]))
	}
	y.AssertTrue(len(prios) == len(s.levels))

	var preLevel int
	for level := t.baseLevel; level < len(s.levels); level++ {
		if prios[preLevel].adjusted >= 1 {
			const minScore = 0.01
			if prios[level].score >= minScore {
				prios[preLevel].adjusted /= prios[level].adjusted
			} else {
				prios[preLevel].adjusted /= minScore
			}
		}
		preLevel = level
	}

	out := prios[:0]
	for _, p := range prios[:len(prios)-1] {
		if p.score >= 1.0 {
			out = append(out, p)
		}
	}
	prios = out
	sort.Slice(prios, func(i, j int) bool {
		return prios[i].adjusted > prios[j].adjusted
	})
	return prios
}

// 返回每一层的目标尺寸以及文件尺寸
// 层与层之间的比例是10倍
func (s *LevelsController) levelTargets() targets {
	//如果当前DB的大小小于默认的10MB，则返回10MB，如果超过10MB则返回真正的DB的大小
	adjust := func(sz int64) int64 {
		if sz < s.kv.Opt.BaseLevelSize {
			return s.kv.Opt.BaseLevelSize
		}
		return sz
	}

	//初始化target 目标尺寸和文件尺寸
	t := targets{
		targetSz: make([]int64, len(s.levels)),
		fileSz:   make([]int64, len(s.levels)),
	}
	// 获取LN层的文件大小
	dbSize := s.lastLevel().getTotalSize()
	//遍历N个层
	for i := len(s.levels) - 1; i > 0; i-- {
		lTarget := adjust(dbSize)
		//将调整后的值写到目标值
		t.targetSz[i] = lTarget
		//如果当前base层是0层并且目标值小于等于Base值则将baseLevel设置成i层
		if t.baseLevel == 0 && lTarget <= s.kv.Opt.BaseLevelSize {
			t.baseLevel = i
		}
		//更新当前DBSize为DBSize处理
		//设置上一层的DB尺寸，按照10倍递减，
		//比如层数是8，设置的尺寸是8000GB，
		//第7层的目标尺寸是800GB
		//第6层的目标尺寸是80GB
		//第5层的目标尺寸是8GB
		//第4层的目标尺寸是800MB
		//第三层的目标尺寸是80MB
		//第二层的目标尺寸是8MB
		//第一层的目标尺寸是800KB
		dbSize /= int64(s.kv.Opt.LevelSizeMultiplier)
	}
	//L1 的大小是10MB
	tsz := s.kv.Opt.BaseLevelSize
	//从L0到LN-1遍历
	for i := 0; i < len(s.levels); i++ {
		if i == 0 {
			//如果是L0,则将memTableSize写到fileSz中
			t.fileSz[i] = s.kv.Opt.MemTableSize
		} else if i <= t.baseLevel {
			//
			t.fileSz[i] = tsz
		} else {
			tsz *= int64(s.kv.Opt.TableSizeMultiplier)
			//fileSz时10MB的10的N次方
			t.fileSz[i] = tsz
		}
	}

	for i := t.baseLevel + 1; i < len(s.levels)-1; i++ {
		if s.levels[i].getTotalSize() > 0 {
			break
		}
		t.baseLevel = i
	}
	b := t.baseLevel
	lvl := s.levels
	if b < len(lvl)-1 && lvl[b].getTotalSize() == 0 && lvl[b+1].getTotalSize() < t.targetSz[b+1] {
		t.baseLevel++
	}
	return t
}

// 返回最下面的一层.
func (s *LevelsController) lastLevel() *levelHandler {
	return s.levels[len(s.levels)-1]
}

// 按照协程ID以及压缩优先级进行压缩
func (s *LevelsController) doCompact(id int, p compactionPriority) error {
	l := p.level
	y.AssertTrue(l < s.kv.Opt.MaxLevels)
	//如果是L0
	//我理解的baseLevel指的是当前DB的最高层是L0即所有数据还存在
	//memtable中，第一次准备写到磁盘中
	//则定义每一层的约定
	if p.t.baseLevel == 0 {
		p.t = s.levelTargets()
	}

	//创建span用来记录trace
	_, span := otrace.StartSpan(context.Background(), "WISC.Compaction")
	//添加defer 函数用来关闭span
	defer span.End()
	cd := compactDef{
		compactorId:  id,
		span:         span,
		p:            p,
		t:            p.t,
		thisLevel:    s.levels[l],
		dropPrefixes: p.dropPrefixes,
	}

	//如果是L0 表示使用memtableflush到disk中
	if l == 0 {
		cd.nextLevel = s.levels[p.t.baseLevel]
		if !s.fillTablesL0(&cd) {
			return errFillTables
		}
	} else {
		cd.nextLevel = cd.thisLevel
		if !cd.thisLevel.isLastLevel() {
			cd.nextLevel = s.levels[l+1]
		}
		if !s.fillTables(&cd) {
			return errFillTables
		}
	}
	defer s.compactStatus.delete(cd)
	span.Annotatef(nil, "Compaction: %+v", cd)

	if err := s.runCompactDef(id, l, cd); err != nil {
		s.kv.Opt.Warningf("[Compactor: %d] LOG Compact FAILED with error: %+v: %+v", id, err, cd)
		return err
	}
	s.kv.Opt.Debugf("[Compactor: %d] Compaction for level: %d DONE", id, cd.thisLevel.level)
	return nil
}

func (s *LevelsController) runCompactDef(id, l int, cd compactDef) (err error) {
	return nil
}

func tableToString(tables []*table.Table) []string {
	return nil
}

func buildChangeSet(cd *compactDef, newTables []*table.Table) pb.ManifestChangeSet {
	return pb.ManifestChangeSet{}
}

func (s *LevelsController) compactBuildTables(lev int, cd compactDef) ([]*table.Table, func() error, error) {
	return nil, func() error {
		return nil
	}, nil
}

func (s *LevelsController) subCompact(it y.Iterator, kr keyRange, cd compactDef, inflightBuilders *y.Throttle, res chan<- *table.Table) {

}

func hasAnyPrefixes(s []byte, listOfPrefixes [][]byte) bool {
	return false
}

func (s *LevelsController) checkOverlap(tables []*table.Table, lev int) bool {
	return false
}
func (s *LevelsController) addSplits(cd *compactDef) {

}
func (s *LevelsController) fillTables(cd *compactDef) bool {
	return false
}

func (s *LevelsController) sortByHeuristic(tables []*table.Table, cd *compactDef) {

}

func (s *LevelsController) fillMaxLevelTables(tables []*table.Table, cd *compactDef) bool {

	return false
}

func (s *LevelsController) sortByStaleDataSize(tables []*table.Table, cd *compactDef) {

}

// 压缩L0层的数据
func (s *LevelsController) fillTablesL0(cd *compactDef) bool {
	if ok := s.fillTablesL0ToLBase(cd); ok {
		return true
	}
	return s.fillTablesL0ToL0(cd)
}
func (s *LevelsController) fillTablesL0ToL0(cd *compactDef) bool {
	if cd.compactorId != 0 {
		return false
	}
	cd.nextLevel = s.levels[0]
	cd.nextRange = keyRange{}
	cd.bot = nil
	y.AssertTrue(cd.thisLevel.level == 0)
	y.AssertTrue(cd.nextLevel.level == 0)
	s.levels[0].RLock()
	defer s.levels[0].RUnlock()

	s.compactStatus.Lock()
	defer s.compactStatus.Unlock()

	// 获取该层的table
	top := cd.thisLevel.tables
	var out []*table.Table
	now := time.Now()
	for _, t := range top {
		// 如果table尺寸超过就ignore
		if t.Size() >= 2*cd.t.fileSz[0] {
			continue
		}
		// 如果超过10秒钟还没有被压缩则ignore
		if now.Sub(t.CreatedAt) < 10*time.Second {
			continue
		}
		// 如果table的状态已经是压缩过了,则ignore
		if _, beingCompacted := s.compactStatus.tables[t.ID()]; beingCompacted {
			continue
		}
		// 将该table添加到out中
		out = append(out, t)
	}

	// 如果out的length小于4则返回false.
	if len(out) < 4 {
		return false
	}
	cd.thisRange = infRange
	cd.top = out

	thisLevel := s.compactStatus.levels[cd.thisLevel.level]
	thisLevel.ranges = append(thisLevel.ranges, infRange)
	for _, t := range out {
		s.compactStatus.tables[t.ID()] = struct{}{}
	}
	cd.t.fileSz[0] = math.MaxUint32

	return true
}

func (s *LevelsController) fillTablesL0ToLBase(cd *compactDef) bool {
	if cd.nextLevel.level == 0 {
		panic("Base level can't be zero.")
	}
	if cd.p.adjusted > 0.0 && cd.p.adjusted < 1.0 {
		return false
	}

	cd.lockLevels()
	defer cd.unlockLevels()
	// top为当前level的所有table
	top := cd.thisLevel.tables
	if len(top) == 0 {
		return false
	}
	var out []*table.Table
	if len(cd.dropPrefixes) > 0 {
		out = top
	} else {
		var kr keyRange
		for _, t := range top {
			dkr := getKeyRange(t)
			if kr.overlapsWith(dkr) {
				out = append(out, t)
				kr.extend(dkr)
			} else {
				break
			}
		}
	}
	cd.thisRange = getKeyRange(out...)
	cd.top = out

	left, right := cd.nextLevel.overlappingTables(levelHandlerRLocked{}, cd.thisRange)
	cd.bot = make([]*table.Table, right-left)
	copy(cd.bot, cd.nextLevel.tables[left:right])

	if len(cd.bot) == 0 {
		cd.nextRange = cd.thisRange
	} else {
		cd.nextRange = getKeyRange(cd.bot...)
	}
	return s.compactStatus.compareAndAdd(thisAndNextLevelRLocked{}, *cd)
}

type thisAndNextLevelRLocked struct{}

func (s *LevelsController) appendIterator(iters []y.Iterator, opt *IteratorOptions) []y.Iterator {
	return nil
}

func (s *LevelsController) keySplits(numPerTable int, prefix []byte) []string {
	return nil
}

func (s *LevelsController) getLevelInfo() []LevelInfo {
	return nil
}

func (s *LevelsController) verifyChecksum() error {
	return nil
}

func (s *LevelsController) AppendIterators(iters []y.Iterator, opt *IteratorOptions) []y.Iterator {
	return nil
}

// 从L0到LN里面查找key
func (s *LevelsController) Get(key []byte, maxVs y.ValueStruct, startLevel int) (y.ValueStruct, error) {
	if s.kv.IsClosed() {
		return y.ValueStruct{}, ErrDBClosed
	}
	version := y.ParseTs(key)
	for _, h := range s.levels {
		if h.level < startLevel {
			continue
		}
		vs, err := h.get(key)
		if err != nil {
			return y.ValueStruct{}, y.Wrapf(err, "get key: %q", key)
		}
		if vs.Value == nil && vs.Meta == 0 {
			continue
		}
		y.NumBytesReadsLSMAdd(s.kv.Opt.MetricsEnabled, int64(len(vs.Value)))
		if vs.Version == version {
			return vs, nil
		}
		if maxVs.Version < vs.Version {
			maxVs = vs
		}
	}

	if len(maxVs.Value) > 0 {
		y.NumGetsWithResultAdd(s.kv.Opt.MetricsEnabled, 1)
	}
	return maxVs, nil
}
func HasAnyPrefixes(s []byte, listOfPrefixes [][]byte) bool {
	return false
}

// open db时被调用
func newLevelsController(db *DB, mf *Manifest) (*LevelsController, error) {
	y.AssertTrue(db.Opt.NumLevelZeroTablesStall > db.Opt.NumLevelZeroTables)
	levelsController := &LevelsController{
		kv:     db,
		levels: make([]*levelHandler, db.Opt.MaxLevels),
	}
	levelsController.compactStatus.tables = make(map[uint64]struct{})
	levelsController.compactStatus.levels = make([]*levelCompactStatus, db.Opt.MaxLevels)

	//为每一层都创建level handler
	for i := 0; i < db.Opt.MaxLevels; i++ {
		levelsController.levels[i] = newLevelHandler(db, i)
		levelsController.compactStatus.levels[i] = new(levelCompactStatus)
	}

	if db.Opt.InMemory {
		return levelsController, nil
	}
	//移除manifest没有fileId的sst文件
	if err := revertToManifest(db, mf, getIDMap(db.Opt.Dir)); err != nil {
		return nil, err
	}
	var mu sync.Mutex
	tables := make([][]*table.Table, db.Opt.MaxLevels)
	var maxFileID uint64
	//设置限流器
	throttle := y.NewThrottle(3)
	start := time.Now()
	var numOpened atomic.Int32
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()

	//遍历manifest 的table
	for fileID, tf := range mf.Tables {
		fName := table.NewFileName(fileID, db.Opt.Dir)
		select {
		case <-tick.C:
			db.Opt.Infof("%d tables out of %d opened in %levelsController\n", numOpened.Load(), len(mf.Tables), time.Since(start).Round(time.Millisecond))
		default:

		}
		// 尝试抢占限流器的资源
		if err := throttle.Do(); err != nil {
			closeAllTables(tables)
			return nil, err
		}
		if fileID > maxFileID {
			maxFileID = fileID
		}
		go func(fName string, tf TableManifest) {
			var rerr error
			defer func() {
				throttle.Done(rerr)
				numOpened.Add(1)
			}()
			dk, err := db.Register.DataKey(tf.KeyID)
			if err != nil {
				rerr = y.Wrapf(err, "Error while reading datakey")
				return
			}
			tOpt := BuildTableOptions(db)
			tOpt.Compression = tf.Compression
			tOpt.DataKey = dk

			mf, err := z.OpenMmapFile(fName, db.Opt.GetFileFlags(), 0)
			if err != nil {
				rerr = y.Wrapf(err, "Opending file: %q", fName)
				return
			}

			t, err := table.OpenTable(mf, tOpt)
			if err != nil {
				if strings.HasSuffix(err.Error(), "CHECKSUM_MISMATCH:") {
					db.Opt.Errorf(err.Error())
					db.Opt.Errorf("Ignoring table %levelsController", mf.Fd.Name())
				} else {
					rerr = y.Wrapf(err, "Opening table: %q", fName)
				}
				return
			}
			mu.Lock()
			tables[tf.Level] = append(tables[tf.Level], t)
			mu.Unlock()
		}(fName, tf)
	}

	if err := throttle.Finish(); err != nil {
		closeAllTables(tables)
		return nil, err
	}
	db.Opt.Infof("All %d tables opened in %levelsController\n", numOpened.Load(), time.Since(start).Round(time.Millisecond))
	levelsController.nextFileID.Store(maxFileID + 1)
	for i, tbls := range tables {
		levelsController.levels[i].initTables(tbls)
	}
	if err := levelsController.validate(); err != nil {
		_ = levelsController.cleanupLevels()
		return nil, y.Wrap(err, "Level validation")
	}

	if err := SyncDir(db.Opt.Dir); err != nil {
		_ = levelsController.close()
		return nil, err
	}
	return levelsController, nil
}

// 将memTable 删除,tables是个二维数组
// TODO: 需要了解下具体实现细节
func closeAllTables(tables [][]*table.Table) {
	for _, tableSlice := range tables {
		for _, table := range tableSlice {
			_ = table.Close(-1)
		}
	}
}

// idMap中的id均是sst文件.
func revertToManifest(kv *DB, mf *Manifest, idMap map[uint64]struct{}) error {
	for id := range mf.Tables {
		//检查manifest中的fileID是否有对应的fileID.sst
		//如果没有则返回
		if _, ok := idMap[id]; !ok {
			return fmt.Errorf("file does not exist for table %d", id)
		}
	}
	//如果manifest的fileID有sst文件
	for id := range idMap {
		//如果sst不在MANIFEST中就重新创建一个文件，
		//map 使用遍历分别是key-value,ok指的是value
		//如果manifest的table中没有改fileId
		if _, ok := mf.Tables[id]; !ok {
			kv.Opt.Debugf("Table file %d not referenced in MANIFEST\n", id)
			//如果Manifest中没有该fileId,计算fileId的完整路径
			//然后删除该文件.
			fileName := table.NewFileName(id, kv.Opt.Dir)
			if err := os.Remove(fileName); err != nil {
				return y.Wrapf(err, "While removing table %d", id)
			}
		}
	}
	return nil
}

// 读取DB的dir目录
func getIDMap(dir string) map[uint64]struct{} {
	//读取DB dir所有的文件信息,除了文件还有目录
	fileInfos, err := os.ReadDir(dir)
	y.Check(err)
	idMap := make(map[uint64]struct{})
	for _, info := range fileInfos {
		//如果是目录，则继续读取
		if info.IsDir() {
			continue
		}
		//如果是文件，则读取文件名对应的文件ID
		fileID, ok := table.ParseFileID(info.Name())
		if !ok {
			//如果没找到sst,则继续遍历查找.
			continue
		}
		//将哪些sst的fileID置位空
		idMap[fileID] = struct{}{}
	}
	return idMap
}
func (s *LevelsController) validate() error {
	for _, l := range s.levels {
		if err := l.validate(); err != nil {
			y.Wrap(err, "Levels Controller")
		}
	}
	return nil
}
