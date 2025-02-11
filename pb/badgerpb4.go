package pb

import (
	"github.com/dgraph-io/ristretto/v2/z"
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
}

type Match struct {
	Prefix      []byte
	IgnoreBytes string
}
