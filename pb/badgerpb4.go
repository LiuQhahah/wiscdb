package pb

import (
	"github.com/dgraph-io/ristretto/v2/z"
	"google.golang.org/protobuf/reflect/protoreflect"
	"wiscdb/options"
)

type DataKey struct {
	KeyId    uint64
	Data     []byte
	Iv       []byte
	CreateAt int64
}

type KV struct {
	Key        []byte
	Value      []byte
	UserMeta   []byte
	Version    uint64
	ExpireAt   uint64
	Meta       []byte
	StreamId   uint32
	StreamDone bool
}

func NewKV(alloc *z.Allocator) *KV {
	return &KV{}
}

type KVList struct {
	Kv       []*KV
	AllocRef uint64
}

type EncryptionAlgo int32

const (
	EncryptionAlgo_aes EncryptionAlgo = iota
)

type ManifestChangeOperation int32

const (
	ManifestChange_CREATE ManifestChangeOperation = iota
	ManifestChange_DELETE
)

type ManifestChangeSet struct {
}

type ManifestChange struct {
	Id             uint64
	Op             ManifestChangeOperation
	Level          uint32
	KeyId          uint64
	EncryptionAlgo EncryptionAlgo
	Compression    uint32
}

func newCreateChange(id uint64, level int, keyID uint64, c options.CompressionType) *ManifestChange {
	return &ManifestChange{}
}

func newDeleteChange(id uint64) *ManifestChange {
	return &ManifestChange{}
}

type ChecksumAlgorithm int

const (
	CheckSumCRC32C ChecksumAlgorithm = iota
	CheckSumXXHash64
)

type CheckSum struct {
	Algo ChecksumAlgorithm
	Sum  uint64
}

func (x *CheckSum) Reset() {

}

func (x *CheckSum) String() string {
	return ""
}

func (*CheckSum) ProtoMessage() {

}
func (x *CheckSum) ProtoReflect() protoreflect.Message {
	return nil
}
func (*CheckSum) Descriptor() ([]byte, []int) {
	return nil, nil
}

func (x *CheckSum) GetAlgo() ChecksumAlgorithm {
	return 0
}
func (x *CheckSum) GetSum() uint64 {
	return 0
}

type Match struct {
	Prefix      []byte
	IgnoreBytes string
}
