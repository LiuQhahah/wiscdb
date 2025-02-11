package internal

import (
	"bytes"
	"github.com/dgraph-io/ristretto/v2/z"
	"sync"
	"sync/atomic"
	"wiscdb/pb"
)

type logFile struct {
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

func (lf *logFile) Truncate(end int64) error {
	return nil
}

func (lf *logFile) encodeEntry(buf *bytes.Buffer, e *Entry, offset uint32) (int, error) {
	return 0, nil
}

func (lf *logFile) writeEntry(buf *bytes.Buffer, e *Entry, opt Options) error {
	return nil
}

func (lf *logFile) encryptionEnabled() bool {
	return false
}

func (lf *logFile) zeroNextEntry() {

}

func (lf *logFile) generateIV(offset uint32) []byte {
	return nil
}

func (lf *logFile) decryptKV(buf []byte, offset uint32) ([]byte, error) {
	return nil, nil
}

func (lf *logFile) keyID() uint64 {
	return 0
}

func (lf *logFile) doneWriting(offset uint32) error {
	return nil
}

func (lf *logFile) open(path string, flags int, fSize int64) error {
	return nil
}

func (lf *logFile) iterate(readOnly bool, offset uint32, fb logEntry) (uint32, error) {
	return 0, nil
}

func (lf *logFile) bootstrap() error {
	return nil
}

func (lf *logFile) read(p valuePointer) (buf []byte, err error) {
	return nil, nil
}
