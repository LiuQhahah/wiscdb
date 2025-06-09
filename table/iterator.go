package table

import (
	"bytes"
	"io"
	"sort"
	"time"

	"wiscdb/fb"
	"wiscdb/y"
)

var (
	REVERSED int = 2
	NOCACHE  int = 4
)

type Iterator struct {
	t    *Table
	bpos int
	bi   blockIterator
	err  error
	opt  int
}

func (itr *Iterator) Next() {

	if itr.opt&REVERSED == 0 {
		itr.Next()
	} else {
		itr.prev()
	}
}

func (itr *Iterator) ReWind() {

	if itr.opt&REVERSED == 0 {
		itr.seekToFirst()
	} else {
		itr.seekToLast()
	}
}

func NewConcatIterator(tbls []*Table, opt int) *concatIterator {
	iters := make([]*Iterator, len(tbls))
	for i := 0; i < len(tbls); i++ {
		tbls[i].IncrRef()
	}
	return &concatIterator{
		options: opt,
		iters:   iters,
		tables:  tbls,
		idx:     -1,
	}
}
func (itr *Iterator) Value() (ret y.ValueStruct) {
	ret.Decode(itr.bi.val)
	return
}

type blockIterator struct {
	data         []byte
	idx          int
	err          error
	baseKey      []byte
	key          []byte
	val          []byte
	entryOffsets []uint32
	block        *Block
	tableID      uint64
	blockID      int
	prevOverlap  uint16
}

func (itr *Iterator) Rewind() {
	if itr.opt&REVERSED == 0 {
		itr.seekToFirst()
	} else {
		itr.seekToLast()
	}
}

// 迭代器中寻找第一个Key
func (itr *Iterator) seekToFirst() {
	numBlocks := itr.t.offsetsLength()
	if numBlocks == 0 {
		itr.err = io.EOF
		return
	}
	itr.bpos = 0
	block, err := itr.t.block(itr.bpos, itr.useCache())
	if err != nil {
		itr.err = err
		return
	}
	itr.bi.tableID = itr.t.id
	itr.bi.blockID = itr.bpos
	itr.bi.setBlock(block)
	itr.bi.seekToFirst()
	itr.err = itr.bi.Error()
}

func (itr *Iterator) useCache() bool {
	return false
}

// 迭代器选招最后一个key
func (itr *Iterator) seekToLast() {

}

var (
	origin  = 0
	current = 1
)

func (itr *Iterator) reset() {
	itr.bpos = 0
	itr.err = nil
}
func (itr *Iterator) seekFrom(key []byte, whence int) {
	itr.err = nil
	switch whence {
	case origin:
		itr.reset()
	case current:

	}
	var ko fb.BlockOffset
	idx := sort.Search(itr.t.offsetsLength(), func(idx int) bool {
		y.AssertTrue(itr.t.offsets(&ko, idx))
		return y.CompareKeys(ko.KeyBytes(), key) > 0
	})

	if idx == 0 {
		itr.seekHelper(0, key)
		return
	}
	itr.seekHelper(idx-1, key)
	if itr.err == io.EOF {
		if idx == itr.t.offsetsLength() {
			return
		}
		itr.seekHelper(idx, key)
	}
}
func (itr *Iterator) seekHelper(blockIdx int, key []byte) {

}
func (itr *Iterator) seek(key []byte) {
	itr.seekFrom(key, origin)
}
func (itr *Iterator) Seek(key []byte) {
	if itr.opt&REVERSED == 0 {
		itr.seek(key)
	} else {
		itr.seekForPrev(key)
	}
}

func (itr *Iterator) prev() {
	itr.err = nil
	if itr.bpos < 0 {
		itr.err = io.EOF
		return
	}
	if len(itr.bi.data) == 0 {
		block, err := itr.t.block(itr.bpos, itr.useCache())
		if err != nil {
			itr.err = err
			return
		}
		itr.bi.tableID = itr.t.id
		itr.bi.blockID = itr.bpos
		itr.bi.setBlock(block)
		itr.bi.seekToLast()
		itr.err = itr.bi.Error()
		return
	}
	itr.bi.prev()

	if !itr.bi.Valid() {
		itr.bpos--
		itr.bi.data = nil
		itr.prev()
		return
	}
}

func (itr *Iterator) seekForPrev(key []byte) {
	itr.seekFrom(key, origin)
	if !bytes.Equal(itr.Key(), key) {
		itr.prev()
	}
}

func (itr *Iterator) ValueCopy() (ret y.ValueStruct) {
	dst := y.Copy(itr.bi.val)
	ret.Decode(dst)
	return
}

func (itr *blockIterator) Valid() bool {
	return itr != nil && itr.err == nil
}
func (itr *blockIterator) prev() {

}

func (itr *blockIterator) setBlock(b *Block) {
	itr.block.decrRef()
	itr.block = b
	itr.idx = 0
	itr.baseKey = itr.baseKey[:0]
	itr.prevOverlap = 0
	itr.key = itr.key[:0]
	itr.val = itr.val[:0]
	itr.data = b.data[:b.entriesIndexStart]
	itr.entryOffsets = b.entryOffsets
}

func (itr *blockIterator) Error() error {
	return nil
}
func (itr *blockIterator) seekToFirst() {

}
func (itr *blockIterator) seekToLast() {

}
func (itr *blockIterator) Close() {
	itr.block.decrRef()
}
func (itr *Iterator) Close() error {
	itr.bi.Close()
	return itr.t.DecrRef()
}
func (itr *Iterator) Valid() bool {
	return itr.err == nil
}

// 拿到block iterator中的key作为该table中最大的key
func (itr *Iterator) Key() []byte {
	return itr.bi.key
}

const BitDelete byte = 1 << 0

func IsDeletedOrExpired(meta byte, expiresAt uint64) bool {
	if meta&BitDelete > 0 {
		return true
	}
	if expiresAt == 0 {
		return false
	}
	return expiresAt <= uint64(time.Now().Unix())
}
