package level

import (
	"fmt"
	"github.com/dgraph-io/ristretto/v2/z"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"wiscdb/internal"
	"wiscdb/pb"
	"wiscdb/table"
	"wiscdb/y"
)

type LevelsController struct {
	nextFileID    atomic.Uint64
	l0stallMs     atomic.Int64
	levels        []*levelHandler
	kv            *internal.DB
	compactStatus compactStatus
}

func newLevelController(db *internal.DB, mf *internal.Manifest) (*LevelsController, error) {
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
	return nil
}

func (s *LevelsController) AddLevel0Table(t *table.Table) error {
	if !t.IsInMemory {
		err := s.kv.Manifest.AddChanges([]*pb.ManifestChange{
			internal.NewCreateChange(t.ID(), 0, t.KeyID(), t.CompressionType()),
		})
		if err != nil {
			return err
		}
	}

	for !s.levels[0].tryAddLevel0Table(t) {
		timeStart := time.Now()
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

}

func (s *LevelsController) runCompactor(id int, lc *z.Closer) {

}

func (s *LevelsController) pickCompactLevels(priosBuffer []compactionPriority) (prios []compactionPriority) {
	return nil
}

func (s *LevelsController) levelTargets() targets {
	return targets{}
}

func (s *LevelsController) lastLevel() *levelHandler {
	return nil
}

func (s *LevelsController) doCompact(id int, p compactionPriority) error {
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

func (s *LevelsController) fillTablesL0(cd *compactDef) bool {
	return false
}
func (s *LevelsController) fillTablesL0ToL0(cd *compactDef) bool {
	return false
}

func (s *LevelsController) fillTablesL0ToLBase(cd *compactDef) bool {
	return false
}

func (s *LevelsController) appendIterator(iters []y.Iterator, opt *internal.IteratorOptions) []y.Iterator {
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

func (s *LevelsController) AppendIterators(iters []y.Iterator, opt *internal.IteratorOptions) []y.Iterator {
	return nil
}

func HasAnyPrefixes(s []byte, listOfPrefixes [][]byte) bool {
	return false
}

// open db时被调用
func newLevelsController(db *internal.DB, mf *internal.Manifest) (*LevelsController, error) {
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
		if err := throttle.Do(); err != nil {
			closeAllTables(tables)
			return nil, err
		}
		if fileID > maxFileID {
			maxFileID = fileID
		}
		go func(fName string, tf internal.TableManifest) {
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
			tOpt := internal.BuildTableOptions(db)
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

	if err := internal.SyncDir(db.Opt.Dir); err != nil {
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
func revertToManifest(kv *internal.DB, mf *internal.Manifest, idMap map[uint64]struct{}) error {
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
	return nil
}
