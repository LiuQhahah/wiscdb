package internal

import (
	"time"
)

type Entry struct {
	Key          []byte
	Value        []byte
	ExpiresAt    uint64
	version      uint64
	offset       uint32
	UserMeta     byte
	meta         byte
	hlen         int
	valThreshold int64
}

func (e *Entry) isZero() bool {
	return false
}

func (e *Entry) estimateSizeAndSetThreshold(threshold int64) int64 {
	//	TODO:
	return 0
}

func (e *Entry) skipVlogAndSetThreshold(threshold int64) bool {
	return false
}

func (e Entry) print(prefix string) {

}

func NewEntry(key, value []byte) *Entry {
	return nil
}

func (e *Entry) WithMeta(meta byte) *Entry {
	return nil
}

func (e *Entry) WithDiscard() *Entry {
	return nil
}

func (e *Entry) WithTTL(dur time.Duration) *Entry {
	return nil
}

func (e *Entry) withMergeBit() *Entry {
	return nil
}
