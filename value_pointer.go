package main

import (
	"encoding/binary"
	"unsafe"
)

// 作用
// 长度指的是entry的key，entry的value，entry的header len以及校验的长度
// fid指的是当前write ahead log
// offset指的是io reader的游标位置
type valuePointer struct {
	Fid    uint32
	Len    uint32
	Offset uint32
}

const vPtrSize = unsafe.Sizeof(valuePointer{})

func (p valuePointer) Less(o valuePointer) bool {
	return p.Fid < o.Fid || p.Offset < o.Offset || p.Len < o.Len
}

func (p valuePointer) IsZero() bool {
	return p.Len == 0
}

func (p valuePointer) Encode() []byte {
	v := make([]byte, vPtrSize)
	index := binary.PutUvarint(v[0:], uint64(p.Fid))
	index += binary.PutUvarint(v[index:], uint64(p.Len))
	index += binary.PutUvarint(v[index:], uint64(p.Offset))
	return v
}

func (p *valuePointer) Decode(b []byte) {
	pFid, sz := binary.Uvarint(b[0:])
	p.Fid = uint32(pFid)
	pLen, sz1 := binary.Uvarint(b[sz:])
	p.Len = uint32(pLen)
	pOffset, _ := binary.Uvarint(b[sz+sz1:])
	p.Offset = uint32(pOffset)
}
