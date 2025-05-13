package y

import (
	"errors"
	"github.com/cespare/xxhash/v2"
	"hash/crc32"
	"wiscdb/pb"
)

var ErrChecksumMismatch = errors.New("checksum mismatch")

// 具体做事的
func VerifyCheckSum(data []byte, expected *pb.CheckSum) error {
	actual := CalculateCheckSum(data, expected.Algo)
	if actual != expected.Sum {
		return Wrapf(ErrChecksumMismatch, "actual: %d, expected: %d", actual, expected.Sum)
	}
	return nil
}

func CalculateCheckSum(data []byte, ct pb.ChecksumAlgorithm) uint64 {
	switch ct {
	case pb.CheckSumCRC32C:
		return uint64(crc32.Checksum(data, CastTagNoLiCrcTable))
	case pb.CheckSumXXHash64:
		return xxhash.Sum64(data)
	default:
		panic("check sum type error")
	}
}
