package internal

import (
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

func (txn *Txn) checkSize(e *Entry) error {
	return nil
}

func (txn *Txn) Discard() {

}

func (txn *Txn) modify(e *Entry) error {
	return nil
}

func (txn *Txn) addReadKey(key []byte) {

}

func (txn *Txn) commitAndSend() (func() error, error) {
	return func() error {
		return nil
	}, nil
}

func (txn *Txn) commitPreCheck() error {
	return nil
}

func (txn *Txn) Commit() error {
	return nil
}

func (txn *Txn) CommitWith(cb func(err error)) {

}

func (txn *Txn) ReadTs() uint64 {
	return 0
}

func (txn *Txn) SetEntry(e *Entry) error {
	return nil
}

func (txn *Txn) Set(key, value []byte) error {
	return nil
}

func (txn *Txn) Get(key []byte) (item *Item, err error) {
	return nil, nil
}

func (txn *Txn) Delete(key []byte) error {
	return nil
}

func (txn *Txn) NewIterator(opt IteratorOptions) *Iterator {
	if txn.discarded {
		panic(ErrDiscardTxn)
	}
	if txn.db.IsClosed() {
		panic(ErrDBClosed)
	}
	y.NumIteratorsCreatedAdd(txn.db.opt.MetricsEnabled, 1)

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
