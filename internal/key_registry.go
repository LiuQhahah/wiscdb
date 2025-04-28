package internal

import (
	"os"
	"sync"
	"time"
	"wiscdb/pb"
)

type KeyRegistry struct {
	sync.RWMutex
	dataKeys    map[uint64]*pb.DataKey
	lastCreated int64
	nextKeyID   uint64
	fp          *os.File
	opt         KeyRegistryOptions
}

type KeyRegistryOptions struct {
	Dir                           string
	ReadOnly                      bool
	EncryptionKey                 []byte
	EncryptionKeyRotationDuration time.Duration
	InMemory                      bool
}

func (kr *KeyRegistry) LatestDataKey() (*pb.DataKey, error) {
	return nil, nil
}

func (kr *KeyRegistry) DataKey(id uint64) (*pb.DataKey, error) {
	return nil, nil
}
