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
	// 指定MB创建跳表
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

// 对request创建pool,可以重复利用
var requestPool = sync.Pool{
	New: func() interface{} {
		return new(request)
	},
}

func (db *DB) sendToWriteCh(entries []*Entry) (*request, error) {
	// 判断当前db的blockWrite状态
	if db.blockWrites.Load() == 1 {
		return nil, ErrBlockedWrites
	}
	var count, size int64
	// 遍历entries,计算存储的尺寸以及entries个数
	for _, e := range entries {
		size += e.estimateSizeAndSetThreshold(db.valueThreshold())
		count++
	}
	y.NumBytesWrittenUserAdd(db.Opt.MetricsEnabled, size)
	// 限制批处理的个数以及批处理的尺寸
	// 批处理的个数和尺寸会决定内存的大小
	if count >= db.Opt.maxBatchCount || size >= db.Opt.maxBatchSize {
		return nil, ErrTxnTooBig
	}
	// 从池中取得request
	req := requestPool.Get().(*request)
	// 重置request
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

maxBatchSize 的作用 : 批处理是一个事务能处理的最大字节数
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
			if i%100 == 0 {
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
		if entry.skipVlogAndSetThreshold(db.valueThreshold(), int64(len(entry.Value))) {
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
	//	将memtable转移到db的immutable中
	//	成功刷到db的immutable channel中
	//	case为true才会执行里面的内容
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

//批量写入（Batching）
//
//收集多个请求后一次性写入，减少 I/O 次数，提高吞吐量。
//
//并发控制（pendingCh）
//
//限制同时只有一个 writeRequests 在执行，避免资源竞争。
//
//优雅关闭（lc.HasBeenClosed()）
//
//关闭时确保所有 pending 请求被处理，避免数据丢失。
//
//动态调整写入触发条件
//
//基于请求堆积量（3*kvWriteChCapacity）和令牌可用性（pendingCh）决定何时写入。
//
//监控（expvar.Int）
//
//暴露 reqLen 指标，便于观察堆积情况。

// DB 被调用后就会执行函数，用于一直监听事件.
// 会通过协程来调用，不会阻塞主线程但是可以一直在后台执行

// 需求	实现方式
// 持续消费请求	case r = <-db.writeCh 逐步读取
// 低延迟写入	pendingCh 空闲时立即触发
// 高吞吐写入	队列满时批量触发
// 避免阻塞	select 多路监听
func (db *DB) doWrites(lc *z.Closer) {
	// 确保在退出时通知 Closer（用于优雅关闭）
	defer lc.Done()
	// 本质是一个二进制信号量
	// 控制并发写入的令牌桶（最大并发=1）
	pendingCh := make(chan struct{}, 1)
	//创建一个闭包函数与直接调用db的writeRequests相比添加了一条处理完后返回pendingCh的输出
	//pendingCh用来标识阻塞，成功后就可以返回表示处理完毕
	writeRequests := func(reqs []*request) {
		if err := db.writeRequests(reqs); err != nil { // 执行批量写入
			db.Opt.Errorf("writeRequests: %v", err) // 出错时打印日志
		}
		//会一直阻塞知道pendingCh的channel中有数据
		//有数据输出就会执行，没有数据就会一直等待.
		// 释放令牌，允许新的写入
		<-pendingCh
	}

	// 用于监控 pending 请求数量
	reqLen := new(expvar.Int)
	y.PendingWritesSet(db.Opt.MetricsEnabled, db.Opt.Dir, reqLen)
	//初始化slice,超过10后会自动扩容
	// 初始化请求缓冲区（容量=10）
	reqs := make([]*request, 0, 10)
	for {
		var r *request
		// 在执行下面for循环时,先扫一个两个channel中的值
		select {
		case r = <-db.writeCh: // 如果db.writeCh有值,则将值赋值给r
		case <-lc.HasBeenClosed(): // 如果收到关闭信号，进入清理逻辑
			goto closedCase
		}
		for {
			reqs = append(reqs, r) // 加入待处理列表
			//set metric
			reqLen.Set(int64(len(reqs)))
			//if reqs's length greater than 3000, will call write method
			if len(reqs) >= 3*kvWriteChCapacity { // 如果堆积过多，立即触发写入
				//第一个条件中的阻塞发送确保信号被接收后才会继续
				// 占用令牌（可能阻塞）
				// 准备处理request时就要占用令牌
				pendingCh <- struct{}{}
				goto writeCase
			}
			select {
			// 连续向writeCh中写入数据
			//
			case r = <-db.writeCh: // 刷新一下db writeCh中的值,如果有貌似并不会写到reqs中,只能等到channel遍历完才会写到reqs中
			//	pendingCh是缓冲为1的通道，每次只能容纳一个信号
			//第二个条件中的非阻塞发送确保不会重复发送
			// 如果pendingCh还有容量,尝试获取令牌
			//则进入写操作，所以至少要向writeCh中写两次，才会去执行写操作
			case pendingCh <- struct{}{}: // 如果能获取令牌，立即触发写入 如果能写入
				goto writeCase
			case <-lc.HasBeenClosed(): // 如果关闭，进入清理
				goto closedCase
			}
		}
	closedCase:
		for {
			select {
			case r = <-db.writeCh: // 继续消费剩余请求
				reqs = append(reqs, r)
			default: // 没有剩余请求时，执行最后一次写入并退出
				pendingCh <- struct{}{} // 占用令牌
				writeRequests(reqs)     // 同步写入（确保所有数据落盘）
				return                  // 退出 goroutine
			}
		}
	writeCase:
		// 异步执行写入
		go writeRequests(reqs)
		// 重置缓冲区
		reqs = make([]*request, 10)
		// 重置监控指标
		reqLen.Set(0)
	}
}

const (
	kvWriteChCapacity = 1000
)

// 会被go goroutine调用作为一个后台一直执行的协程在运行.
func (db *DB) flushMemTable(lc *z.Closer) {
	defer lc.Done()

	// 如果channel不为空就会进入到for 循环
	//函数会一直阻塞在 for mt := range db.flushChan，等待新的 MemTable 刷新任务
	//Channel 行为
	//阻塞读取：如果没有数据可读，循环会阻塞等待
	//
	//自动终止：只有当 channel 被显式关闭 (close(db.flushChan)) 时循环才会结束
	//
	//零值处理：如果 channel 为 nil，循环会永久阻塞
	//不会执行一次就退出：函数会一直运行，直到 db.flushChan 被显式关闭
	//range 循环特性：for range 会持续监听 channel，直到它被关闭
	//消息不丢失：所有通过 db.flushChan <- msg 成功发送的消息都会被 for range 接收到
	//for range 能可靠接收所有消息的保证源于：
	//
	//精心设计的线程安全数据结构
	//
	//严格遵循的 happens-before 规则
	//
	//运行时调度器的深度配合
	//
	//对各种边界条件的完备处理
	//for range 循环能够可靠接收所有通过 db.flushChan <- msg 发送的消息，这一保证背后是 Go 语言 channel 设计的核心机制。让我们深入分析其工作原理：
	//Go 的 channel 本质上是一个带锁的环形队列（有缓冲）或同步点（无缓冲），其关键组成：
	//type hchan struct {
	//    qcount   uint           // 当前队列中元素数量
	//    dataqsiz uint           // 环形队列大小
	//    buf      unsafe.Pointer // 指向环形队列
	//    sendx    uint           // 发送索引
	//    recvx    uint           // 接收索引
	//    lock     mutex          // 互斥锁
	//    recvq    waitq          // 阻塞的接收者队列
	//    sendq    waitq          // 阻塞的发送者队列
	//}
	//消息传递的原子性保证
	//当执行 db.flushChan <- msg 时：
	//
	//获取锁：首先获取 channel 的互斥锁
	//
	//直接交付（fast path）：
	//
	//如果已有接收者在等待（recvq 不为空），直接将消息拷贝到接收方
	//
	//这是无缓冲 channel 的典型情况
	//
	//缓冲写入（有缓冲 channel）：
	//
	//如果缓冲区有空位，将消息存入缓冲区并更新索引
	//
	//阻塞等待（slow path）：
	//
	//无缓冲且无接收者/有缓冲且缓冲区满时
	//
	//将当前 goroutine 加入 sendq 队列并挂起

	// for range 的接收保证
	//for v := range ch 编译后相当于：
	//
	//
	//for v, ok := <-ch; ok; v, ok = <-ch {
	//    // 循环体
	//}
	//接收过程：
	//
	//检查 channel 状态：
	//
	//如果 channel 已关闭且无数据，结束循环
	//
	//直接获取（fast path）：
	//
	//如果缓冲区有数据（qcount > 0），直接取出
	//
	//如果有发送者在等待（sendq 不为空），直接从发送方获取数据
	//
	//阻塞等待（slow path）：
	//
	//无数据时，将当前 goroutine 加入 recvq 队列并挂起

	//4. 可靠性的三大支柱
	//(1) Happens-before 原则
	//Go 内存模型保证：
	//
	//第 n 次发送 happens before 第 n 次接收完成
	//
	//Channel 的关闭 happens before 接收端收到零值
	//
	//(2) 完全同步的队列管理
	//所有操作受 mutex 保护
	//
	//通过 sudog 结构精确跟踪每个阻塞的 goroutine
	//
	//唤醒顺序严格遵循 FIFO
	//
	//(3) 运行时调度器的深度集成
	//当 goroutine 因 channel 操作阻塞时：
	//
	//被放入 sendq 或 recvq 队列
	//
	//调度器将其 G 结构体从运行队列移除
	//
	//当条件满足时，由对方操作将其重新激活

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
		// 做的事情就是为了内存的数据库
		data := builder.CutDownBuilder()
		// 如果是内存的数据库,需要将block memory table中
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
