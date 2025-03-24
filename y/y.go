package y

import (
	"fmt"
	"github.com/pkg/errors"
	"hash/crc32"
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

func ParseTs(key []byte) uint64 {
	return 0
}

func CompareKeys(key1, key2 []byte) int {
	return 0
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
