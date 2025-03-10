package table

import (
	"github.com/dgraph-io/ristretto/v2"
	"github.com/dgraph-io/ristretto/v2/z"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
	"wiscdb/fb"
	"wiscdb/options"
	"wiscdb/pb"
)

type Table struct {
	sync.Mutex
	*z.MmapFile
	tableSize         int
	_index            *fb.TableIndex
	_cheap            *cheapIndex
	ref               atomic.Int32
	smallest, biggest []byte
	id                uint64
	Checksum          []byte
	CreatedAt         time.Time
	indexStart        int
	indexLen          int
	hasBloomFilter    bool
	IsInMemory        bool
	opt               *Options
}

type Options struct {
	ReadOnly           bool
	MetricsEnabled     bool
	TableSize          uint64
	tableCapacity      uint64
	ChkMode            options.ChecksumVerificationMode
	BloomFalsePositive float64
	BlockSize          int
	DataKey            *pb.DataKey
	Compression        options.CompressionType
	BlockCache         *ristretto.Cache[[]byte, *Block]
}

type cheapIndex struct {
	MaxVersion        uint64
	KeyCount          uint32
	UnCompressionSize uint32
	OnDiskSize        uint32
	BloomFilterLength int
	OffsetsLength     int
}

type Block struct {
	offset            int
	data              []byte
	checksum          []byte
	entriesIndexStart int
	entryOffsets      []uint32
	chkLen            int
	freeMe            bool
	ref               atomic.Int32
}

var NumBlocks atomic.Int32

func (b *Block) incrRef() bool {
	ref := b.ref.Load()
	return b.ref.CompareAndSwap(ref, ref+1)
}

func (b *Block) decrRef() {
	ref := b.ref.Load()
	if ref == 0 {
		return
	}

	b.ref.CompareAndSwap(ref, ref-1)
}

const intSize = int(unsafe.Sizeof(int(0)))

func (b *Block) size() int64 {
	return int64(3*intSize + cap(b.data) + cap(b.entryOffsets)*4 + cap(b.checksum))
}

func (b *Block) verifyCheckSum() error {

	return nil
}

func CreateTable(fName string, builder *Builder) {

}

func OpenTable(mf *z.MmapFile, opts Options) (*Table, error) {
	return nil, nil
}

/*
*
input file path
return file ID
*/
func ParseFileID(filePath string) (uint64, bool) {
	fileName := strings.TrimSuffix(filePath, ".sst")
	id, err := strconv.Atoi(fileName)
	if err != nil {
		return 0, false
	}
	return uint64(id), true
}

func (t *Table) offsets(ko *fb.BlockOffset, i int) bool {
	return false
}

func (t *Table) block(idx int, useCache bool) (*Block, error) {
	return nil, nil
}

func (t *Table) decompress(b *Block) error {
	return nil
}

type TableInterface interface {
	Smallest() []byte
	Biggest() []byte
	DoesNotHave(hash uint32) bool
	MaxVersion() uint64
}

func (t *Table) KeySplit(n int, prefix []byte) []string {
	return nil
}

func OpenInMemoryTable(data []byte, id uint64, opt *Options) (*Table, error) {
	return &Table{}, nil
}

func NewFileName(id uint64, dir string) string {
	return ""
}

func IDToFilename(id uint64) string {
	return ""
}
