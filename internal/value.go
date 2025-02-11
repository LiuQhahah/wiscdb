package internal

import (
	"hash"
	"hash/crc32"
	"io"
	"math"
	"wiscdb/y"
)

var maxVlogFileSize uint32 = math.MaxUint32

const (
	bitDelete                byte  = 1 << 0
	bitValuePointer          byte  = 1 << 1
	bitDiscardEarlierVersion byte  = 1 << 2
	bitMergeEntry            byte  = 1 << 3
	bitTxn                   byte  = 1 << 6
	bitFinTxn                byte  = 1 << 7
	mi                       int64 = 1 << 20
	vlogHeaderSize                 = 20
)

type logEntry func(e Entry, vp valuePointer) error
type hashReader struct {
	r         io.Reader
	h         hash.Hash32
	bytesRead int
}

func newHashReader(r io.Reader) *hashReader {
	hash := crc32.New(y.CastTagNoLiCrcTable)
	return &hashReader{
		r: r,
		h: hash,
	}
}

func (t *hashReader) Read(p []byte) (int, error) {
	return 0, nil
}

func (t *hashReader) ReadByte() (byte, error) {
	return 0, nil
}

func (t *hashReader) Sum32() uint32 {
	return 0
}
