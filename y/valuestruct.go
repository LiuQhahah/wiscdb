package y

import "bytes"

type ValueStruct struct {
	Meta      byte
	UserMeta  byte
	ExpiresAt uint64
	Value     []byte
	Version   uint64
}

func (v *ValueStruct) EncodeTo(buf *bytes.Buffer) {

}

func (v *ValueStruct) Encode(b []byte) uint32 {
	return 0
}
func (v *ValueStruct) Decode(b []byte) {

}

func (v *ValueStruct) EncodedSize() uint32 {
	return 0
}
