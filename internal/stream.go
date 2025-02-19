package internal

import (
	"context"
	"github.com/dgraph-io/ristretto/v2/z"
	"sync/atomic"
	"wiscdb/pb"
)

const batchSize = 16 << 20

var maxStreamSize = uint64(100 << 20)

type Stream struct {
	Prefix       []byte
	NumGo        int
	LogPrefix    string
	ChooseKey    func(item *Item) bool
	MaxSize      uint64
	KeyToList    func(key []byte, itr *Iterator) (*pb.KVList, error)
	Send         func(buf *z.Buffer) error
	SinceTs      uint64
	readTs       uint64
	db           *DB
	rangeCh      chan keyRange
	kvChan       chan *z.Buffer
	nextStreamId atomic.Uint32
	doneMarkers  bool
	scanned      atomic.Uint64
	numProducers atomic.Int32
}

func (st *Stream) SendDoneMakers(done bool) {

}
func (st *Stream) ToList(key []byte, itr *Iterator) (*pb.KVList, error) {
	return &pb.KVList{}, nil
}

func (st *Stream) produceRanges(ctx context.Context) {

}

func (st *Stream) produceKVs(ctx context.Context, threadId int) error {
	return nil
}

func (st *Stream) streamKVs(ctx context.Context) error {
	return nil
}

func KVToBuffer(kv *pb.KV, buf *z.Buffer) {

}

func (st *Stream) Orchestrate(ctx context.Context) error {
	return nil
}

func BufferToKVList(buf *z.Buffer) (*pb.KVList, error) {
	return &pb.KVList{}, nil
}
