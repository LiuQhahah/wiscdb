package internal

import "unsafe"

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
