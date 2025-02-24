package skl

import (
	"math"
	"sync/atomic"
	"unsafe"
	"wiscdb/y"
)

const (
	offsetSize = int(unsafe.Sizeof(uint32(0)))
	//对齐大小（例如 8 字节对齐）
	nodeAlign = int(unsafe.Sizeof(uint64(0))) - 1
)

type SkipList struct {
	height  atomic.Int32
	head    *node
	ref     atomic.Int32
	arena   *Arena
	OnClose func()
}

const (
	maxHeight      = 20
	heightIncrease = math.MaxUint32 / 3
	MaxNodeSize    = int(unsafe.Sizeof(node{}))
)

type node struct {
	value     atomic.Uint64
	ketOffset uint32
	keySize   uint16
	height    uint16
	tower     [18]atomic.Uint32
}

type Arena struct {
	n   atomic.Uint32
	buf []byte
}

// return the value of n
func (s *Arena) size() int64 {
	return int64(s.n.Load())
}
func newArena(n int64) *Arena {
	out := &Arena{
		buf: make([]byte, n),
	}
	out.n.Store(1)
	return out
}

// 这个函数的主要作用是在一个预分配的内存缓冲区中分配一块内存，并返回该内存块的偏移量
// 它考虑了内存对齐和未使用的空间，确保分配的内存块是有效的。
// 如果内存不足，函数会触发 panic。
// 该函数通常用于实现自定义的内存分配器，特别是在需要高效管理内存的场景中（如数据库、缓存系统等）。
func (s *Arena) putNode(height int) uint32 {
	unusedSize := (maxHeight - height) * offsetSize
	//本次分配的内存块的大小
	l := uint32(MaxNodeSize - unusedSize + nodeAlign)
	//当前内存分配器的偏移量（即下一个可用内存的起始位置）
	n := s.n.Add(l)
	//计算对齐后的内存偏移量。
	//它的目的是确保分配的内存块的起始地址是对齐的（通常是按照 nodeAlign 对齐）
	//n - l：计算出本次分配的内存块的起始位置（未对齐）
	//n - l + uint32(nodeAlign)：在未对齐的起始位置基础上，加上对齐大小。
	//这一步是为了确保后续的对齐操作能够正确进行
	//&^ uint32(nodeAlign)：这是一个按位清除操作（AND NOT），用于将低位的对齐掩码清零，从而得到对齐后的地址
	m := (n - l + uint32(nodeAlign)) &^ uint32(nodeAlign)
	//对齐后的内存偏移量
	//如果 m 为 0，意味着对齐后的内存偏移量是 0，即分配的内存块的起始地址位于 Arena 内存缓冲区的起始位置（s.buf 的开头）
	return m
}

// 将v写到v中
func (s *Arena) putVal(v y.ValueStruct) uint32 {
	l := v.EncodedSize()
	n := s.n.Add(l)
	m := n - l
	v.Encode(s.buf[m:])
	return 0
}

func (s *Arena) putKey(key []byte) uint32 {
	l := uint32(len(key))
	n := s.n.Add(l)
	m := n - l
	copy(s.buf[m:n], key)
	return m
}

func (s *Arena) getNode(offset uint32) *node {
	return &node{}
}

/*return key */
func (s *Arena) getKey(offset uint32, size uint16) []byte {
	return s.buf[offset : offset+uint32(size)]
}

func (s *Arena) getVal(offset uint32, size uint32) (ret y.ValueStruct) {
	val := s.buf[offset : offset+size]
	v := y.ValueStruct{}
	v.Decode(val)
	return v
}

func (s *Arena) getNodeOffset(nd *node) uint32 {
	return 0
}

func newNode(arena *Arena, key []byte, v y.ValueStruct, height int) *node {
	return nil
}
func encodeValue(valOffset uint32, valSize uint32) uint64 {
	return 0
}

func decodeValue(value uint64) (valOffset uint32, valSize uint32) {
	return 0, 0
}

func NewSkipList(arenaSize int64) *SkipList {
	return nil
}

func (s *node) setValue(arena *Arena, v y.ValueStruct) {

}
func (s *node) getNextOffset(h int) uint32 {
	return 0
}

func (s *node) casNextOffset(h int, old, val uint32) bool {
	return false
}

func (s *SkipList) randomHeight() int {
	return 0
}

func (s *SkipList) getNext(nd *node, height int) *node {
	return nil
}

func (s *SkipList) getHeight() int32 {
	return 0
}
func (s *SkipList) findNear(key []byte, less bool, allowEqual bool) (*node, bool) {
	return nil, false
}

func (s *SkipList) MemSize() int64 {
	return 0
}

func (s *SkipList) findSplitForLevel(key []byte, before *node, less int) (*node, *node) {
	return nil, nil
}

func (s *SkipList) Put(key []byte, v y.ValueStruct) {

}

func (s *SkipList) Empty() bool {
	return false
}

func (s *SkipList) findLast() *node {
	return nil
}

func (s *SkipList) Get(key []byte) y.ValueStruct {
	return y.ValueStruct{}
}

func (s *SkipList) IncrRef() {

}

func (s *SkipList) DecrRef() {

}

type Iterator struct {
	list *SkipList
	n    *node
}

func (i *Iterator) Close() error {
	return nil
}
func (i *Iterator) Valid() bool {
	return false
}

func (i *Iterator) Key() []byte {
	return nil
}

func (i *Iterator) Value() y.ValueStruct {
	return y.ValueStruct{}
}

func (i *Iterator) ValueUint64() uint64 {
	return 0
}

func (i *Iterator) Next() {

}

func (i *Iterator) Prev() {

}

func (i *Iterator) Seek(target []byte) {

}

func (i *Iterator) SeekForPrev(target []byte) {

}

func (i *Iterator) SeekToFirst() {

}

func (i *Iterator) SeekToLast() {

}

func (i *Iterator) NewIterator() *Iterator {
	return &Iterator{}
}

type UniIterator struct {
	iter     *Iterator
	reversed bool
}

func (s *SkipList) NewUniIterator(reversed bool) *UniIterator {
	return &UniIterator{}
}

func (s *UniIterator) Next() {

}

func (s *UniIterator) ReWind() {

}

func (s *UniIterator) Seek(key []byte) {

}
func (s *UniIterator) Key() []byte {
	return nil
}

func (s *UniIterator) Value() y.ValueStruct {
	return y.ValueStruct{}
}

func (s *UniIterator) Valid() bool {
	return false
}

func (s *UniIterator) Close() error {
	return nil
}
