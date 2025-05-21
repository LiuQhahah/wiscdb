package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"hash/crc32"
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

// 获取当前可用的数据加密密钥（DataKey），并在必要时生成新密钥
func (kr *KeyRegistry) LatestDataKey() (*pb.DataKey, error) {
	if len(kr.opt.EncryptionKey) == 0 {
		return nil, nil
	}
	validKey := func() (*pb.DataKey, bool) {
		diff := time.Since(time.Unix(kr.lastCreated, 0))
		// 密钥过期时间,默认值是10天
		// 验证dataKey的值,如果小于设定的10天则直接返回
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
	// 再次核实下
	key, valid = validKey()
	if valid {
		return key, nil
	}
	// 如果过期了,则重新产生
	k := make([]byte, len(kr.opt.EncryptionKey))
	//产生类似于token的信息
	// 生成新密钥
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

	// 如果支持本地存储则写道本地
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
	xor := func() error {
		if len(storageKey) == 0 {
			return nil
		}
		var err error
		k.Data, err = y.XORBlockAllocate(k.Data, storageKey, k.Iv)
		return err
	}
	var err error
	if err = xor(); err != nil {
		return y.Wrapf(err, "Error while encrypting datakey in storeDataKey")
	}
	var data []byte
	if data, err = proto.Marshal(k); err != nil {
		err = y.Wrapf(err, "Error while marshling datakey in storeDatakey")
		var err2 error
		if err2 = xor(); err2 != nil {
			return y.Wrapf(err, y.Wrapf(err2, "Error while decrypting datakey in storeDatakey").Error())
		}
		return err
	}

	var lenCrcBuf [8]byte
	binary.BigEndian.PutUint32(lenCrcBuf[0:4], uint32(len(data)))
	binary.BigEndian.PutUint32(lenCrcBuf[4:8], crc32.Checksum(data, y.CastTagNoLiCrcTable))

	y.Check2(buf.Write(lenCrcBuf[:]))
	y.Check2(buf.Write(data))
	return xor()
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
