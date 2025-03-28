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
	hashCrc := crc32.New(y.CastTagNoLiCrcTable)
	return &hashReader{
		r: r,
		h: hashCrc,
	}
}

// 在将read写到[]byte中还更新已读自己的数量以及hash值
func (t *hashReader) Read(p []byte) (int, error) {
	n, err := t.r.Read(p)
	if err != nil {
		return n, err
	}
	t.bytesRead += n
	return t.h.Write(p[:n])
}

// 只读取一个字节大小
func (t *hashReader) ReadByte() (byte, error) {
	b := make([]byte, 1)
	_, err := t.Read(b)
	return b[0], err
}

func (t *hashReader) Sum32() uint32 {
	return 0
}

type safeRead struct {
	k          []byte
	v          []byte
	readOffset uint32
	lf         *valueLogFile
}

// return entry
func (r *safeRead) Entry(reader io.Reader) (*Entry, error) {
	hReader := newHashReader(reader)
	var h header
	hlen, err := h.DecodeFrom(hReader)
	if err != nil {
		return nil, err
	}

	if h.kLen > uint32(1<<16) {
		return nil, errTruncate
	}
	kl := int(h.kLen)
	if cap(r.k) < kl {
		r.k = make([]byte, 2*kl)
	}
	vl := int(h.vLen)
	if cap(r.v) < vl {
		r.v = make([]byte, 2*vl)
	}
	e := &Entry{}
	e.offset = r.readOffset
	e.hlen = hlen
	buf := make([]byte, h.kLen+h.vLen)

	if _, err := io.ReadFull(hReader, buf[:]); err != nil {
		if err == io.EOF {
			err = errTruncate
		}
		return nil, err
	}

	return &Entry{}, nil
}
