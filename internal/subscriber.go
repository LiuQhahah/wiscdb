package internal

import (
	"github.com/dgraph-io/ristretto/v2/z"
	"sync/atomic"
	"wiscdb/pb"
)

type subscriber struct {
	id        uint64
	matches   []pb.Match
	sendCh    chan *pb.KVList
	subCloser *z.Closer
	active    *atomic.Uint64
}
