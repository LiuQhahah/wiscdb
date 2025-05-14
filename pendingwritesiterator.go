package main

import (
	"wiscdb/y"
)

type pendingWritesIterator struct {
	entries  []*Entry
	nextIdx  int
	readTs   uint64
	reversed bool
}

func (pi *pendingWritesIterator) Next() {

}
func (pi *pendingWritesIterator) ReWind() {

}

func (pi *pendingWritesIterator) Seek(key []byte) {

}

func (pi *pendingWritesIterator) Valid() bool {
	return false
}

func (pi *pendingWritesIterator) Key() []byte {
	return nil
}

func (pi *pendingWritesIterator) Value() y.ValueStruct {
	return y.ValueStruct{}
}

func (pi *pendingWritesIterator) Close() error {
	return nil
}
