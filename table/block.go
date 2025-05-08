package table

import (
	"sync/atomic"
	"unsafe"
)

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
