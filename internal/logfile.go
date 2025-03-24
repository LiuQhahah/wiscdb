package internal

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/dgraph-io/ristretto/v2/z"
	"hash/crc32"
	"io"
	"sync"
	"sync/atomic"
	"wiscdb/pb"
)

type writeAheadLog struct {
	*z.MmapFile
	path     string
	lock     sync.RWMutex
	fid      uint32
	size     atomic.Uint32
	dataKey  *pb.DataKey
	baseIV   []byte
	registry *KeyRegistry
	writeAt  uint32
	opt      Options
}

func (lf *writeAheadLog) Truncate(end int64) error {
	return nil
}

func (lf *writeAheadLog) encodeEntry(buf *bytes.Buffer, e *Entry, offset uint32) (int, error) {
	return 0, nil
}

func (lf *writeAheadLog) writeEntry(buf *bytes.Buffer, e *Entry, opt Options) error {
	return nil
}

func (lf *writeAheadLog) encryptionEnabled() bool {
	return false
}

func (lf *writeAheadLog) zeroNextEntry() {

}

func (lf *writeAheadLog) generateIV(offset uint32) []byte {
	return nil
}

func (lf *writeAheadLog) decryptKV(buf []byte, offset uint32) ([]byte, error) {
	return nil, nil
}

func (lf *writeAheadLog) keyID() uint64 {
	return 0
}

func (lf *writeAheadLog) doneWriting(offset uint32) error {
	return nil
}

func (lf *writeAheadLog) open(path string, flags int, fSize int64) error {
	return nil
}

var errTruncate = errors.New("Do truncate")

func (lf *writeAheadLog) iterate(readOnly bool, offset uint32, fb logEntry) (uint32, error) {
	if offset == 0 {
		offset = vlogHeaderSize
	}
	reader := bufio.NewReader(lf.NewReader(int(offset)))
	read := &safeRead{
		k:          make([]byte, 10),
		v:          make([]byte, 10),
		readOffset: offset,
		lf:         lf,
	}

	var lastCommit uint64
	var validEndoffset uint32 = offset
	var entries []*Entry
	var vptrs []valuePointer

loop:
	for {
		e, err := read.Entry(reader)
		switch {
		case err == io.EOF:
			break loop
		case err == io.ErrUnexpectedEOF || err == errTruncate:
			break loop
		case err != nil:
			return 0, err
		case e == nil:
			continue
		case e.isZero():
			break loop
		}

		var vp valuePointer
		vp.Len = uint32(len(e.Key) + len(e.Value) + crc32.Size + e.hlen)
		read.readOffset += vp.Len
		vp.Offset = e.offset
		vp.Fid = lf.fid

		switch {
		case e.meta&bitTxn > 0:

		case e.meta&bitFinTxn > 0:
		default:

		}
	}
	return validEndoffset, nil
}

func (lf *writeAheadLog) bootstrap() error {
	return nil
}

func (lf *writeAheadLog) read(p valuePointer) (buf []byte, err error) {
	return nil, nil
}
