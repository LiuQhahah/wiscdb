package y

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"hash/crc32"
	"math"
	"os"
	"time"
)

var (
	ErrEOF               = errors.New("ErrEOF: End of file")
	ErrCommitAfterFinish = errors.New("batch commit not permitted after finish")
)

type Flags int

const (
	Sync Flags = 1 << iota
	ReadOnly
)

var (
	datasyncFileFlag    = 0x0
	CastTagNoLiCrcTable = crc32.MakeTable(crc32.Castagnoli)
)

func OpenExistingFile(filename string, flags Flags) (*os.File, error) {
	return nil, nil
}
func CreateSyncedFile(filename string, sync bool) (*os.File, error) {
	return nil, nil
}

func OpenTruncateFile(filename string, sync bool) (*os.File, error) {
	return nil, nil
}
func SafeCopy(a, src []byte) []byte {
	return nil
}

func Copy(a []byte) []byte {
	return nil
}

func KeyWithTs(key []byte, ts uint64) []byte {
	return nil
}

// 解释key 字节数据的内容,从第8位开始解析
func ParseTs(key []byte) uint64 {
	if len(key) <= 8 {
		return 0
	}
	return math.MaxUint64 - binary.BigEndian.Uint64(key[len(key)-8:])
}

func CompareKeys(key1, key2 []byte) int {
	//先比较后8位，如果部位空直接返回
	if cmp := bytes.Compare(key1[:len(key1)-8], key2[:len(key2)-8]); cmp != 0 {
		return cmp
	}
	//否则比较前8位
	return bytes.Compare(key1[len(key1)-8:], key2[len(key2)-8:])
}

func ParseKey(key []byte) []byte {
	return nil
}

func SameKey(str, dst []byte) bool {
	return false
}

func FixedDuration(d time.Duration) string {
	return ""
}

// Wrapf is Wrap with extra info.
func Wrapf(err error, format string, args ...interface{}) error {
	if !debugMode {
		if err == nil {
			return nil
		}
		return fmt.Errorf(format+" error: %+v", append(args, err)...)
	}
	return errors.Wrapf(err, format, args...)
}

var debugMode = false
