package table

import "io"

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
