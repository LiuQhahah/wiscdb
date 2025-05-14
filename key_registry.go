package main

import (
	"bytes"
	"crypto/rand"
	"os"
	"sync"
	"time"
	"wiscdb/pb"
	"wiscdb/y"
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
	if len(kr.opt.EncryptionKey) == 0 {
		return nil, nil
	}
	validKey := func() (*pb.DataKey, bool) {
		diff := time.Since(time.Unix(kr.lastCreated, 0))
		if diff < kr.opt.EncryptionKeyRotationDuration {
			return kr.dataKeys[kr.nextKeyID], true
		}
		return nil, false
	}
	kr.RLock()
	key, valid := validKey()
	kr.RUnlock()
	if valid {
		return key, nil
	}
	kr.Lock()
	defer kr.Unlock()
	key, valid = validKey()
	if valid {
		return key, nil
	}
	k := make([]byte, len(kr.opt.EncryptionKey))
	//产生类似于token的信息
	iv, err := y.GenerateIV()
	if err != nil {
		return nil, err
	}
	_, err = rand.Read(k)
	if err != nil {
		return nil, err
	}
	kr.nextKeyID++
	dk := &pb.DataKey{
		KeyId:    kr.nextKeyID,
		Data:     k,
		CreateAt: time.Now().Unix(),
		Iv:       iv,
	}
	if !kr.opt.InMemory {
		buf := &bytes.Buffer{}
		if err = storeDataKey(buf, kr.opt.EncryptionKey, dk); err != nil {
			return nil, err
		}
		if _, err = kr.fp.Write(buf.Bytes()); err != nil {
			return nil, err
		}

	}
	dk.Data = k
	kr.lastCreated = dk.CreateAt
	kr.dataKeys[kr.nextKeyID] = dk
	return dk, nil
}

func storeDataKey(buf *bytes.Buffer, storageKey []byte, k *pb.DataKey) error {
	return nil
}

// 基于id获取datakey的值.
func (kr *KeyRegistry) DataKey(id uint64) (*pb.DataKey, error) {
	kr.RLock()
	defer kr.RUnlock()
	if id == 0 {
		return nil, nil
	}
	dk, ok := kr.dataKeys[id]
	if !ok {
		return nil, y.Wrapf(ErrInvalidKey, "Error for the KEY ID %d", id)
	}
	return dk, nil
}
