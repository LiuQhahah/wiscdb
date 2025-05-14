package main

import (
	"bytes"
	"errors"
	"expvar"
	"fmt"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/dgraph-io/ristretto/v2/z"
	"math"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
	"wiscdb/fb"
	"wiscdb/skl"
	"wiscdb/table"
	"wiscdb/y"
)

var (
	wiscPrefix  = []byte("!wisc!")
	txnKey      = []byte("!wisc!txn")
	bannedNsKey = []byte("!wisc!banned")
)

type DB struct {
	testOnlyExtensions
	lock          sync.RWMutex
	dirLockGuard  *directoryLockGuard
	valueDisGuard *directoryLockGuard
	closers       closer
	mt            *memTable
	Opt           Options
	registry      *KeyRegistry
	imm           []*memTable
	vlog          valueLog
	lc            *LevelsController
	flushChan     chan *memTable

	threshold        *vlogThreshold
	orc              *oracle
	blockWrites      atomic.Int32
	writeCh          chan *request
	pub              *publisher
	Register         *KeyRegistry
	blockCache       *ristretto.Cache[[]byte, *table.Block]
	indexCache       *ristretto.Cache[uint64, *fb.TableIndex]
	allocPool        *z.AllocatorPool
	Manifest         *manifestFile
	nextMemFid       int
	bannedNamespaces *lockedKeys
}

type lockedKeys struct {
	sync.RWMutex
	keys map[uint64]struct{}
}

func (lk *lockedKeys) add(key uint64) {

}
func (lk *lockedKeys) has(key uint64) bool {
	return false
}
func (lk *lockedKeys) all() []uint64 {
	lk.RLock()
	defer lk.RUnlock()
	keys := make([]uint64, len(lk.keys))
	for key, _ := range lk.keys {
		keys = append(keys, key)
	}
	return keys
}

type closer struct {
	updateSize  *z.Closer
	compactors  *z.Closer
	memTable    *z.Closer
	writes      *z.Closer
	valueGC     *z.Closer
	pub         *z.Closer
	cacheHealth *z.Closer
}

func (db *DB) openMemTables(opt Options) error {
	return nil
}

// create new memtable
func (db *DB) newMemTable() (*memTable, error) {
	mt, err := db.openMemTable(db.nextMemFid, os.O_CREATE|os.O_RDWR)
	if err == z.NewFile {
		db.nextMemFid++
		return mt, nil
	}
	if err != nil {
		db.Opt.Errorf("Got error: %v for id: %d\n", err, db.nextMemFid)
		return nil, y.Wrapf(err, "newMemTable")
	}

	errMessage := fmt.Sprintf("File %s already exists", mt.wal.Fd.Name())
	return mt, errors.New(errMessage)
}

// create/create memtable with file id.
func (db *DB) openMemTable(fid, flags int) (*memTable, error) {
	filePath := db.getFilePathWithFid(fid)
	s := skl.NewSkipList(arenaSize(db.Opt))
	memtable := &memTable{
		sl:  s,
		opt: db.Opt,
		buf: &bytes.Buffer{},
	}
	memtable.wal = &writeAheadLog{
		path:     filePath,
		fid:      uint32(fid),
		opt:      db.Opt,
		writeAt:  vlogHeaderSize,
		registry: db.registry,
	}
	//open write ahead log
	walError := memtable.wal.open(filePath, flags, 2*db.Opt.MemTableSize)
	if walError != nil && !errors.Is(walError, z.NewFile) {
		return nil, y.Wrapf(walError, "While opening memtable: %s", filePath)
	}

	if errors.Is(walError, z.NewFile) {
		return memtable, walError
	}
	err := memtable.UpdateSkipList()
	return memtable, y.Wrapf(err, "while updating skiplist")
}

// return file name with file if
// 1: filepath package
// 2. fmt package
// 3. const memTableExt:=".mem"
func (db *DB) getFilePathWithFid(fid int) string {
	return filepath.Join(db.Opt.Dir, fmt.Sprintf("%05d%s", fid, memFileExt))
}

// DB的view函数用于查看DB中的值
func (db *DB) View(fb func(txn *Txn) error) error {
	if db.IsClosed() {
		return ErrDBClosed
	}
	var txn *Txn
	if db.Opt.managedTxns {
		txn = db.NewTransactionAt(math.MaxUint64, false)
	} else {
		// 是否是属于更新的transaction,查询 update设置成false
		txn = db.NewTransaction(false)
	}
	defer txn.Discard()
	return fb(txn)
}

func (db *DB) NewTransactionAt(readTs uint64, update bool) *Txn {
	if !db.Opt.managedTxns {
		panic("Cannot use NewTransactionAt with managedDB=false, use NewTransaction instead")
	}
	txn := db.newTransaction(update, true)
	txn.readTs = readTs
	return txn
}
func (db *DB) NewTransaction(update bool) *Txn {
	return db.newTransaction(update, false)
}

func (db *DB) newTransaction(update, isManaged bool) *Txn {
	if db.Opt.ReadOnly && update {
		update = false
	}
	txn := &Txn{
		update: update,
		db:     db,
		count:  1,
		size:   int64(len(txnKey) + 10),
	}
	// 如果是更新的transaction,则需要检查冲突
	if update {
		if db.Opt.DetectConflicts {
			txn.conflictKeys = make(map[uint64]struct{})
		}
		//如果是更新的transaction 需要初始化pendingWrite
		txn.pendingWrites = make(map[string]*Entry)
	}
	if !isManaged {
		txn.readTs = db.orc.readTs()
	}
	return txn
}

func (db *DB) Update(fn func(txn *Txn) error) error {
	return nil
}

func (db *DB) initBannedNamespaces() error {
	return nil
}

func (db *DB) IsClosed() bool {
	return false
}

func (db *DB) isBanned(key []byte) error {
	if db.Opt.NamespaceOffset < 0 {
		return nil
	}
	if len(key) <= db.Opt.NamespaceOffset+8 {
		return nil
	}
	if db.bannedNamespaces.has(y.BytesToU64(key[db.Opt.NamespaceOffset:])) {
		return ErrBannedKey
	}
	return nil
}

func (db *DB) BanNamespace(ns uint64) error {
	return nil
}

var requestPool = sync.Pool{
	New: func() interface{} {
		return new(request)
	},
}

func (db *DB) sendToWriteCh(entries []*Entry) (*request, error) {
	if db.blockWrites.Load() == 1 {
		return nil, ErrBlockedWrites
	}
	var count, size int64
	for _, e := range entries {
		size += e.estimateSizeAndSetThreshold(db.valueThreshold())
		count++
	}
	y.NumBytesWrittenUserAdd(db.Opt.MetricsEnabled, size)
	if count >= db.Opt.maxBatchCount || size >= db.Opt.maxBatchSize {
		return nil, ErrTxnTooBig
	}
	req := requestPool.Get().(*request)
	req.reset()
	req.Entries = entries
	req.Wg.Add(1)
	req.IncrRef()
	db.writeCh <- req
	y.NumPutsAdd(db.Opt.MetricsEnabled, int64(len(entries)))
	return req, nil
}

func (db *DB) valueThreshold() int64 {
	return db.threshold.valueThreshold.Load()
}

func (db *DB) newStream() *Stream {
	return &Stream{}
}

func (db *DB) NewStream() *Stream {
	return db.newStream()
}

func (db *DB) NewStreamAt(readTs uint64) *Stream {
	stream := db.newStream()
	stream.readTs = readTs
	return stream
}

/*
*

	arena Size: MemTableSize 指的是内存中的memtable

默认大小： 64<<20(64MiB) + 10 + 88*100

TODO: maxBatchSize 的作用
TODO: MaxNodeSize指的是什么
default size: MemTableSize:64MiB.
*/
func arenaSize(opt Options) int64 {
	return opt.MemTableSize + opt.maxBatchSize + opt.maxBatchCount*int64(skl.MaxNodeSize)
}

// 返回DB struct中的memtable
func (db *DB) getMemTables() ([]*memTable, func()) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	var memTables []*memTable
	if !db.Opt.ReadOnly {
		memTables = append(memTables, db.mt)
		db.mt.IncrRef()
	}

	last := len(db.imm) - 1

	for i := range db.imm {
		memTables = append(memTables, db.imm[last-i])
		db.imm[last-i].IncrRef()
	}
	//后一个参数是一个钩子函数，用来处理完毕后调用的函数.
	return memTables, func() {
		for _, tbl := range memTables {
			tbl.DecrRef()
		}
	}
}

func (db *DB) StreamDB(outOptions Options) error {
	return nil
}

func (db *DB) writeRequests(reqs []*request) error {
	if len(reqs) == 0 {
		return nil
	}
	done := func(err error) {
		for _, r := range reqs {
			r.Err = err
			r.Wg.Done()
		}
	}

	db.Opt.Debugf("writeRequests called. Writing to value log")
	// TODO: 没有得到是否写到本地文件中
	err := db.vlog.write(reqs)
	if err != nil {
		done(err)
		return err
	}

	db.Opt.Debugf("Writing to memtable")
	var count int
	for _, b := range reqs {
		if len(b.Entries) == 0 {
			continue
		}
		count += len(b.Entries)
		var i uint64
		var err error
		for err = db.ensureRoomForWrite(); err == errNoRoom; err = db.ensureRoomForWrite() {
			i++
			for i%100 == 0 {
				db.Opt.Debugf("Making room for writes")
			}
			time.Sleep(10 * time.Millisecond)
		}
		if err != nil {
			done(err)
			return y.Wrap(err, "writeRequests")
		}
		if err := db.writeToLSM(b); err != nil {
			done(err)
			return y.Wrap(err, "writeRequests")
		}
	}

	db.Opt.Debugf("Sending updates to subscribers")
	db.pub.sendUpdates(reqs)
	done(nil)
	db.Opt.Debugf("%s entries written", count)
	return nil
}

func (db *DB) writeToLSM(b *request) error {
	if !db.Opt.InMemory && len(b.Ptrs) != len(b.Entries) {
		errMessage := fmt.Sprintf("Ptrs and Entries don't match: %+v", b)
		return errors.New(errMessage)
	}
	for i, entry := range b.Entries {
		var err error
		if entry.skipVlogAndSetThreshold(db.valueThreshold()) {
			err = db.mt.Put(entry.Key, y.ValueStruct{
				Value:     entry.Value,
				Meta:      entry.meta &^ bitValuePointer,
				UserMeta:  entry.UserMeta,
				ExpiresAt: entry.ExpiresAt,
			})
		} else {
			//writer pointer to memTable
			err = db.mt.Put(entry.Key, y.ValueStruct{
				Value:     b.Ptrs[i].Encode(),
				Meta:      entry.meta | bitValuePointer,
				UserMeta:  entry.UserMeta,
				ExpiresAt: entry.ExpiresAt,
			})
		}

		if err != nil {
			return y.Wrapf(err, "while writing to memTable")
		}
	}

	if db.Opt.SyncWrites {
		return db.mt.SyncWAL()
	}

	return nil
}

var errNoRoom = errors.New("No Room for write")

func (db *DB) ensureRoomForWrite() error {
	var err error
	db.lock.Lock()
	defer db.lock.Unlock()

	y.AssertTrue(db.mt != nil)
	if !db.mt.isFull() {
		return nil
	}

	select {
	//如果db的memTable channel中有值
	//	如果db的memtable有值，就会执行case里的函数体
	case db.flushChan <- db.mt:
		db.Opt.Debugf("Flushing memtable, mt.size=%d size of flushChan: %d\n", db.mt.sl.MemSize(), len(db.flushChan))
		//将memtable写到immutable中
		db.imm = append(db.imm, db.mt)
		//将memtable的值写到flushchan中，同事创建新的memtable
		db.mt, err = db.newMemTable()
		if err != nil {
			return y.Wrapf(err, "cannot create new mem table")
		}
		return nil
	//	什么情况下会执行到default？
	//flushChan 已满，无法接收新的刷盘任务,flushChan是个没有buffer channel
	//系统负载过高，刷盘速度低于写入速度
	//如果负责消费 flushChan 的后台刷盘协程意外退出或死锁，flushChan 会一直处于满状态,所有后续写入尝试都会直接走 default 分支
	//即使 flushChan 未满，但如果后台刷盘协程因为某些原因（如磁盘 I/O 慢、资源竞争等）无法及时消费 flushChan 中的任务，也可能导致 flushChan 积压
	default:
		return errNoRoom
	}
}

// DB 被调用后就会执行函数，用于一直监听事件.
// 会通过协程来调用，不会阻塞主线程但是可以一直在后台执行
func (db *DB) doWrites(lc *z.Closer) {
	defer lc.Done()
	pendingCh := make(chan struct{}, 1)
	//创建一个闭包函数与直接调用db的writeRequests相比添加了一条处理完后返回pendingCh的输出
	//pendingCh用来标识阻塞，成功后就可以返回表示处理完毕
	writeRequests := func(reqs []*request) {
		if err := db.writeRequests(reqs); err != nil {
			db.Opt.Errorf("writeRequests: %v", err)
		}
		//会一直阻塞知道pendingCh的channel中有数据
		//有数据输出就会执行，没有数据就会一直等待.
		<-pendingCh
	}

	reqLen := new(expvar.Int)
	y.PendingWritesSet(db.Opt.MetricsEnabled, db.Opt.Dir, reqLen)
	//初始化slice,超过10后会自动扩容
	reqs := make([]*request, 0, 10)
	for {
		var r *request
		select {
		case r = <-db.writeCh: //writechan中有值就会做一个赋值处理,交给下面append来做
		case <-lc.HasBeenClosed():
			goto closedCase
		}
		for {
			reqs = append(reqs, r)
			//set metric
			reqLen.Set(int64(len(reqs)))
			//if reqs's length greater than 3000, will call write method
			if len(reqs) >= 3*kvWriteChCapacity {
				//第一个条件中的阻塞发送确保信号被接收后才会继续
				pendingCh <- struct{}{}
				goto writeCase
			}
			select {
			case r = <-db.writeCh:
			//	pendingCh是缓冲为1的通道，每次只能容纳一个信号
			//第二个条件中的非阻塞发送确保不会重复发送
			case pendingCh <- struct{}{}:
				goto writeCase
			case <-lc.HasBeenClosed():
				goto closedCase
			}
		}
	closedCase:
		for {
			select {
			case r = <-db.writeCh:
				reqs = append(reqs, r)
			default:
				pendingCh <- struct{}{}
				writeRequests(reqs)
				return
			}
		}
	writeCase:
		go writeRequests(reqs)
		//new request
		reqs = make([]*request, 10)
		reqLen.Set(0)
	}
}

const (
	kvWriteChCapacity = 1000
)

// 会被go goroutine调用作为一个后台一直执行的协程在运行.
func (db *DB) flushMemTable(lc *z.Closer) {
	defer lc.Done()
	for mt := range db.flushChan {
		if mt == nil {
			continue
		}
		for {
			if err := db.handleMemTableFlush(mt, nil); err != nil {
				db.Opt.Errorf("error flushing memtable to disk: %v, retrying", err)
				time.Sleep(time.Second)
				continue
			}
			db.lock.Lock()
			//将memTable的内容写到immutable中
			y.AssertTrue(mt == db.imm[0])
			//更新immutable
			// TODO: immutable最后一个变成了什么？
			db.imm = db.imm[1:]
			mt.DecrRef()
			db.lock.Unlock()
			break
		}
	}
}

func (db *DB) handleMemTableFlush(mt *memTable, dropPrefixes [][]byte) error {
	bOpts := BuildTableOptions(db)
	itr := mt.sl.NewUniIterator(false)
	//将memtable的值写到buildTable中
	builder := buildL0Table(itr, nil, bOpts)
	defer builder.Close()

	if builder.Empty() {
		builder.CutDownBuilder()
		return nil
	}
	fileID := db.lc.ReserveFileID()
	var tbl *table.Table
	var err error
	if db.Opt.InMemory {
		data := builder.CutDownBuilder()
		tbl, err = table.OpenInMemoryTable(data, fileID, &bOpts)
	} else {
		//创建table
		tbl, err = table.CreateTable(table.NewFileName(fileID, db.Opt.Dir), builder)
	}
	if err != nil {
		return y.Wrapf(err, "error while creating table")
	}

	//将memTable的table加到Level 0中
	err = db.lc.AddLevel0Table(tbl)
	_ = tbl.DecrRef()

	return err
}

func buildL0Table(iter y.Iterator, dropPrefixes [][]byte, bOpts table.Options) *table.Builder {
	defer iter.Close()
	b := table.NewTableBuilder(bOpts)
	for iter.ReWind(); iter.Valid(); iter.Next() {
		if len(dropPrefixes) > 0 && HasAnyPrefixes(iter.Key(), dropPrefixes) {
			continue
		}
		vs := iter.Value()
		var vp valuePointer
		//判断当前value是存储的是值还是value的地址
		if vs.Meta&bitValuePointer > 0 {
			//如果是地址还需要解码
			vp.Decode(vs.Value)
		}
		//从迭代器中得到key value写到table builder中
		b.Add(iter.Key(), iter.Value(), vp.Len)
	}
	return b
}
func BuildTableOptions(db *DB) table.Options {
	opt := db.Opt
	dk, err := db.registry.LatestDataKey()
	y.Check(err)
	return table.Options{
		ReadOnly:             opt.ReadOnly,
		MetricsEnabled:       db.Opt.MetricsEnabled,
		TableSize:            uint64(opt.BaseTableSize),
		BlockSize:            opt.BlockSize,
		BloomFalsePositive:   opt.BloomFalsePositive,
		ChkMode:              opt.ChecksumVerificationMode,
		Compression:          opt.Compression,
		ZSTDCompressionLevel: opt.ZSTDCompressionLevel,
		BlockCache:           db.blockCache,
		IndexCache:           db.indexCache,
		AllocPool:            db.allocPool,
		DataKey:              dk,
	}
}

func (db *DB) Get(key []byte) (y.ValueStruct, error) {
	if db.IsClosed() {
		return y.ValueStruct{}, ErrDBClosed
	}
	tables, decr := db.getMemTables()
	defer decr()

	var maxVs y.ValueStruct
	//查询的key后8个字节带有查询的时间
	version := y.ParseTs(key)
	y.NumGetsAdd(db.Opt.MetricsEnabled, 1)
	for i := 0; i < len(tables); i++ {
		vs := tables[i].sl.Get(key)
		y.NumMemtableGetsAdd(db.Opt.MetricsEnabled, 1)
		if vs.Meta == 0 && vs.Value == nil {
			continue
		}
		if vs.Version == version {
			y.NumGetsWithResultAdd(db.Opt.MetricsEnabled, 1)
			return vs, nil
		}
		if maxVs.Version < vs.Version {
			maxVs = vs
		}
	}
	return db.lc.Get(key, maxVs, 0)
}
