package table

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/dgraph-io/ristretto/v2/z"
	"github.com/klauspost/compress/snappy"
	"github.com/klauspost/compress/zstd"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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
	BlockSize            int         // default 4KB
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

func (t *Table) MaxVersion() uint64 {
	return t.cheapIndex().MaxVersion
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
func (t *Table) Size() int64 {
	return int64(t.tableSize)
}

func (t *Table) StaleDataSize() uint32 {
	return t.fetchIndex().StaleDataSize()
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

	//初始化索引，直接读取table的索引
	index, err := t.readTableIndex()
	if err != nil {
		return nil, err
	}
	if !t.shouldDecrypt() {
		t._index = index
	}
	//构建table的_cheap 包含key的数量,布隆过滤器的长度
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
	//将索引0返回即最小的Key
	y.AssertTrue(index.Offsets(&bo, 0))
	return &bo, nil
}

// 读取table的索引
func (t *Table) readTableIndex() (*fb.TableIndex, error) {
	//在table中读取索引
	//根据table的索引开始时以及索引的长度读取索引值
	data := t.readNoFail(t.indexStart, t.indexLen)
	var err error
	if t.shouldDecrypt() {
		if data, err = t.decrypt(data, false); err != nil {
			return nil, y.Wrapf(err, "Error while decrpting table index for the table %d in readTableIndex", t.id)
		}
	}
	//根据解密后的数据得到索引
	return fb.GetRootAsTableIndex(data, 0), nil
}

// decrypt 对给定数据进行解密。只有在检查了 shouldDecrypt 后才能调用它。
func (t *Table) decrypt(data []byte, viaCalloc bool) ([]byte, error) {
	iv := data[len(data)-aes.BlockSize:]
	data = data[:len(data)-aes.BlockSize]

	var dst []byte
	if viaCalloc {
		dst = z.Calloc(len(data), "Table.Decrypt")
	} else {
		dst = make([]byte, len(data))
	}
	//使用 AES 加密算法在 CTR (计数器) 模式下进行 XOR 操作
	// 将加密好的内容写到dst中
	if err := y.XORBlock(dst, data, t.opt.DataKey.Data, iv); err != nil {
		return nil, y.Wrapf(err, "while decrypt")
	}
	return dst, nil
}

// table 校验和
func (t *Table) VerifyChecksum() error {
	//读取table 中的索引
	ti := t.fetchIndex()
	for i := 0; i < ti.OffsetsLength(); i++ {
		b, err := t.block(i, true)
		if err != nil {
			return y.Wrapf(err, "checksum validation failed for table: %s, block: %d,offset:%d",
				t.Filename(), i, b.offset)
		}
		defer b.decrRef()

		if !(t.opt.ChkMode == options.OnBlockRead || t.opt.ChkMode == options.OnTableAndBlockRead) {
			// 核实block的校验和
			if err = b.verifyCheckSum(); err != nil {
				return y.Wrapf(err, "checksum validation failed for table: %s, block: %d, offset: %d",
					t.Filename(), i, b.offset)
			}
		}
	}
	return nil
}

func (t *Table) fetchIndex() *fb.TableIndex {
	if !t.shouldDecrypt() {
		return t._index
	}
	if t.opt.IndexCache == nil {
		panic("Index cache must be set for encrypted workloads")
	}
	//一个table有一个id这个id会用来表示indexcache的内容,返回值是tableIndex
	//如果有则直接从cache返回
	if val, ok := t.opt.IndexCache.Get(t.indexKey()); ok && val != nil {
		return val
	}
	//如果cache中没有该id对应的索引,则需要从文件读取table index
	index, err := t.readTableIndex()
	y.Check(err)
	t.opt.IndexCache.Set(t.indexKey(), index, int64(t.indexLen))
	return index
}

// table中的id指示的是索引的key
func (t *Table) indexKey() uint64 {
	return t.id
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
	return t.fetchIndex().Offsets(ko, i)
}

func (t *Table) block(idx int, useCache bool) (*Block, error) {
	y.AssertTruef(idx >= 0, "idx=%d", idx)
	//idx的值不能超过block的长度
	if idx >= t.offsetsLength() {
		return nil, errors.New("block out of index")
	}
	//如果将block存储在cache中来
	if t.opt.BlockCache != nil {
		key := t.blockCacheKey(idx)
		//根据key拿到缓存中的block
		blk, ok := t.opt.BlockCache.Get(key)
		//如果不为空则直接返回blk
		if ok && blk != nil {
			if blk.incrRef() {
				return blk, nil
			}
		}
	}
	var ko fb.BlockOffset
	y.AssertTrue(t.offsets(&ko, idx))
	//blockOffset更新到block中
	blk := &Block{offset: int(ko.Offset())}
	//设置引用加1
	blk.ref.Store(1)
	defer blk.decrRef()
	NumBlocks.Add(1)

	var err error
	//在table中读取block的data
	if blk.data, err = t.read(blk.offset, int(ko.Len())); err != nil {
		return nil, y.Wrapf(err, "failed to read from file: %s at offset: %d, len: %d", t.Fd.Name(), blk.offset, ko.Len())
	}
	if t.shouldDecrypt() {
		if blk.data, err = t.decrypt(blk.data, true); err != nil {
			return nil, err
		}
		blk.freeMe = true
	}

	// 解压table中block的内容
	if err = t.decompress(blk); err != nil {
		return nil, y.Wrapf(err, "failed to decode compressed data in file: %s at offset: %d,len: %d", t.Fd.Name(), blk.offset, ko.Len())
	}

	readPos := len(blk.data) - 4
	//取后四个字节的到校验和的长度
	blk.chkLen = int(y.BytesToU32(blk.data[readPos : readPos+4]))
	if blk.chkLen > len(blk.data) {
		return nil, errors.New("invalid checksum length,Either the data is corrupted or the table options are incorrectly set")
	}

	readPos -= blk.chkLen
	//读取校验和
	blk.checksum = blk.data[readPos : readPos+blk.chkLen]
	readPos -= 4
	//读取entry的数量即key-value个数
	numEntries := int(y.BytesToU32(blk.data[readPos : readPos+4]))
	entriesIndexStart := readPos - (numEntries * 4)
	entriesIndexEnd := entriesIndexStart + numEntries*4

	//读取entry的偏移量
	blk.entryOffsets = y.BytesToU32Slice(blk.data[entriesIndexStart:entriesIndexEnd])
	//读取entry 开始时的索引
	blk.entriesIndexStart = entriesIndexStart

	//读取data值
	blk.data = blk.data[:readPos+4]

	if t.opt.ChkMode == options.OnBlockRead || t.opt.ChkMode == options.OnTableAndBlockRead {
		if err = blk.verifyCheckSum(); err != nil {
			return nil, err
		}
	}
	blk.incrRef()
	//如果需要使用cache,则将block写到cache中
	if useCache && t.opt.BlockCache != nil {
		key := t.blockCacheKey(idx)
		y.AssertTrue(blk.incrRef())
		if !t.opt.BlockCache.Set(key, blk, blk.size()) {
			blk.decrRef()
		}
	}
	return blk, nil
}

func (t *Table) decompress(b *Block) error {
	var dst []byte
	var err error

	src := b.data
	switch t.opt.Compression {
	case options.None:
		return nil
	case options.Snappy:
		if sz, err := snappy.DecodedLen(b.data); err == nil {
			dst = z.Calloc(sz, "Table.Decompress")
		} else {
			dst = z.Calloc(len(b.data)*4, "Table.Decompress")
		}
		b.data, err = snappy.Decode(dst, b.data)
		if err != nil {
			z.Free(dst)
			return y.Wrap(err, "failed to decompress")
		}
	case options.ZSTD:
		sz := int(float64(t.opt.BlockSize) * 1.2)
		var hdr zstd.Header
		if err := hdr.Decode(b.data); err == nil && hdr.HasFCS && hdr.FrameContentSize < uint64(t.opt.BlockSize*2) {
			sz = int(hdr.FrameContentSize)
		}
		dst = z.Calloc(sz, "Table.Decompress")
		b.data, err = y.ZSTDDeCompress(dst, b.data)
		if err != nil {
			z.Free(dst)
			return y.Wrap(err, "failed to decompress")
		}
	default:
		return errors.New("Unsupported compression type")
	}
	if b.freeMe {
		z.Free(src)
		b.freeMe = false
	}

	if len(b.data) > 0 && len(dst) > 0 && &dst[0] != &b.data[0] {
		z.Free(dst)
	} else {
		b.freeMe = true
	}
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

func (t *Table) DoesNotHave(hash uint32) bool {
	return false
}
