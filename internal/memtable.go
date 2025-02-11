package internal

import (
	"bytes"
	"wiscdb/skl"
	"wiscdb/y"
)

type memTable struct {
	sl         *skl.SkipList
	wal        *logFile
	maxVersion uint64
	opt        Options
	buf        *bytes.Buffer
}

func (mt *memTable) isFull() bool {
	return false
}

func (mt *memTable) Put(key []byte, value y.ValueStruct) error {
	return nil
}

func (mt *memTable) IncrRef() {
	mt.sl.IncrRef()
}

func (mt *memTable) DecrRef() {

}

func (mt *memTable) replayFunction(opt Options) func(Entry, valuePointer) error {
	return func(entry Entry, pointer valuePointer) error {
		return nil
	}
}

func (mt *memTable) UpdateSkipList() error {
	return nil
}

func (mt *memTable) SyncWAL() error {
	return nil
}

const memFileExt string = ".mem"
