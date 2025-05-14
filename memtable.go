package main

import (
	"bytes"
	"errors"
	"wiscdb/skl"
	"wiscdb/y"
)

type memTable struct {
	sl         *skl.SkipList
	wal        *writeAheadLog
	maxVersion uint64
	opt        Options
	buf        *bytes.Buffer
}

func (mt *memTable) isFull() bool {
	if mt.sl.MemSize() >= mt.opt.MemTableSize {
		return true
	}
	if mt.opt.InMemory {
		return false
	}
	return int64(mt.wal.writeAt) >= mt.opt.MemTableSize
}

// 将key value写到skiplist中
func (mt *memTable) Put(key []byte, value y.ValueStruct) error {
	entry := &Entry{
		Key:       key,
		Value:     value.Value,
		UserMeta:  value.UserMeta,
		meta:      value.Meta,
		ExpiresAt: value.ExpiresAt,
	}
	if mt.wal != nil {
		if err := mt.wal.writeEntry(mt.buf, entry, mt.opt); err != nil {
			return y.Wrapf(err, "cannot write entry to WAL file")
		}
	}

	if entry.meta&bitFinTxn > 0 {
		return nil
	}

	mt.sl.Put(key, value)

	//更新下version
	if ts := y.ParseTs(entry.Key); ts > mt.maxVersion {
		mt.maxVersion = ts
	}
	y.NumBytesWrittenToL0Add(mt.opt.MetricsEnabled, entry.estimateSizeAndSetThreshold(mt.opt.ValueThreshold))
	return nil
}

func (mt *memTable) IncrRef() {
	mt.sl.IncrRef()
}

func (mt *memTable) DecrRef() {
	mt.sl.DecrRef()
}

// 回放函数，返回函数
// 第一参数是Entry，按照
func (mt *memTable) replayFunction(opt Options) func(Entry, valuePointer) error {
	return func(entry Entry, pointer valuePointer) error {
		opt.Logger.Debugf("First key=%q\n", entry.Key)
		//不停解析value log file中的key的事务时间戳，同时更新memtable中的最大版本
		if ts := y.ParseTs(entry.Key); ts > mt.maxVersion {
			mt.maxVersion = ts
		}
		value := y.ValueStruct{
			Value:    entry.Value,
			UserMeta: entry.UserMeta,
			Meta:     entry.meta,
			Version:  entry.version,
		}
		//将解析出来的entry写到跳表中，key为entry的key，value就是valuestruct，包含value以及各种meta信息
		mt.sl.Put(entry.Key, value)
		return nil
	}
}

// ErrTruncateNeeded is returned when the value log gets corrupt, and requires truncation of
// corrupt data to allow Badger to run properly.
var ErrTruncateNeeded = errors.New("Log truncate required to run DB. This might result in data loss")

// 将wal文件中的信息迭代到memtable中.
// 更新跳表
func (mt *memTable) UpdateSkipList() error {
	// 从0开始遍历
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
	return mt.wal.Sync()
}

// const有type 和untype之分,加了string指的是有typeed的constant.
const memFileExt string = ".mem"
