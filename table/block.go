package table

import (
	"google.golang.org/protobuf/proto"
	"sync/atomic"
	"unsafe"
	"wiscdb/pb"
	"wiscdb/y"
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
	cs := &pb.CheckSum{}
	// 将b的checksum 反序列到cs中
	if err := proto.Unmarshal(b.checksum, cs); err != nil {
		return y.Wrapf(err, "unable to unmarshal checksum for block")
	}
	return y.VerifyCheckSum(b.data, cs)
}
