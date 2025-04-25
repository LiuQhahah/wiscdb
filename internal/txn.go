package internal

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto/v2/z"
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
	if count >= txn.db.opt.maxBatchCount || size >= txn.db.opt.maxBatchSize {
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
	case int64(len(e.Value)) > txn.db.opt.ValueLogFileSize:
		return exceedsSize("Value", txn.db.opt.ValueLogFileSize, e.Key)
	case txn.db.opt.InMemory && int64(len(e.Value)) > txn.db.valueThreshold():
		return exceedsSize("Value", txn.db.valueThreshold(), e.Value)
	}

	if err := txn.db.isBanned(e.Key); err != nil {
		return err
	}
	if err := txn.checkTransactionCountAndSize(e); err != nil {
		return err
	}

	if txn.db.opt.DetectConflicts {
		fp := z.MemHash(e.Key)
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
	//1. check pendingWrites
	if len(txn.pendingWrites) == 0 {
		txn.Discard()
		return nil
	}
	return nil
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
	return nil, nil
}

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
