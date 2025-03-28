package internal

import (
	"bytes"
	"errors"
	"wiscdb/skl"
	"wiscdb/y"
)

type memTable struct {
	sl         *skl.SkipList
	wal        *valueLogFile
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

// 回放函数，返回函数
func (mt *memTable) replayFunction(opt Options) func(Entry, valuePointer) error {
	return func(entry Entry, pointer valuePointer) error {
		value := y.ValueStruct{
			Value:    entry.Value,
			UserMeta: entry.UserMeta,
			Meta:     entry.meta,
			Version:  entry.version,
		}
		mt.sl.Put(entry.Key, value)
		return nil
	}
}

// ErrTruncateNeeded is returned when the value log gets corrupt, and requires truncation of
// corrupt data to allow Badger to run properly.
var ErrTruncateNeeded = errors.New("Log truncate required to run DB. This might result in data loss")

// 将wal文件中的信息迭代到memtable中.
func (mt *memTable) UpdateSkipList() error {
	endOff, err := mt.wal.iterate(true, 0, mt.replayFunction(mt.opt))
	if err != nil {
		return y.Wrapf(err, "while iterating wal: %s", mt.wal.Fd.Name())
	}
	if endOff < mt.wal.size.Load() && mt.wal.opt.ReadOnly {
		return y.Wrapf(ErrTruncateNeeded, "end offset: %d < size: %d", endOff, mt.wal.size.Load())
	}
	return mt.wal.Truncate(int64(endOff))
}

func (mt *memTable) SyncWAL() error {
	return nil
}

// const有type 和untype之分,加了string指的是有typeed的constant.
const memFileExt string = ".mem"
