package internal

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/dgraph-io/ristretto/v2/z"
	"hash/crc32"
	"io"
	"strconv"
	"sync"
	"sync/atomic"
	"wiscdb/pb"
	"wiscdb/y"
)

type valueLogFile struct {
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

func (lf *valueLogFile) Truncate(end int64) error {
	return nil
}

func (lf *valueLogFile) encodeEntry(buf *bytes.Buffer, e *Entry, offset uint32) (int, error) {
	return 0, nil
}

func (lf *valueLogFile) writeEntry(buf *bytes.Buffer, e *Entry, opt Options) error {
	return nil
}

func (lf *valueLogFile) encryptionEnabled() bool {
	return false
}

func (lf *valueLogFile) zeroNextEntry() {

}

func (lf *valueLogFile) generateIV(offset uint32) []byte {
	return nil
}

func (lf *valueLogFile) decryptKV(buf []byte, offset uint32) ([]byte, error) {
	return nil, nil
}

func (lf *valueLogFile) keyID() uint64 {
	return 0
}

func (lf *valueLogFile) doneWriting(offset uint32) error {
	return nil
}

func (lf *valueLogFile) open(path string, flags int, fSize int64) error {
	return nil
}

var errTruncate = errors.New("Do truncate")

// 指的是vlog file存储的是value中的实际内容
func (lf *valueLogFile) iterate(readOnly bool, offset uint32, fn logEntry) (uint32, error) {
	if offset == 0 {
		offset = vlogHeaderSize
	}
	//返回reader Struct
	//	调用bufio package返回buffer io的reader
	reader := bufio.NewReader(lf.NewReader(int(offset)))
	read := &safeRead{
		k:            make([]byte, 10),
		v:            make([]byte, 10),
		recordOffset: offset,
		lf:           lf,
	}

	var lastCommit uint64
	var validEndoffset uint32 = offset
	var entries []*Entry
	var vptrs []valuePointer

loop:
	for {
		//读取value log中的Entry
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
		case e.isEmpty():
			break loop
		}

		var vp valuePointer
		vp.Len = uint32(len(e.Key) + len(e.Value) + crc32.Size + e.hlen)
		read.recordOffset += vp.Len
		vp.Offset = e.offset
		vp.Fid = lf.fid

		switch {
		case e.meta&bitTxn > 0:
			txnTs := y.ParseTs(e.Key)
			if lastCommit == 0 {
				lastCommit = txnTs
			}
			if lastCommit != txnTs {
				break loop
			}
			entries = append(entries, e)
			vptrs = append(vptrs, vp)

		case e.meta&bitFinTxn > 0:
			txnTs, err := strconv.ParseUint(string(e.Value), 10, 64)
			if err != nil || lastCommit != txnTs {
				break loop
			}
			lastCommit = 0
			validEndoffset = read.recordOffset
			for i, e := range entries {
				vp := vptrs[i]
				if err := fn(*e, vp); err != nil {
					if err == errStop {
						break
					}
					return 0, errFile(err, lf.path, "Iteration function")
				}
			}
			entries = entries[:0]
			vptrs = vptrs[:0]
		default:
			if lastCommit != 0 {
				break loop
			}
			validEndoffset = read.recordOffset
			if err := fn(*e, vp); err != nil {
				if err == errStop {
					break
				}
				return 0, errFile(err, lf.path, "Iteration function")
			}
		}
	}
	return validEndoffset, nil
}

var errStop = errors.New("Stop iteration")

func (lf *valueLogFile) bootstrap() error {
	return nil
}

func (lf *valueLogFile) read(p valuePointer) (buf []byte, err error) {
	return nil, nil
}
