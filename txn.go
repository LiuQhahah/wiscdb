package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto/v2/z"
	"strconv"
	"sync"
	"sync/atomic"
	"wiscdb/table"
	"wiscdb/y"
)

type Txn struct {
	readTs          uint64
	commitTs        uint64
	size            int64
	count           int64
	db              *DB
	reads           []uint64
	conflictKeys    map[uint64]struct{}
	readsLock       sync.Mutex
	pendingWrites   map[string]*Entry
	duplicateWrites []*Entry
	numIterators    atomic.Int32
	discarded       bool
	doneRead        bool
	update          bool
}

func (txn *Txn) newPendingWritesIterator(reversed bool) *pendingWritesIterator {
	return &pendingWritesIterator{}
}

func (txn *Txn) checkTransactionCountAndSize(e *Entry) error {
	count := txn.count + 1
	size := txn.size + e.estimateSizeAndSetThreshold(txn.db.valueThreshold()) + 10
	if count >= txn.db.Opt.maxBatchCount || size >= txn.db.Opt.maxBatchSize {
		return ErrTxnTooBig
	}
	txn.count, txn.size = count, size
	return nil
}

func (txn *Txn) Discard() {
	if txn.discarded {
		return
	}
	if txn.numIterators.Load() > 0 {
		panic("Unclosed iterator at time of Txn.Discard.")
	}
	txn.discarded = true
	if !txn.db.orc.isManaged {
		txn.db.orc.doneRead(txn)
	}

}

func (txn *Txn) modify(e *Entry) error {
	const maxKeySize = 65000
	switch {
	case !txn.update:
		return ErrReadOnlyTxn
	case !txn.discarded:
		return ErrDiscardedTxn
	case len(e.Key) == 0:
		return ErrEmptyKey
	case bytes.HasPrefix(e.Key, wiscPrefix):
		return ErrInvalidKey
	case len(e.Key) > maxKeySize:
		return exceedsSize("Key", maxKeySize, e.Key)
	case int64(len(e.Value)) > txn.db.Opt.ValueLogFileSize:
		return exceedsSize("Value", txn.db.Opt.ValueLogFileSize, e.Key)
	case txn.db.Opt.InMemory && int64(len(e.Value)) > txn.db.valueThreshold():
		return exceedsSize("Value", txn.db.valueThreshold(), e.Value)
	}

	if err := txn.db.isBanned(e.Key); err != nil {
		return err
	}
	if err := txn.checkTransactionCountAndSize(e); err != nil {
		return err
	}

	// 如果开了冲突检测，则将更改的key写到事务的conflictKeys中,供后续事务提交时检测冲突
	if txn.db.Opt.DetectConflicts {
		fp := z.MemHash(e.Key)
		// 将所有的key写到冲突检测map中
		txn.conflictKeys[fp] = struct{}{}
	}

	if oldEntry, ok := txn.pendingWrites[string(e.Key)]; ok && oldEntry.version != e.version {
		txn.duplicateWrites = append(txn.duplicateWrites, oldEntry)
	}
	txn.pendingWrites[string(e.Key)] = e
	return nil
}

func exceedsSize(prefix string, max int64, key []byte) error {
	errMessage := fmt.Sprintf("%s with size %d exceeded %d limit. %s:\n%s", prefix, len(key), max, prefix, hex.Dump(key[:1<<10]))
	return errors.New(errMessage)
}

// 将 需要读取的key写到事务的reads中
func (txn *Txn) addReadKey(key []byte) {
	if txn.update {
		fp := z.MemHash(key)
		txn.readsLock.Lock()
		// 将key的hash值写到事务的reads中,可以用来做冲突检测
		txn.reads = append(txn.reads, fp)
		txn.readsLock.Unlock()
	}
}

// 用户执行完多次Set方法后 开始提交请求
func (txn *Txn) commitAndSend() (func() error, error) {
	orc := txn.db.orc
	orc.writeChLock.Lock()
	defer orc.writeChLock.Unlock()

	//创建提交时间戳
	commitTs, conflict := orc.newCommitTs(txn)
	if conflict {
		return nil, ErrConflict
	}
	keepTogether := true

	// 如果entry的version已经设置了version,则keepTogether为false
	setVersion := func(e *Entry) {
		if e.version == 0 {
			e.version = commitTs
		} else {
			keepTogether = false
		}
	}
	//遍历事务的pendingWrites的值,设置Entry的version是commitTs
	for _, e := range txn.pendingWrites {
		setVersion(e)
	}

	// 遍历事务的duplicateWrites,设置Entry的version为commitTs
	for _, e := range txn.duplicateWrites {
		setVersion(e)
	}

	// 构建entry list,长度为pendingWrite+duplicateWrite另外再+1
	entries := make([]*Entry, 0, len(txn.pendingWrites)+len(txn.duplicateWrites)+1)

	processEntry := func(e *Entry) {
		e.Key = y.KeyWithTs(e.Key, e.version)
		if keepTogether {
			e.meta |= bitTxn
		}
		entries = append(entries, e)
	}

	// 更新完version,调用processEntry处理pendingWrites中的Entry
	for _, e := range txn.pendingWrites {
		processEntry(e)
	}

	// 更新完version,调用processEntry处理duplicateWrites中的Entry
	for _, e := range txn.duplicateWrites {
		processEntry(e)
	}

	// TODO: keepTogether意味着什么?
	if keepTogether {
		y.AssertTrue(commitTs != 0)
		// 构建Entry key为wiscDB并且携带提交时间戳
		e := &Entry{
			Key: y.KeyWithTs(txnKey, commitTs),
			// value存储出的是提交时间戳
			Value: []byte(strconv.FormatUint(commitTs, 10)),
			meta:  bitFinTxn, // meta表FIN标志
		}
		// 这是entries最后一个entries用来表示报文的结束位
		// 这也是entries长度+1的原因
		//	entries := make([]*Entry, 0, len(txn.pendingWrites)+len(txn.duplicateWrites)+1)
		entries = append(entries, e)
	}
	//将这些entries发送给writeCh中
	req, err := txn.db.sendToWriteCh(entries)
	// 无论是否有错都会在oracle中完成该事务提交
	if err != nil {
		orc.doneCommit(commitTs)
		return nil, err
	}
	// 提交后就在等待,由上游调用函数执行
	ret := func() error {
		err := req.Wait()
		orc.doneCommit(commitTs)
		return err
	}

	return ret, nil
}

func (txn *Txn) commitPreCheck() error {
	if txn.discarded {
		return errors.New("the txn has been discarded")
	}
	keepTogether := true
	for _, e := range txn.pendingWrites {
		if e.version != 0 {
			keepTogether = false
		}
	}
	if keepTogether && txn.db.Opt.managedTxns && txn.commitTs == 0 {
		return errors.New("CommitTs cannot be zero. Please use commitAt instead")
	}
	return nil
}

func (txn *Txn) Commit() error {
	//1. check pendingWrites
	if len(txn.pendingWrites) == 0 {
		txn.Discard()
		return nil
	}
	if err := txn.commitPreCheck(); err != nil {
		return err
	}
	defer txn.Discard()
	txnCb, err := txn.commitAndSend()
	if err != nil {
		return err
	}
	// 执行txnCb()函数等待协程完成
	return txnCb()
}

func (txn *Txn) CommitWith(cb func(err error)) {

}

func (txn *Txn) ReadTs() uint64 {
	return 0
}

func (txn *Txn) SetEntry(e *Entry) error {
	return txn.modify(e)
}

func (txn *Txn) Set(key, value []byte) error {
	return txn.SetEntry(NewEntry(key, value))
}

func (txn *Txn) Get(key []byte) (item *Item, err error) {
	if len(key) == 0 {
		return nil, ErrEmptyKey
	} else if txn.discarded {
		return nil, ErrDiscardedTxn
	}
	//判断该可以是不是属于被禁止的key
	if err := txn.db.isBanned(key); err != nil {
		return nil, err
	}
	//分配内存空间
	item = new(Item)
	//如果是写的事务 ，则需要查询事务的pendingWrite中的值.
	if txn.update {
		//如果pendingWrite中有值并且是需求的key
		//同时该key也没有被删除或者过期则直接返回.
		//这是获取的第一层
		if e, has := txn.pendingWrites[string(key)]; has && bytes.Equal(key, e.Key) {
			if table.IsDeletedOrExpired(e.meta, e.ExpiresAt) {
				item.meta = e.meta
				item.val = e.Value
				item.userMeta = e.UserMeta
				item.key = key
				item.status = int(prefetched)
				item.version = txn.readTs
				item.expiresAt = e.ExpiresAt
				return item, nil
			}
			txn.addReadKey(key)
		}
	}
	seek := y.KeyWithTs(key, txn.readTs)
	// 查询key带上了读取的时间戳:一切都是报文
	vs, err := txn.db.Get(seek)
	if err != nil {
		return nil, y.Wrapf(err, "DB::Get Key: %q", key)
	}
	if vs.Value == nil && vs.Meta == 0 {
		return nil, ErrBannedKey
	}
	if table.IsDeletedOrExpired(vs.Meta, vs.ExpiresAt) {
		return nil, ErrKeyNotFound
	}
	item.key = key
	item.version = vs.Version
	item.meta = vs.Meta
	item.userMeta = vs.UserMeta
	item.vptr = y.SafeCopy(item.vptr, vs.Value)
	item.txn = txn
	item.expiresAt = vs.ExpiresAt
	return item, nil
}

type prefetchStatus int

const (
	prefetched prefetchStatus = iota + 1
)

func (txn *Txn) Delete(key []byte) error {
	return nil
}

func (txn *Txn) NewIterator(opt IteratorOptions) *Iterator {
	if txn.discarded {
		panic(ErrDiscardedTxn)
	}
	if txn.db.IsClosed() {
		panic(ErrDBClosed)
	}
	y.NumIteratorsCreatedAdd(txn.db.Opt.MetricsEnabled, 1)

	//创建迭代器时，将事务中迭代器加1
	txn.numIterators.Add(1)
	// TODO: If Prefix is set, only pick those memtables which have keys with the prefix.
	tables, decr := txn.db.getMemTables()
	//在函数退出时调用.
	defer decr()
	txn.db.vlog.incrIteratorCount()
	var iters []y.Iterator
	if itr := txn.newPendingWritesIterator(opt.Reverse); itr != nil {
		iters = append(iters, itr)
	}
	//遍历db中的mem table, 在找到跳表中的迭代器
	for i := 0; i < len(tables); i++ {
		iters = append(iters, tables[i].sl.NewUniIterator(opt.Reverse))
	}
	iters = txn.db.lc.AppendIterators(iters, &opt)
	res := &Iterator{
		txn:    txn,
		iitr:   table.NewMergeIterator(iters, opt.Reverse),
		opt:    opt,
		readTs: txn.readTs,
	}
	return res
}
