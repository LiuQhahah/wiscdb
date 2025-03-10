package internal

import (
	"sync"
	"sync/atomic"
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
