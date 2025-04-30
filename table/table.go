package table

import (
	"github.com/dgraph-io/ristretto/v2"
	"github.com/dgraph-io/ristretto/v2/z"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
	"wiscdb/fb"
	"wiscdb/options"
	"wiscdb/pb"
	"wiscdb/y"
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
	ReadOnly             bool
	MetricsEnabled       bool
	TableSize            uint64
	tableCapacity        uint64
	ChkMode              options.ChecksumVerificationMode
	BloomFalsePositive   float64
	BlockSize            int
	DataKey              *pb.DataKey // TODO: DataKey的作用
	Compression          options.CompressionType
	BlockCache           *ristretto.Cache[[]byte, *Block]
	IndexCache           *ristretto.Cache[uint64, *fb.TableIndex]
	AllocPool            *z.AllocatorPool
	ZSTDCompressionLevel int
}

// TODO: cheapIndex的作用是什么
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

// 根据文件名和table builder来创建table
func CreateTable(fName string, builder *Builder) (*Table, error) {
	//将tableBuilder转化为buildData
	bd := builder.Done()
	mf, err := z.OpenMmapFile(fName, os.O_CREATE|os.O_RDWR|os.O_EXCL, bd.Size)
	if err == z.NewFile {

	} else if err != nil {
		return nil, y.Wrapf(err, "while creating table: %s", fName)
	} else {
		return nil, errors.Errorf("file already exists: %s", fName)
	}
	//将buildData写到mMapFile中
	written := bd.Copy(mf.Data)
	y.AssertTrue(written == len(mf.Data))
	if err := z.Msync(mf.Data); err != nil {
		return nil, y.Wrapf(err, "while calling msync on %s", fName)
	}
	return OpenTable(mf, *builder.opts)
}

func OpenTable(mf *z.MmapFile, opts Options) (*Table, error) {
	if opts.BlockSize == 0 && opts.Compression != options.None {
		return nil, errors.New("block size cannot be zero")
	}
	//return mmapfile file info
	fileInfo, err := mf.Fd.Stat()
	if err != nil {
		mf.Close(-1)
		return nil, y.Wrap(err, "")
	}
	filename := fileInfo.Name()
	id, ok := ParseFileID(filename)
	if !ok {
		mf.Close(-1)
		return nil, errors.Errorf("invalid filename: %s", filename)
	}
	t := &Table{
		MmapFile:   mf,
		id:         id,
		opt:        &opts,
		IsInMemory: false,
		tableSize:  int(fileInfo.Size()),
		CreatedAt:  fileInfo.ModTime(),
	}
	t.ref.Store(1)

	if err := t.initBiggestAndSmallest(); err != nil {
		return nil, y.Wrapf(err, "failed to initialize table")
	}

	if opts.ChkMode == options.OnTableRead || opts.ChkMode == options.OnTableAndBlockRead {
		if err := t.VerifyChecksum(); err != nil {
			mf.Close(-1)
			return nil, y.Wrapf(err, "failed to verify checksum")
		}
	}
	return t, nil
}

func (t *Table) initBiggestAndSmallest() error {

	return nil
}

func (t *Table) VerifyChecksum() error {
	return nil
}
func (t *Table) fetchIndex() *fb.TableIndex {
	return nil
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
	mf := &z.MmapFile{
		Data: data,
		Fd:   nil,
	}
	t := &Table{
		MmapFile:   mf,
		opt:        opt,
		tableSize:  len(data),
		IsInMemory: true,
		id:         id,
	}
	t.ref.Store(1)
	if err := t.initBiggestAndSmallest(); err != nil {
		return nil, err
	}
	return t, nil
}

func NewFileName(id uint64, dir string) string {
	return ""
}

func IDToFilename(id uint64) string {
	return ""
}

func (t *Table) IncrRef() {
	t.ref.Add(1)
}

// 垃圾回收，如果该表没有被引用，那么就会从内存中回收
func (t *Table) DecrRef() error {
	newRef := t.ref.Add(-1)
	if newRef == 0 {
		//	GC
		//从缓存中删除key
		for i := 0; i < t.offsetsLength(); i++ {
			t.opt.BlockCache.Del(t.blockCacheKey(i))
		}
		if err := t.Delete(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Table) blockCacheKey(idx int) []byte {
	return nil
}
func (t *Table) offsetsLength() int {
	return 0
}

func (t *Table) cheapIndex() *cheapIndex {
	return t._cheap
}

func (t *Table) ID() uint64 {
	return t.id
}

func (t *Table) KeyID() uint64 {
	if t.opt.DataKey != nil {
		return t.opt.DataKey.KeyId
	}
	return 0
}
func (t *Table) CompressionType() options.CompressionType {
	return t.opt.Compression
}
