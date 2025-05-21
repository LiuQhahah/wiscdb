package skl

import (
	"github.com/dgraph-io/ristretto/v2/z"
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

// +---------+---------+---------+---------+
// | value   | keyOffset | keySize | height | tower...
// +---------+---------+---------+---------+
// value 占用 8 字节(uint64)
// keyOffset 占用 4 字节(uint32)
// keySize 占用 2 字节(uint16)
// height 占用 2 字节(uint16)
// tower 占用 72 字节([18]uint32，每个 uint32 占用 4 字节)
// node一共有88个字节.
type node struct {
	value     atomic.Uint64     //存储value的地址
	keyOffset uint32            //存储key的偏移量
	keySize   uint16            //存储key的大小
	height    uint16            //存储高度
	tower     [18]atomic.Uint32 // tower的含义: 存储偏移量
}

// 内存分配中的分配区概念
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
	//  Store 等同于Set操作,Load等同于Get操作
	// TODO: n为什么要设定为1?
	out.n.Store(1)
	return out
}

// 这个函数的主要作用是在一个预分配的内存缓冲区中分配一块内存，并返回该内存块的偏移量
// 它考虑了内存对齐和未使用的空间，确保分配的内存块是有效的。
// 如果内存不足，函数会触发 panic。
// 该函数通常用于实现自定义的内存分配器，特别是在需要高效管理内存的场景中（如数据库、缓存系统等）。

// 在 Arena 内存池中分配一个跳表节点（node），并返回其 对齐后的内存偏移量
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
	//返回value的偏移值
	return m
}

func (s *Arena) putKey(key []byte) uint32 {
	l := uint32(len(key))
	n := s.n.Add(l)
	m := n - l
	copy(s.buf[m:n], key)
	return m
}

// 根据偏移量返回node struct
// 根据 偏移量使用unsafe在内存中解析出node struct
func (s *Arena) getNode(offset uint32) *node {
	if offset == 0 {
		return nil
	}
	v := &s.buf[offset]
	//如果 offset 指向的内存区域与 node 结构体的内存布局一致，
	//那么 (*node)(unsafe.Pointer(v)) 就可以正确地将这块内存解释为一个 node 结构体
	// v仅仅是内存位置, 编译器会根据node的结构来访问对应内存大小，然后将这段内存序列化成node struct.
	//内存的大小是由 目标类型 决定的。当你将 unsafe.Pointer 转换为一个具体的指针类型（例如 *node）时，
	//Go 编译器会根据目标类型的大小来访问内存

	//node 结构体的大小是由其字段的类型和对齐规则决定的
	//node存储在分配区Arena中，通过unsafe 包进行直接转化
	return (*node)(unsafe.Pointer(v))
}

/*return key */
//根据偏移量以及key的长度返回字节数组
func (s *Arena) getKey(offset uint32, size uint16) []byte {
	return s.buf[offset : offset+uint32(size)]
}

func (s *Arena) getVal(offset uint32, size uint32) (ret y.ValueStruct) {
	val := s.buf[offset : offset+size]
	v := y.ValueStruct{}
	v.Decode(val)
	return v
}

// 获取node的偏移量
func (s *Arena) getNodeOffset(nd *node) uint32 {
	if nd == nil {
		return 0
	}
	//强制将nd转化成无符号整型指针,强制将s.buf[0]转化成无符号指针
	return uint32(uintptr(unsafe.Pointer(nd)) - uintptr(unsafe.Pointer(&s.buf[0])))
}

func newNode(arena *Arena, key []byte, v y.ValueStruct, height int) *node {
	offset := arena.putNode(height)
	node := arena.getNode(offset)
	node.keyOffset = arena.putKey(key)
	node.keySize = uint16(len(key))
	node.height = uint16(height)
	node.value.Store(encodeValue(arena.putVal(v), v.EncodedSize()))
	return node
}

// save value address
// 使用4个字节存储value的起始位置以及value的尺寸,
// offset: 0x01234567, size:0xABCDEFAB
// return value: 0x01234567ABCDEFAB
func encodeValue(valOffset uint32, valSize uint32) uint64 {
	return uint64(valOffset)<<32 | uint64(valSize)
}

// 与encodeVale反过来.
// example:0x01234567ABCDEFAB
// decode的return值
// valOffset:0xABCDEFAB
// valSize: 0x01234567
// value存储的是value的偏移量地址以及value大小的值
func decodeValue(value uint64) (valOffset uint32, valSize uint32) {
	valSize = uint32(value)
	valOffset = uint32(value >> 32)
	return
}

/*
*
arenaSize the size of skip list. default value: 64MB.
*/
func NewSkipList(arenaSize int64) *SkipList {
	arena := newArena(arenaSize)
	head := newNode(arena, nil, y.ValueStruct{}, maxHeight)
	s := &SkipList{head: head, arena: arena}
	s.height.Store(1)
	s.ref.Store(1)
	return s
}

// set value into node.
// 先将value存储到内存arena，同时返回当前的游标
func (s *node) setValue(arena *Arena, v y.ValueStruct) {
	valOffset := arena.putVal(v)
	//value存储的是当前value的偏移量及其长度，使用8个字节表示，也是一种编码方式
	value := encodeValue(valOffset, v.EncodedSize())
	s.value.Store(value)
}

// h指的是链表的高度，
// 根据高度得到当前层次的偏移量
// 偏移量存储在node的tower中
func (s *node) getNextOffset(h int) uint32 {
	return s.tower[h].Load()
}

func (s *node) getValueOffset() (uint32, uint32) {
	value := s.value.Load()
	return decodeValue(value)
}
func (s *node) casNextOffset(h int, old, val uint32) bool {
	return s.tower[h].CompareAndSwap(old, val)
}

// 获取node总分配区的key，返回值是个字节数组
func (s *node) key(arena *Arena) []byte {
	//根据key的起始位置以及key的大小返回对应的字节数
	//s为当前链表的游标位置,每一次插入时都会更新当前的key的偏移量以及key的长度，
	//从而方便解码出key的值
	return arena.getKey(s.keyOffset, s.keySize)
}

// 随机生成高度，加高一层的概率是1/3
func (s *SkipList) randomHeight() int {
	h := 1
	for h < maxHeight && z.FastRand() <= heightIncrease {
		h++
	}
	return h
}

// 返回指定层高链表的下一个节点
func (s *SkipList) getNext(nd *node, height int) *node {
	return s.arena.getNode(nd.getNextOffset(height))
}

// 加载跳表的高度
func (s *SkipList) getHeight() int32 {
	return s.height.Load()
}

// 在跳表中查找key
func (s *SkipList) findNear(key []byte, less bool, allowEqual bool) (*node, bool) {
	//第一层的第一个链表
	x := s.head
	//从最高层开始遍历
	level := int(s.getHeight() - 1)
	for {
		//返回的是一个链表.
		next := s.getNext(x, level)
		if next == nil { //表明已经遍历到最后一位了
			if level > 0 {
				level--
				continue
			}
			if !less {
				return nil, false
			}
			if x == s.head {
				return nil, false
			}
			return x, false
		}

		//将内存块作为参数传入进入从而获取key
		nextKey := next.key(s.arena)
		cmp := y.CompareKeys(key, nextKey)
		if cmp > 0 {
			x = next
			continue
		}
		if cmp == 0 {
			if allowEqual {
				return next, true
			}
			if !less {
				return s.getNext(next, 0), false
			}
			if level > 0 {
				level--
				continue
			}
			//此时level会等于0,那么就需要重置跳表的头
			if x == s.head {
				return nil, false
			}
			return x, false
		}
		if level > 0 {
			level--
			continue
		}
		if !less {
			return next, false
		}
		if x == s.head {
			return nil, false
		}
		return x, false
	}
}

func (s *SkipList) MemSize() int64 {
	return s.arena.size()
}

// 根据key找到这一层插入的位置.
func (s *SkipList) findSpliceForLevel(key []byte, before *node, level int) (*node, *node) {
	//遍历当前层的链表
	for {
		//得到当前层的下一个node
		next := s.getNext(before, level)
		if next == nil {
			return before, next
		}
		//得到next存储的key是什么
		nextKey := next.key(s.arena)
		//比如两个key分别是i，k中间是j
		cmp := y.CompareKeys(key, nextKey)
		if cmp == 0 { //如果两个key相等,表明找打这层是有这个key的
			return next, next
		}
		if cmp < 0 { //如果key小于下一个key
			return before, next // 则插入在before中next之间
		}
		before = next //更新下before位置
	}
}

// 链表中Put操作
func (s *SkipList) Put(key []byte, value y.ValueStruct) {
	listHeight := s.getHeight() //获取当前跳表的高度
	var prev [maxHeight + 1]*node
	var next [maxHeight + 1]*node
	prev[listHeight] = s.head // 跳表的头为prev指针，设置最高层的node的节点
	next[listHeight] = nil
	//从最高层开始遍历
	for i := int(listHeight) - 1; i >= 0; i-- { //遍历当前跳表的高度
		//第一次传入的参数是跳表的头，返回链表头与下一次节点的值
		//一层一层的查找key
		prev[i], next[i] = s.findSpliceForLevel(key, prev[i+1], i)
		//这是tricky的设置，找到该层的key后，会返回两个next
		if prev[i] == next[i] {
			//会直接返回，此时下一层的对应的key并没有更新，这样也是对的，不需要每一层都更新
			//原则就是最新找到的值是最新的，妙计
			prev[i].setValue(s.arena, value)
			return
		}
	}
	//如果该key不存在与链表中，即最底层也不存在，则走下面的逻辑

	height := s.randomHeight()
	//利用提供的key，value，height以及分配区创建一个新的node.
	node1 := newNode(s.arena, key, value, height)
	//如果随机生成的高度大于当前链表的高度为真
	for height > int(listHeight) {
		//更新下当前跳表的属性：高度
		//线程安全情况下更新跳表的高度
		//CAS 失败：其他线程可能已修改高度，需重新获取 listHeight 并重试
		if s.height.CompareAndSwap(listHeight, int32(height)) {
			break
		}
		//重新获取跳表高度
		listHeight = s.getHeight()
	}
	//如果当前跳表中没有改key，则在每一层都要插入改key？
	for i := 0; i < height; i++ {
		for {
			if prev[i] == nil {
				//不做最底层
				y.AssertTrue(i > 1)
				//在高度为i中找到可以插入的地方
				prev[i], next[i] = s.findSpliceForLevel(key, s.head, i)
				y.AssertTrue(prev[i] != next[i])
			}

			//传入node来得到node在分配区arena的偏移量
			nextOffset := s.arena.getNodeOffset(next[i])
			//将偏移量存储在tower中
			node1.tower[i].Store(nextOffset)
			// 插入新节点后要更新prev[i]的值
			//需要将 prev[i] 的第 i 层指针从原来的 nextOffset（指向 next[i]）改为指向 node1。
			//同时，node1 的第 i 层指针应指向 next[i]
			//原子地将 prev[i] 的第 i 层指针从 nextOffset 改为指向 node1。
			//CAS 失败处理 原因：其他线程可能已修改 prev[i] 的指针
			//恢复：重新调用 findSpliceForLevel 获取最新的 prev[i] 和 next[i]
			if prev[i].casNextOffset(i, nextOffset, s.arena.getNodeOffset(node1)) {
				break
			}
			//恢复：重新调用 findSpliceForLevel 获取最新的 prev[i] 和 next[i]
			prev[i], next[i] = s.findSpliceForLevel(key, prev[i], i)
			if prev[i] == next[i] {
				//只能发生在第0层
				y.AssertTruef(i == 0, "Euality can happend only on base level: %d", i)
				prev[i].setValue(s.arena, value)
				return
			}
		}
	}
}

func (s *SkipList) Empty() bool {
	return s.findLast() == nil
}

// 返回最后一个元素,从跳表中返回node
func (s *SkipList) findLast() *node {
	n := s.head
	//从最高层遍历
	level := int(s.getHeight()) - 1
	for {
		//根据层次找到下一个node，然后持续调用Next
		next := s.getNext(n, level)
		if next != nil {
			n = next
			continue
		}
		//当level为0时并且就是跳表中第0层最后一个node
		if level == 0 {
			//如果此时的node是跳表头表示没有找到last node，否则改node 就是最后一个node
			if n == s.head {
				return nil
			}
			return n
		}
		level--
	}
}

// 根据key 找到key对应的value struct
func (s *SkipList) Get(key []byte) y.ValueStruct {
	sNode, _ := s.findNear(key, true, true)
	if sNode == nil {
		return y.ValueStruct{}
	}
	k := s.arena.getKey(sNode.keyOffset, sNode.keySize)
	if !y.SameKey(k, key) {
		return y.ValueStruct{}
	}
	valueOffset, valueSize := sNode.getValueOffset()
	vs := s.arena.getVal(valueOffset, valueSize)
	vs.Version = y.ParseTs(k)
	return vs
}

func (s *SkipList) IncrRef() {
	s.ref.Add(1)
}

// 回收资源
func (s *SkipList) DecrRef() {
	newRef := s.ref.Add(-1)
	if newRef > 0 {
		return
	}
	if s.OnClose != nil {
		s.OnClose()
	}
	s.arena = nil
	s.head = nil

}

// 用来遍历跳表, memtable,immutable的合并的操作都需要遍历
type Iterator struct {
	list *SkipList
	n    *node
}

func (i *Iterator) Close() error {
	i.list.DecrRef()
	return nil
}
func (i *Iterator) Valid() bool {
	return i.n != nil
}

// 获取当前游标的key
func (i *Iterator) Key() []byte {
	return i.list.arena.getKey(i.n.keyOffset, i.n.keySize)
}

func (i *Iterator) Value() y.ValueStruct {
	valOffset, valSize := i.n.getValueOffset()
	return i.list.arena.getVal(valOffset, valSize)
}

// 返回当前迭代器游标的value地址
func (i *Iterator) LoadValueAddress() uint64 {
	return i.n.value.Load()
}

func (i *Iterator) Next() {
	y.AssertTrue(i.Valid())
	i.n = i.list.getNext(i.n, 0)
}

func (i *Iterator) Prev() {
	y.AssertTrue(i.Valid())
	i.n, _ = i.list.findNear(i.Key(), true, false)
}

func (i *Iterator) Seek(target []byte) {
	i.n, _ = i.list.findNear(target, false, true)
}

func (i *Iterator) SeekForPrev(target []byte) {
	i.n, _ = i.list.findNear(target, true, true)
}

func (i *Iterator) SeekToFirst() {
	i.n = i.list.getNext(i.list.head, 0)
}

func (i *Iterator) SeekToLast() {
	i.n = i.list.findLast()
}

func (i *Iterator) NewIterator() *Iterator {
	return &Iterator{}
}
func (s *UniIterator) ReWind() {
	if !s.reversed {
		s.iter.SeekToFirst()
	} else {
		s.iter.SeekToLast()
	}
}

type UniIterator struct {
	iter     *Iterator
	reversed bool
}

// 包装了一层，利用参数reversed来数据是顺序还是逆序
func (s *SkipList) NewUniIterator(reversed bool) *UniIterator {
	return &UniIterator{
		iter:     s.NewIterator(),
		reversed: reversed,
	}
}

func (s *SkipList) NewIterator() *Iterator {
	s.IncrRef()
	return &Iterator{list: s}
}
func (s *UniIterator) Next() {
	if !s.reversed {
		s.iter.Next()
	} else {
		s.iter.Prev()
	}
}

func (s *UniIterator) Seek(key []byte) {
	if !s.reversed {
		s.iter.Seek(key)
	} else {
		s.iter.SeekForPrev(key)
	}
}
func (s *UniIterator) Key() []byte {
	return s.iter.Key()
}

func (s *UniIterator) Value() y.ValueStruct {
	return s.iter.Value()
}

func (s *UniIterator) Valid() bool {
	return s.iter.Valid()
}

func (s *UniIterator) Close() error {
	return s.iter.Close()
}
