package pb

import (
	"github.com/dgraph-io/ristretto/v2/z"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"wiscdb/options"
)

type DataKey struct {
	KeyId    uint64
	Data     []byte
	Iv       []byte
	CreateAt int64
}

var file_wiscdb_proto_msgTypes = make([]protoimpl.MessageInfo, 7)

func (d *DataKey) ProtoReflect() protoreflect.Message {

	mi := &file_wiscdb_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && d != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(d))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(d)
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
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	Changes       []*ManifestChange `protobuf:"bytes,1,rep,name=changes" json:"changes,omitempty"`
}

func (m *ManifestChangeSet) ProtoReflect() protoreflect.Message {
	mi := &file_wiscdb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && m != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(m))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(m)
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
