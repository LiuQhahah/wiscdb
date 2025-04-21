package table

import (
	"github.com/dgraph-io/ristretto/v2/z"
	fbs "github.com/google/flatbuffers/go"
	"sync"
	"sync/atomic"
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

	return nil
}

func (b *Builder) finishBlock() {

}

func (b *Builder) Add(key []byte, value y.ValueStruct, valueLen uint32) {

}

func (b *Builder) shouldFinishBlock(key []byte, value y.ValueStruct) bool {
	return false
}

func (b *Builder) AddStaleKey(key []byte, v y.ValueStruct, valueLen uint32) {

}

func (b *Builder) keyDiff(newKey []byte) []byte {
	return nil
}

type buildData struct {
	blockList []*bblock
	index     []byte
	checksum  []byte
	Size      int
	alloc     *z.Allocator
}

func (bd *buildData) Copy(dst []byte) int {
	var written int

	return written
}

func (b *Builder) addHelper(key []byte, v y.ValueStruct, vpLen uint32) {

}

func (b *Builder) addInternal(key []byte, value y.ValueStruct, valueLen uint32, isStale bool) {

}

func (b *Builder) Empty() bool {
	return false
}
func (b *Builder) Done() buildData {
	return buildData{}
}

func (b *Builder) writeBlockOffsets(builder *fbs.Builder) ([]fbs.UOffsetT, uint32) {
	return nil, 0
}

func (b *Builder) writeBlockOffset(builder *fbs.Builder, bl *bblock, startOffset uint32) fbs.UOffsetT {
	return 0
}

func (b *Builder) buildIndex(bloom []byte) ([]byte, uint32) {
	return nil, 0
}

func (b *Builder) allocate(need int) []byte {
	return nil
}

func (b *Builder) append(data []byte) {

}

func (b *Builder) DataKey() *pb.DataKey {
	return nil
}

func (b *Builder) shouldEncrypt() bool {
	return false
}

func (b *Builder) encrypt(data []byte) ([]byte, error) {
	return nil, nil
}

func (b *Builder) compressData(data []byte) ([]byte, error) {
	return nil, nil
}

func (b *Builder) handleBlock() {

}

func NewTableBuilder(opts Options) *Builder {

	return &Builder{}
}

func (b *Builder) ReachedCapacity() bool {

	return false
}

func (b *Builder) Finish() []byte {
	return nil
}

func (b *Builder) Opts() *Options {
	return nil
}

type bblock struct {
	data         []byte
	baseKey      []byte
	entryOffsets []uint32
	end          int
}
