package table

import (
	"bytes"
	"fmt"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/dgraph-io/ristretto/v2/z"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
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
	BlockSize            int         // default 4MB
	DataKey              *pb.DataKey // DataKey的作用: 判断table是否有加密
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
	UncompressedSize  uint32
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
	bd := builder.CutDoneBuildData()
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
	//sync到 memtable中
	if err := z.Msync(mf.Data); err != nil {
		return nil, y.Wrapf(err, "while calling msync on %s", fName)
	}
	return OpenTable(mf, *builder.opts)
}

// 将memtable写到table中
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

// 根据判断table的选项有没有datakey来判断是否需要解密
func (t *Table) shouldDecrypt() bool {
	return t.opt.DataKey != nil
}

// 从table中读取字节数组,读取的逻辑是按照偏移量+文件尺寸
func (t *Table) read(off, sz int) ([]byte, error) {
	return t.Bytes(off, sz)
}

// 按照偏移量+文件尺寸读取字节数组
func (t *Table) readNoFail(off, sz int) []byte {
	res, err := t.read(off, sz)
	y.Check(err)
	return res
}

//	biggest 和smallest的作用
//
// biggest和smallest作为一张SST表的属性方便用于query.
func (t *Table) initBiggestAndSmallest() error {
	defer func() {
		//r 指的是发生了panic返回的内容,比如panic("error") 那么string error就是r
		if r := recover(); r != nil {
			//r!=nil 则表明有panic发生 输出索引信息
			var debugBuf bytes.Buffer
			defer func() {
				panic(fmt.Sprintf("%s\n== Recovered ==\n", debugBuf.String()))
			}()

			count := 0
			for i := len(t.Data) - 1; i >= 0; i-- {
				if t.Data[i] != 0 {
					break
				}
				count++
			}

			fmt.Fprintf(&debugBuf, "\n== Recovering from initIndex crash ==\n")
			fmt.Fprintf(&debugBuf, "File info: [ID:%d, Size: %d, Zeros:%d]\n", t.id, t.tableSize, count)
			fmt.Fprintf(&debugBuf, "isEnrypted: %v", t.shouldDecrypt())
			readPos := t.tableSize

			readPos -= 4
			//读取四个字节,拿到校验合的长度
			buf := t.readNoFail(readPos, 4)
			checksumLen := int(y.BytesToU32(buf))
			fmt.Fprintf(&debugBuf, "checksumLen: %d", checksumLen)

			checksum := &pb.CheckSum{}
			readPos -= checksumLen
			//读取校验和并解析校验和
			buf = t.readNoFail(readPos, checksumLen)
			_ = proto.Unmarshal(buf, checksum)
			fmt.Fprintf(&debugBuf, "checksum: %+v", checksum)

			readPos -= 4
			//读取索引,用4个字节表示索引
			buf = t.readNoFail(readPos, 4)
			indexLen := int(y.BytesToU32(buf))
			fmt.Fprintf(&debugBuf, "indexLen: %d ", indexLen)

			readPos -= t.indexLen
			t.indexStart = readPos
			//读取索引值
			indexData := t.readNoFail(readPos, t.indexLen)
			fmt.Fprintf(&debugBuf, "index: %v ", indexData)
		}
	}()
	var err error
	var ko *fb.BlockOffset
	// 初始化索引
	if ko, err = t.initIndex(); err != nil {
		return y.Wrapf(err, "failed to read index")
	}
	//将keyBytes的内容复制到smallest中
	//找到索引中的key
	t.smallest = y.Copy(ko.KeyBytes())
	//为table创建迭代器是为着找到最大的值
	it2 := t.NewIterator(REVERSED | NOCACHE)
	defer it2.Close()
	//读取table中的值
	it2.Rewind()
	if !it2.Valid() {
		return y.Wrapf(it2.err, "failed to initialize biggest for table %s", t.Filename())
	}
	//将迭代器的Key写到biggest中
	t.biggest = y.Copy(it2.Key())
	return nil
}

func (t *Table) NewIterator(opt int) *Iterator {
	t.IncrRef()
	ti := &Iterator{
		t:   t,
		opt: opt,
	}
	return ti
}

func (t *Table) Filename() string {
	return ""
}

// TODO: 如何构造table 的索引 key是key value时存储的位置吗？
func (t *Table) initIndex() (*fb.BlockOffset, error) {
	//从后向前读取
	readPos := t.tableSize
	readPos -= 4
	buf := t.readNoFail(readPos, 4)
	//校验和长度的长度占据4个字节因此读到buf中并解析出校验和的长度
	checksumLen := int(y.BytesToU32(buf))
	if checksumLen < 0 {
		return nil, errors.New("checksum  length less than zero. Data corrupted")
	}
	readPos -= checksumLen
	buf = t.readNoFail(readPos, checksumLen)
	expectedChk := &pb.CheckSum{}
	//将校验和写到expectedChk struct中
	if err := proto.Unmarshal(buf, expectedChk); err != nil {
		return nil, err
	}

	readPos -= 4
	buf = t.readNoFail(readPos, 4)
	t.indexLen = int(y.BytesToU32(buf))

	readPos -= t.indexLen
	t.indexStart = readPos
	data := t.readNoFail(readPos, t.indexLen)

	if err := y.VerifyCheckSum(data, expectedChk); err != nil {
		return nil, y.Wrapf(err, "failed to verify checksum for this table: %s", t.Filename())
	}

	index, err := t.readTableIndex()
	if err != nil {
		return nil, err
	}
	if !t.shouldDecrypt() {
		t._index = index
	}
	t._cheap = &cheapIndex{
		MaxVersion:        index.MaxVersion(),
		KeyCount:          index.KeyCount(),
		UncompressedSize:  index.UncompressedSize(),
		OnDiskSize:        index.OnDiskSize(),
		OffsetsLength:     index.OffsetsLength(),
		BloomFilterLength: index.BloomFilterLength(),
	}
	t.hasBloomFilter = len(index.BloomFilterBytes()) > 0
	var bo fb.BlockOffset
	y.AssertTrue(index.Offsets(&bo, 0))
	return &bo, nil
}

// 读取table的索引
func (t *Table) readTableIndex() (*fb.TableIndex, error) {
	//在table中读取索引
	data := t.readNoFail(t.indexStart, t.indexLen)
	var err error
	if t.shouldDecrypt() {
		if data, err = t.decrypt(data, false); err != nil {
			return nil, y.Wrapf(err, "Error while decrpting table index for the table %d in readTableIndex", t.id)
		}
	}
	return fb.GetRootAsTableIndex(data, 0), nil
}

func (t *Table) decrypt(data []byte, viaCalloc bool) ([]byte, error) {
	return nil, nil
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

const fileSuffix = ".sst"

// 对应有sst的文件，返回true及其文件id,否则返回false
func ParseFileID(name string) (uint64, bool) {
	name = filepath.Base(name)
	if !strings.HasSuffix(name, fileSuffix) {
		return 0, false
	}
	name = strings.TrimSuffix(name, fileSuffix)
	id, err := strconv.Atoi(name)
	if err != nil {
		return 0, false
	}
	y.AssertTrue(id >= 0)
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
	return filepath.Join(dir, IDToFilename(id))
}

// 根据file ID返回对应的sst文件.
func IDToFilename(id uint64) string {
	return fmt.Sprintf("%06d", id) + fileSuffix
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

func (t *Table) Smallest() []byte {
	return t.smallest
}
func (t *Table) Biggest() []byte {
	return t.biggest
}
