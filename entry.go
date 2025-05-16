package main

import (
	"fmt"
	"time"
	"unsafe"
)

const (
	//含义丢弃较早版本
	bitDiscardEarlierVersion byte = 1 << 2
	//合并Entry
	bitMergeEntry byte = 1 << 3
)

// 存储在value log 中的Entry，主要是value的值，但是还包含了额外的信息
// 字节存储存储的key和value，存储过期时间、存储版本信息，
// 存储偏移量，存储usermeta信息，存储meta信息
// 存储header的长度以及val的门限值
type Entry struct {
	Key          []byte
	Value        []byte
	ExpiresAt    uint64
	version      uint64
	offset       uint32
	UserMeta     byte
	meta         byte //meta会用来记载该key是否被delete
	hlen         int
	valThreshold int64 // 设置一个阈值,超过阈值则只写入value的地址,小于阈值就直接存储value到内存中
}

// 判断当前Entry为空
func (e *Entry) isEmpty() bool {
	return len(e.Key) == 0
}

// 返回当前Entry的大小，Entry大小包含key和value的长度以及两个meta字节的长度.
func (e *Entry) estimateSizeAndSetThreshold(threshold int64) int64 {
	if e.valThreshold == 0 {
		e.valThreshold = threshold
	}
	k := int64(len(e.Key))
	v := int64(len(e.Value))
	if v < e.valThreshold {
		return k + v + 2
	}
	//如果value的大小大于等于设定的value threshold,那么就用valuePointer的指作为返回值.
	// 如果key对应的value值大于设定的阈值,那么value的值就是value地址的值了此时就需要用int64(unsafe.Sizeof(valuePointer{}))
	// TODO: 2的含义是什么
	return k + int64(unsafe.Sizeof(valuePointer{})) + 2
}

// 判断value的长度是否超过了threshold的门限值
func (e *Entry) checkValueThreshold(threshold int64) bool {
	if e.valThreshold == 0 {
		e.valThreshold = threshold
	}
	return int64(len(e.Value)) < e.valThreshold
}

// 为什么是指传递，因为print函数不会更改Entry中的内容，因此要使用指传递
// 引用传递一般都会修改struct中的内容.
func (e Entry) print(prefix string) {
	fmt.Printf("%s Key: %s Meta: %d UserMeta: %d Offset: %d len(val)=%d",
		prefix, e.Key, e.meta, e.UserMeta, e.offset, len(e.Value))
}

// 传入key和value构建Entry
func NewEntry(key, value []byte) *Entry {
	return &Entry{
		Key: key, Value: value,
	}
}

// 传入meta,加到Entry中
// TODO: usermeta存储的是什么
func (e *Entry) WithUserMeta(meta byte) *Entry {
	e.UserMeta = meta
	return e
}

// 给出TTL值，返回expire 时间
func (e *Entry) WithTTL(dur time.Duration) *Entry {
	e.ExpiresAt = uint64(time.Now().Add(dur).Unix())
	return e
}

// 指的是meta这个参数
func (e *Entry) withMergeBit() *Entry {
	e.meta = bitMergeEntry
	return e
}

// 指的是meta
func (e *Entry) WithDiscard() *Entry {
	e.meta = bitDiscardEarlierVersion
	return e
}

func (e *Entry) skipVlogAndSetThreshold(threshold int64) bool {
	if e.valThreshold == 0 {
		e.valThreshold = threshold
	}
	return int64(len(e.Value)) < e.valThreshold
}
