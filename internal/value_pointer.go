package internal

import "unsafe"

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
	return false
}

func (p valuePointer) IsZero() bool {
	return false
}

func (p valuePointer) Encode() []byte {
	return nil
}

func (p *valuePointer) Decode(b []byte) {

}
