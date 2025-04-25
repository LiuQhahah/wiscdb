package internal

import (
	"github.com/dgraph-io/ristretto/v2/z"
	"sync"
	"wiscdb/table"
	"wiscdb/y"
)

type StreamWriter struct {
	writeLock  sync.Mutex
	db         *DB
	done       func()
	throttle   *y.Throttle
	maxVersion uint64
	writers    map[uint32]*sortedWriter
	prevLevel  int
}

func (db *DB) NewStreamWriter() *StreamWriter {
	return &StreamWriter{
		db:       db,
		throttle: y.NewThrottle(16), //channel的长度是16
		writers:  make(map[uint32]*sortedWriter),
	}
}

func (sw *StreamWriter) Prepare() error {
	return nil
}

func (sw *StreamWriter) PrepareIncremental() {

}
func (sw *StreamWriter) Write(buf *z.Buffer) error {
	return nil
}

func (sw *StreamWriter) Flush() error {
	return nil
}

func (sw *StreamWriter) Cancel() {

}
func (sw *StreamWriter) newWriter(streamID uint32) (*sortedWriter, error) {
	return nil, nil
}

type sortedWriter struct {
	db       *DB
	throttle *y.Throttle
	opts     table.Options
	builder  *table.Builder
	lastKey  []byte
	level    int
	streamID uint32
	reqCh    chan *request
	closer   *z.Closer
}

func (sw *sortedWriter) handleRequests() {

}

func (sw *sortedWriter) Add(key []byte, vs y.ValueStruct) error {
	return nil
}

func (sw *sortedWriter) send(done bool) error {
	return nil
}

func (sw *sortedWriter) Done() error {
	return nil
}

func (sw *sortedWriter) createTable(builder *table.Builder) error {
	return nil
}
