package table

import (
	"crypto/aes"
	"encoding/binary"
	"errors"
	"github.com/dgraph-io/ristretto/v2/z"
	fbs "github.com/google/flatbuffers/go"
	"github.com/klauspost/compress/s2"
	"google.golang.org/protobuf/proto"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"wiscdb/fb"
	"wiscdb/options"
	"wiscdb/pb"
	"wiscdb/y"
)

type Builder struct {
	alloc            *z.Allocator
	curBlock         *bblock
	compressedSize   atomic.Uint32
	uncompressedSize atomic.Uint32
	lenOffsets       uint32
	keyHashes        []uint32
	opts             *Options
	maxVersion       uint64
	onDiskSize       uint32
	staleDataSize    int
	wg               sync.WaitGroup
	blockChan        chan *bblock
	blockList        []*bblock
}

func (b *Builder) calculateCheckSum(data []byte) []byte {
	checksum := pb.CheckSum{
		Sum:  y.CalculateCheckSum(data, pb.CheckSumCRC32C),
		Algo: pb.CheckSumCRC32C,
	}

	chkSum, err := proto.Marshal(&checksum)
	y.Check(err)
	return chkSum
}

/*
Structure of Block.
+-------------------+---------------------+--------------------+--------------+------------------+
| Entry1            | Entry2              | Entry3             | Entry4       | Entry5           |
+-------------------+---------------------+--------------------+--------------+------------------+
| Entry6            | ...                 | ...                | ...          | EntryN           |
+-------------------+---------------------+--------------------+--------------+------------------+
| Block Meta(contains list of offsets used| Block Meta Size    | Block        | Checksum Size    |
| to perform binary search in the block)  | (4 Bytes)          | Checksum     | (4 Bytes)        |
+-----------------------------------------+--------------------+--------------+------------------+
*/
// In case the data is encrypted, the "IV" is added to the end of the block.
// 将curBlock的信息追加到curBlock
//如果curBlock 超过了设定的4KB，就会发送给blockChan中
func (b *Builder) cutDownBlock() {

	if len(b.curBlock.entryOffsets) == 0 {
		return
	}
	b.append(y.U32SliceToBytes(b.curBlock.entryOffsets))
	b.append(y.U32ToBytes(uint32(len(b.curBlock.entryOffsets))))
	checksum := b.calculateCheckSum(b.curBlock.data[:b.curBlock.end])
	b.append(checksum)
	b.append(y.U32ToBytes(uint32(len(checksum))))
	b.blockList = append(b.blockList, b.curBlock)
	b.uncompressedSize.Add(uint32(b.curBlock.end))
	b.lenOffsets += uint32(int(math.Ceil(float64(len(b.curBlock.baseKey))/4))*4) + 40
	if b.blockChan != nil {
		b.blockChan <- b.curBlock
	}
}

// 将key-value添加到block中
// block会根据当前大小进行判断，决定是否要新增一个block
func (b *Builder) Add(key []byte, value y.ValueStruct, valueLen uint32) {
	b.addInternal(key, value, valueLen, false)
}

// 判断curBlock加了key-value后是否超过了设置block的大小4MB，如果超过则改block就完成了
func (b *Builder) checkBlockSizeGreaterThreshold(key []byte, value y.ValueStruct) bool {
	if len(b.curBlock.entryOffsets) <= 0 {
		return false
	}
	y.AssertTrue((uint32(len(b.curBlock.entryOffsets))+1)*4+4+8+4 < math.MaxUint32)
	entriesOffsetsSize := uint32((len(b.curBlock.entryOffsets)+1)*4 + 4 + 8 + 4)
	estimatedSize := uint32(b.curBlock.end) + 6 + uint32(len(key)) + value.EncodedSize() + entriesOffsetsSize
	if b.shouldEncrypt() {
		estimatedSize += aes.BlockSize
	}
	y.AssertTrue(uint64(b.curBlock.end)+uint64(estimatedSize) < math.MaxUint32)
	return estimatedSize > uint32(b.opts.BlockSize)
}

func (b *Builder) AddStaleKey(key []byte, v y.ValueStruct, valueLen uint32) {

}

func (b *Builder) keyDiff(newKey []byte) []byte {
	var i int
	//遍历新key,返回不等于curBlock中的baseKey的key
	//从第i起就不同
	for i = 0; i < len(newKey) && i < len(b.curBlock.baseKey); i++ {
		if newKey[i] != b.curBlock.baseKey[i] {
			break
		}
	}
	return newKey[i:]
}

// buildData包含bblock的数据，包含索引，包含校验合，包含文件尺寸大小以及分配
type buildData struct {
	blockList []*bblock
	index     []byte
	checksum  []byte
	Size      int
	alloc     *z.Allocator
}

func (bd *buildData) Copy(dst []byte) int {
	var written int
	for _, bl := range bd.blockList {
		written += copy(dst[written:], bl.data[:bl.end])
	}
	written += copy(dst[written:], bd.index)
	written += copy(dst[written:], y.U32ToBytes(uint32(len(bd.index))))

	written += copy(dst[written:], bd.checksum)
	written += copy(dst[written:], y.U32ToBytes(uint32(len(bd.checksum))))
	return written
}

// TODO: 表明key-value存储到builder的数据结构

func (b *Builder) addHelper(key []byte, v y.ValueStruct, vpLen uint32) {

	//将key-value的key进行murmurhash得到key hash写到builder中
	// keyHashes存储builder所有被hash过的key,包含所有bblock的hash后的key
	// 5个keyHashes就有5个值.每一个key都用uint32来表示他的hash值
	b.keyHashes = append(b.keyHashes, y.Hash(y.ParseKey(key)))

	// maxVersion只有最新的一个,是该builder中所有bblock的key中最大的那个version
	if version := y.ParseTs(key); version > b.maxVersion {
		b.maxVersion = version
	}
	var diffKey []byte
	if len(b.curBlock.baseKey) == 0 {
		// 将key追加到当前block的baseKey中
		// 如果curBlock的baseKey为0表明所有的key都是不同的key
		b.curBlock.baseKey = append(b.curBlock.baseKey[:], key...)
		diffKey = key
	} else {
		//diffKey表示为与curBlock中相比新的key
		diffKey = b.keyDiff(key)
	}

	y.AssertTrue(len(key)-len(diffKey) < math.MaxUint16)
	y.AssertTrue(len(diffKey) <= math.MaxUint16)
	//diff为新key
	//overlap表示不是新的key的数量
	h := header{
		overlap: uint16(len(key) - len(diffKey)),
		diff:    uint16(len(diffKey)),
	}
	b.curBlock.entryOffsets = append(b.curBlock.entryOffsets, uint32(b.curBlock.end))

	//对header进行编码
	b.append(h.Encode())
	//builder中先是header然后是不同的key
	b.append(diffKey)

	//申请空间
	dst := b.allocate(int(v.EncodedSize()))
	//将value写到目的地址中
	v.Encode(dst)
	b.onDiskSize += vpLen
}

// 使用4个字节来存储header的信息,前两个为overlap,后两个为diff
func (h header) Encode() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b[:2], h.overlap)
	binary.BigEndian.PutUint16(b[2:], h.diff)
	return b[:]
}

type header struct {
	overlap uint16
	diff    uint16
}

// table中存储的key和value是具体的值而不是地址
func (b *Builder) addInternal(key []byte, value y.ValueStruct, valueLen uint32, isStale bool) {
	if b.checkBlockSizeGreaterThreshold(key, value) {
		if isStale {
			b.staleDataSize += len(key) + 4 + 4
		}
		b.cutDownBlock()
		b.curBlock = &bblock{
			data: b.alloc.Allocate(b.opts.BlockSize + padding),
		}
	}
	b.addHelper(key, value, valueLen)
}

func (b *Builder) Empty() bool {
	return len(b.keyHashes) == 0
}
func (b *Builder) CutDoneBuildData() buildData {
	b.cutDownBlock()
	if b.blockChan != nil {
		close(b.blockChan)
	}

	b.wg.Wait()
	if len(b.blockList) == 0 {
		return buildData{}
	}

	bd := buildData{
		blockList: b.blockList,
		alloc:     b.alloc,
	}
	var f y.Filter
	//如果误报率大于0，表明开启了布隆过滤器,
	if b.opts.BloomFalsePositive > 0 {
		// 调用cutDone方法时,此时已经知道block中的key的数量了,每个block会有不同的key数量
		// 再结合误报率最终算法需要多个位的哈希函数.
		bits := y.BloomBitsPerKey(len(b.keyHashes), b.opts.BloomFalsePositive)
		f = y.NewFilter(b.keyHashes, bits)
	}
	//根据布隆过滤器来构建索引
	index, dataSize := b.buildIndex(f)

	var err error
	if b.shouldEncrypt() {
		index, err = b.encrypt(index)
		y.Check(err)
	}

	checksum := b.calculateCheckSum(index)
	bd.index = index
	bd.checksum = checksum
	bd.Size = int(dataSize) + len(index) + len(checksum) + 4 + 4

	return bd
}

func (b *Builder) writeBlockOffsets(builder *fbs.Builder) ([]fbs.UOffsetT, uint32) {
	var startOffset uint32
	var uOffs []fbs.UOffsetT

	// 将每一个bblock通过flatbuffer进行序列化.
	for _, bl := range b.blockList {
		uOff := b.writeBlockOffset(builder, bl, startOffset)
		uOffs = append(uOffs, uOff)
		startOffset += uint32(bl.end)
	}
	return uOffs, startOffset
}

func (b *Builder) writeBlockOffset(builder *fbs.Builder, bl *bblock, startOffset uint32) fbs.UOffsetT {
	k := builder.CreateByteVector(bl.baseKey)

	fb.BlockOffsetStart(builder)
	fb.BlockOffsetAddKey(builder, k)
	fb.BlockOffsetAddOffset(builder, startOffset)
	fb.BlockOffsetAddLen(builder, uint32(bl.end))
	return fb.BlockOffsetEnd(builder)
}

// 构建索引
func (b *Builder) buildIndex(bloom []byte) ([]byte, uint32) {
	builder := fbs.NewBuilder(3 << 20) //创建3MB
	boList, dataSize := b.writeBlockOffsets(builder)
	//boList为builder中bblock的个数
	fb.TableIndexStartOffsetsVector(builder, len(boList))

	for i := len(boList) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(boList[i])
	}
	boEnd := builder.EndVector(len(boList))
	var bfoff fbs.UOffsetT
	if len(bloom) > 0 {
		bfoff = builder.CreateByteVector(bloom)
	}
	fb.TableIndexStart(builder)
	fb.TableIndexAddOffsets(builder, boEnd)
	fb.TableIndexAddBloomFilter(builder, bfoff)
	fb.TableIndexAddMaxVersion(builder, b.maxVersion)
	fb.TableIndexAddUncompressedSize(builder, b.uncompressedSize.Load())
	fb.TableIndexAddKeyCount(builder, uint32(len(b.keyHashes)))
	fb.TableIndexAddOnDiskSize(builder, b.onDiskSize)
	fb.TableIndexAddStaleDataSize(builder, uint32(b.staleDataSize))
	builder.Finish(fb.TableIndexEnd(builder))

	buf := builder.FinishedBytes()
	index := fb.GetRootAsTableIndex(buf, 0)
	y.AssertTrue(index.MutateOnDiskSize(index.OnDiskSize() + uint32(len(buf))))
	return buf, dataSize
}

// 扩容，对bblock data字段进行扩容
// 扩容的原则是最大只有1GB
// 扩容的原则: 扩大两倍，但是如果请求的大小比原来的大，那么就去申请的大小+原有的大小
// 返回的字节就是能用的字节
func (b *Builder) allocate(need int) []byte {
	bb := b.curBlock
	//如果bblock的数据类型小于需要的类型，将sz设置成data长度的两倍
	if len(bb.data[bb.end:]) < need {
		sz := 2 * len(bb.data)
		if sz > (1 << 30) {
			sz = 1 << 30 // 1GB
		}
		//如果bb的最后一位加上需要的值大于sz即need超过end的长度那么取较大的值
		sz = max(sz, bb.end+need)
		//分配空间
		tmp := b.alloc.Allocate(sz)
		//将bblock的数据存储到tmp中
		//然后将tmp的内存复制到bblock data中
		//go可以通过unsafe拿到dst的内存地址，然后底层的汇编函数实现复制
		//实现原理： src/runtime/memmove_arm.s
		/* example code
		func main() {
			dst :=make([]int,10)
			soure:=[]int{1,2,3}

			copy(dst[2:2+len(soure)],soure)
			for a,b:=range dst{
				fmt.Printf("%d,:%d\n",a,b)
			}
		}

		*/
		copy(tmp, bb.data)
		bb.data = tmp
	}
	//更新下bb end的值
	bb.end += need
	return bb.data[bb.end-need : bb.end]
}

// 将新的字节追加到bblock的中
func (b *Builder) append(data []byte) {
	//dst为bblock的当前游标
	dst := b.allocate(len(data))
	//将[]byte添加到bblock中
	//此时是将data复制到dst，而dst是在bblock的第i和第j位因此可以成功复制，
	//猜测:将data写到字节数据的地址中.
	copyLen := copy(dst, data)
	y.AssertTrue(len(data) == copyLen)
	//	疑问: 使用copy能否将data复制到builder 的bblock中?
	//	解答，可以的,使用runtime中的汇编实现复制操作
}

func (b *Builder) DataKey() *pb.DataKey {
	return nil
}

func (b *Builder) shouldEncrypt() bool {
	return b.opts.DataKey != nil
}

func (b *Builder) encrypt(data []byte) ([]byte, error) {
	iv, err := y.GenerateIV()
	if err != nil {
		return data, y.Wrapf(err, "Error while generating IV in Builder.encrypt")
	}
	needSz := len(data) + len(iv)
	dst := b.alloc.Allocate(needSz)

	if err = y.XORBlock(dst[:len(data)], data, b.DataKey().Data, iv); err != nil {
		return data, y.Wrapf(err, "Error while encrypting in Builder.encrypt")
	}
	y.AssertTrue(len(iv) == copy(dst[len(data):], iv))
	return dst, nil
}

// 对数据进行压缩
func (b *Builder) compressData(data []byte) ([]byte, error) {
	switch b.opts.Compression {
	case options.None:
		return data, nil
	case options.Snappy:
		sz := s2.MaxEncodedLen(len(data))
		//分配空间
		dst := b.alloc.Allocate(sz)
		return s2.EncodeSnappy(dst, data), nil
	case options.ZSTD:
		sz := y.ZSTDCompressBound(len(data))
		dst := b.alloc.Allocate(sz)
		return y.ZSTDCompress(dst, data, b.opts.ZSTDCompressionLevel)
	}
	return nil, errors.New("Unsupported compression type")
}

// 对需要的bblock进行压缩和加密
func (b *Builder) handleBlock() {
	defer b.wg.Done()
	doCompress := b.opts.Compression != options.None

	//遍历bblock channel
	for item := range b.blockChan {
		blockBuf := item.data[:item.end]
		if doCompress {
			//将压缩完的数据写到out中
			out, err := b.compressData(blockBuf)
			y.Check(err)
			blockBuf = out
		}
		if b.shouldEncrypt() {
			out, err := b.encrypt(blockBuf)
			y.Check(y.Wrapf(err, "Error while encrying block in table builder."))
			blockBuf = out
		}
		allocatedSpace := maxEncodedLen(b.opts.Compression, item.end) + padding + 1
		y.AssertTrue(len(blockBuf) <= allocatedSpace)

		//将加密好的数据写到data中
		item.data = blockBuf
		item.end = len(blockBuf)
		b.compressedSize.Add(uint32(len(blockBuf)))
	}
}

func maxEncodedLen(ctype options.CompressionType, sz int) int {
	switch ctype {
	case options.Snappy:
		return s2.MaxEncodedLen(sz)
	case options.ZSTD:
		y.ZSTDCompressBound(sz)
	}
	return sz
}

const (
	KB      = 1 << 10
	MB      = KB << 10
	padding = 256
)
const maxAllocatorInitialSz = 256 << 20 // 256MB

func NewTableBuilder(opts Options) *Builder {
	//default opts.TableSize : 2MB
	// 申请的空间是目标table的两倍即4MB,如果超过了最大的256MB,则以256为准
	sz := 2 * int(opts.TableSize)
	if sz > maxAllocatorInitialSz {
		sz = maxAllocatorInitialSz
	}
	b := &Builder{
		// Allocator 通过大块分配内存来分摊小块分配的成本。
		//它在内部使用 z.Calloc 来分配内存。一旦分配完毕，内存就不会移动，
		//因此可以安全地使用已分配的字节，以不安全的方式将其转换为 Go 结构指针。
		//维护自由列表的速度很慢。因此，Allocator 只分配内存，最后释放整个 Allocator。
		// get实则是创建这样的内存区域
		// 这样做的目的是什么?
		alloc: opts.AllocPool.Get(sz, "TableBuilder"),
		opts:  &opts,
	}
	b.alloc.Tag = "Builder"
	b.curBlock = &bblock{
		data: b.alloc.Allocate(opts.BlockSize + padding),
	}
	b.opts.tableCapacity = uint64(float64(b.opts.TableSize) * 0.95)
	if b.opts.Compression == options.None && b.opts.DataKey == nil {
		return b
	}
	count := 2 * runtime.NumCPU()
	b.blockChan = make(chan *bblock, count*2)
	b.wg.Add(count)
	for i := 0; i < count; i++ {
		go b.handleBlock()
	}
	return b
}

func (b *Builder) ReachedCapacity() bool {

	return false
}

// CutDownBuilder finishes the table by appending the index.
/*
The table structure looks like
+---------+------------+-----------+---------------+
| Block 1 | Block 2    | Block 3   | Block 4       |
+---------+------------+-----------+---------------+
| Block 5 | Block 6    | Block ... | Block N       |
+---------+------------+-----------+---------------+
| Index   | Index Size | Checksum  | Checksum Size |
+---------+------------+-----------+---------------+
*/

// 将bblock中的内存都更新成字节数组
func (b *Builder) CutDownBuilder() []byte {
	bd := b.CutDoneBuildData()
	buf := make([]byte, bd.Size)
	written := bd.Copy(buf)
	y.AssertTrue(written == len(buf))
	return buf
}

func (b *Builder) Opts() *Options {
	return nil
}

// TODO: 含义是什么?
func (b *Builder) Close() {
	b.opts.AllocPool.Return(b.alloc)
}

// TODO: bblock的作用是什么？
type bblock struct {
	data         []byte
	baseKey      []byte // TODO: baseKey的含义是什么?
	entryOffsets []uint32
	end          int
}
