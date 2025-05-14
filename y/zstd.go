package y

import (
	"github.com/klauspost/compress/zstd"
	"sync"
)

var (
	decoder *zstd.Decoder
	encoder *zstd.Encoder

	encOnce, decOnce sync.Once
)

func ZSTDDeCompress(dst, src []byte) ([]byte, error) {
	return nil, nil
}
func ZSTDCompress(dst, src []byte, compressionLevel int) ([]byte, error) {
	encOnce.Do(func() {
		var err error
		level := zstd.EncoderLevelFromZstd(compressionLevel)
		encoder, err = zstd.NewWriter(nil, zstd.WithEncoderLevel(level))
		Check(err)
	})
	return encoder.EncodeAll(src, dst[:0]), nil
}

func ZSTDCompressBound(srcSize int) int {
	lowLimit := 128 << 10
	var margin int
	if srcSize < lowLimit {
		margin = (lowLimit - srcSize) >> 11
	}
	return srcSize + (srcSize >> 8) + margin
}
