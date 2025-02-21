package y

import (
	"bytes"
	"encoding/binary"
)

type ValueStruct struct {
	Meta      byte
	UserMeta  byte
	ExpiresAt uint64
	Value     []byte
	Version   uint64
}

// 无人调用
func (v *ValueStruct) EncodeTo(buf *bytes.Buffer) {
	buf.WriteByte(v.Meta)
	buf.WriteByte(v.UserMeta)
	var enc [binary.MaxVarintLen64]byte
	sz := binary.PutUvarint(enc[:], v.ExpiresAt)
	buf.Write(enc[:sz])
	buf.Write(v.Value)
}

// 不通过buffer，没有写到buffer中
func (v *ValueStruct) Encode(b []byte) uint32 {
	b[0] = v.Meta
	b[1] = v.UserMeta
	sz := binary.PutUvarint(b[2:], v.ExpiresAt)
	n := copy(b[2+sz:], v.Value)
	return uint32(2 + sz + n)
}
func (v *ValueStruct) Decode(b []byte) {
	v.Meta = b[0]
	v.UserMeta = b[1]
	var sz int
	v.ExpiresAt, sz = binary.Uvarint(b[2:])
	v.Value = b[2+sz:]

}

func (v *ValueStruct) EncodedSize() uint32 {
	sz := binary.PutUvarint(make([]byte, 10), v.ExpiresAt)
	return uint32(2 + sz + len(v.Value))
}
