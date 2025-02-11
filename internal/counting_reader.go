package internal

import (
	"bufio"
	"errors"
)

type countingReader struct {
	wrapped *bufio.Reader
	count   int64
}

func (r *countingReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (r *countingReader) ReadByte() (b byte, err error) {
	return 0, nil
}

var (
	errBadMagic    = errors.New("manifest has bad magic")
	errBadChecksum = errors.New("manifest has checksum mismatch")
)

type errorString struct {
	s string
}

var (
	s = errorString{"ssss"}
)
