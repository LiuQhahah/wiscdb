package y

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"hash/crc32"
	"math"
	"os"
	"reflect"
	"time"
	"unsafe"
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
	b := make([]byte, len(a))
	copy(b, a)
	return b
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

// 为了效率高，并不会将key解码出来，而是直接用字节比
// 那么对于像MySQL这些数据库使用like 或者> 后台是不是也是进行自己数组比较而不是翻译成int或者string呢
func CompareKeys(key1, key2 []byte) int {
	//先比较前面部位，比如key1和key2的长度分别是10和12,
	//那么就比较key1[0:2]和key2[0:4]很显然key2大，
	//采用这个策略可以提高效率.
	if cmp := bytes.Compare(key1[:len(key1)-8], key2[:len(key2)-8]); cmp != 0 {
		return cmp
	}
	//否则比较最后的8位
	return bytes.Compare(key1[len(key1)-8:], key2[len(key2)-8:])
}

// 所有key的前8个字节都不是真正的开头
func ParseKey(key []byte) []byte {
	if key == nil {
		return nil
	}
	return key[:(len(key) - 8)]
}

func U32SliceToBytes(u32s []uint32) []byte {
	if len(u32s) == 0 {
		return nil
	}
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Len = len(u32s) * 4
	hdr.Cap = hdr.Len
	hdr.Data = uintptr(unsafe.Pointer(&u32s[0]))
	return b
}
func BytesToU32Slice(b []byte) []uint32 {
	return nil
}

func SameKey(src, dst []byte) bool {
	if len(src) != len(dst) {
		return false
	}
	return bytes.Equal(ParseKey(src), ParseKey(dst))
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
