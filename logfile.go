package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	cryptorand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/dgraph-io/ristretto/v2/z"
	"hash/crc32"
	"io"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"wiscdb/pb"
	"wiscdb/y"
)

type writeAheadLog struct {
	*z.MmapFile
	path     string
	lock     sync.RWMutex
	fid      uint32
	size     atomic.Uint32
	dataKey  *pb.DataKey
	baseIV   []byte // TODO: 存储的是什么
	registry *KeyRegistry
	writeAt  uint32
	opt      Options
}

/*
*
传入
*/
func (vLogFile *writeAheadLog) Truncate(end int64) error {
	if fi, err := vLogFile.Fd.Stat(); err != nil {
		return fmt.Errorf("while file.stat on file: %s,error: %v\n", vLogFile.Fd.Name(), err)
	} else if fi.Size() == end {
		//无需截断
		return nil
	}
	y.AssertTrue(!vLogFile.opt.ReadOnly)
	//记录当前文件大小
	vLogFile.size.Store(uint32(end))
	//执行文件截断
	return vLogFile.MmapFile.Truncate(end)
}

// 对entry进行编码
func (vLogFile *writeAheadLog) encodeEntry(buf *bytes.Buffer, e *Entry, offset uint32) (int, error) {
	h := header{
		kLen:      uint32(len(e.Key)),
		vLen:      uint32(len(e.Value)),
		expiresAt: e.ExpiresAt,
		meta:      e.meta,
		userMeta:  e.UserMeta,
	}
	hash := crc32.New(y.CastTagNoLiCrcTable)
	//使用io的mulit writer可以实现buf和hash串行写到文件中
	writer := io.MultiWriter(buf, hash)

	var headerEnc [maxHeaderSize]byte
	//对header进行编码
	sz := h.Encode(headerEnc[:])
	//将header写到writer中
	y.Check2(writer.Write(headerEnc[:sz]))
	//判断value log是否开启加密的条件是Datakey是否为null，这个条件有点怪异
	if vLogFile.encryptionEnabled() {
		eBuf := make([]byte, 0, len(e.Key)+len(e.Value))
		//将e.Key和e.Value展开并append到eBuf中
		// 使用 的是go的语法糖
		eBuf = append(eBuf, e.Key...)
		eBuf = append(eBuf, e.Value...)
		//传入参数计算校验合
		if err := y.XORBlockStream(
			writer, eBuf, vLogFile.dataKey.Data, vLogFile.generateIV(offset)); err != nil {
			return 0, y.Wrapf(err, "Error while encoding entry for vlog.")
		}
	} else {
		//不采用校验合的方法就是直接写入key和value
		y.Check2(writer.Write(e.Key))
		y.Check2(writer.Write(e.Value))
	}
	var crcBuf [crc32.Size]byte
	binary.BigEndian.PutUint32(crcBuf[:], hash.Sum32())
	y.Check2(buf.Write(crcBuf[:]))

	return len(headerEnc[:sz]) + len(e.Key) + len(e.Value) + len(crcBuf), nil
}

func (vLogFile *writeAheadLog) writeEntry(buf *bytes.Buffer, e *Entry, opt Options) error {
	buf.Reset()
	pLen, err := vLogFile.encodeEntry(buf, e, vLogFile.writeAt)
	if err != nil {
		return err
	}
	y.AssertTrue(pLen == copy(vLogFile.Data[vLogFile.writeAt:], buf.Bytes()))
	vLogFile.writeAt += uint32(pLen)
	vLogFile.zeroNextEntry()
	return nil
}

func (vLogFile *writeAheadLog) encryptionEnabled() bool {
	return vLogFile.dataKey != nil
}

// 为下一个Entry写入做到清扫工作
// 在预写日志（Write-Ahead Log, WAL）中清零下一个条目的头部区域，以确保后续写入不会读到脏数据或残留数据
func (vLogFile *writeAheadLog) zeroNextEntry() {
	z.ZeroOut(vLogFile.Data, int(vLogFile.writeAt), int(vLogFile.writeAt+maxHeaderSize))
}

// IV指的是 Initialization vector
// 16个字节数组中存储内容是，前12个字节存储的baseIV的值，后四个字节存储的偏移量
func (vLogFile *writeAheadLog) generateIV(offset uint32) []byte {
	iv := make([]byte, aes.BlockSize)
	//前12个字节用来存储baseIV
	y.AssertTrue(12 == copy(iv[:12], vLogFile.baseIV))
	//后四个字节存储偏移量
	binary.BigEndian.PutUint32(iv[12:], offset)
	return iv
}

// 对KV进行解密
func (vLogFile *writeAheadLog) decryptKV(buf []byte, offset uint32) ([]byte, error) {
	return y.XORBlockAllocate(buf, vLogFile.dataKey.Data, vLogFile.generateIV(offset))
}

func (vLogFile *writeAheadLog) keyID() uint64 {
	return 0
}

// 执行完后，该文件id就结束了
func (vLogFile *writeAheadLog) doneWriting(offset uint32) error {
	if vLogFile.opt.SyncWrites {
		if err := vLogFile.Sync(); err != nil {
			return y.Wrapf(err, "Unable to sync value log: %q", vLogFile.path)
		}
	}
	vLogFile.lock.Lock()
	defer vLogFile.lock.Unlock()
	if err := vLogFile.Truncate(int64(offset)); err != nil {
		return y.Wrapf(err, "Unable to truncate file: %q", vLogFile.path)
	}
	return nil
}

func (vLogFile *writeAheadLog) open(path string, flags int, fSize int64) error {
	mf, ferr := z.OpenMmapFile(path, flags, int(fSize))
	vLogFile.MmapFile = mf

	// 如果是个新文件,初始化valueLog
	if ferr == z.NewFile {
		if err := vLogFile.bootstrap(); err != nil {
			os.Remove(path)
			return err
		}
		vLogFile.size.Store(vlogHeaderSize)
	} else if ferr != nil {
		return y.Wrapf(ferr, "while opening file: %s", path)
	}
	vLogFile.size.Store(uint32(len(vLogFile.Data)))

	if vLogFile.size.Load() < vlogHeaderSize {
		return nil
	}
	buf := make([]byte, vlogHeaderSize)
	y.AssertTruef(vlogHeaderSize == copy(buf, vLogFile.Data), "Unable to copy from %s, size %d", path, vLogFile.size.Load())
	keyID := binary.BigEndian.Uint64(buf[:8])

	if dk, err := vLogFile.registry.DataKey(keyID); err != nil {
		return y.Wrapf(err, "While opening vlog file %d", vLogFile.fid)
	} else {
		vLogFile.dataKey = dk
	}
	vLogFile.baseIV = buf[8:]
	y.AssertTrue(len(vLogFile.baseIV) == 12)

	return ferr
}

var errTruncate = errors.New("Do truncate")

// 指的是vlog file存储的是value中的实际内容
// 输入 offset指的是从哪个位置开始迭代、可以从头开始迭代，logEntry是个函数用于执行具体迭代工作
// 输出
func (vLogFile *writeAheadLog) iterate(readOnly bool, offset uint32, fn logEntry) (uint32, error) {
	//如果偏移量为0 表明从头开始迭代，但是在valuelog中头并不是0，而是从文件headersize开始后读取，
	//用20个字节来描述文件头，包含魔数等信息
	if offset == 0 {
		offset = vlogHeaderSize
	}
	//返回reader Struct
	//	调用bufio package返回buffer io的reader
	reader := bufio.NewReader(vLogFile.NewReader(int(offset)))
	read := &safeRead{
		k:            make([]byte, 10),
		v:            make([]byte, 10),
		recordOffset: offset,
		vLogFile:     vLogFile,
	}

	var lastCommit uint64
	var validEndoffset uint32 = offset
	var entries []*Entry
	var vptrs []valuePointer

	//	直到读到文件截止符才终止循环，或者读到entry的值为null
loop:
	for {
		//读取value log中的Entry
		// 解析Entry
		entryResult, err := read.Entry(reader)
		switch {
		case err == io.EOF:
			break loop
		case err == io.ErrUnexpectedEOF || err == errTruncate:
			break loop
		case err != nil:
			return 0, err
		case entryResult == nil:
			continue
		case entryResult.isEmpty():
			break loop
		}

		var vp valuePointer
		vp.Len = uint32(len(entryResult.Key) + len(entryResult.Value) + crc32.Size + entryResult.hlen)
		read.recordOffset += vp.Len
		vp.Offset = entryResult.offset
		vp.Fid = vLogFile.fid

		// go中的switch不需要break如果满足直接执行case中内容，如果不满足才会执行default中操作
		switch {
		//如果entry的mete值 第6位有值
		//	表明是正常的Entry
		case entryResult.meta&bitTxn > 0:
			//获取事务提交的时间戳
			txnTs := y.ParseTs(entryResult.Key)

			if lastCommit == 0 {
				lastCommit = txnTs
			}
			// 如果上次提交的时间戳和读取到的时间戳一致，那么就跳出循环.
			if lastCommit != txnTs {
				break loop
			}
			//将读取到的entry存储在entries数组中
			entries = append(entries, entryResult)
			vptrs = append(vptrs, vp)

		//如果entry的meta值第7位有值
		// 表明是事务的最后一个Entry
		// 会调用传入的函数fn执行
		case entryResult.meta&bitFinTxn > 0:
			txnTs, err := strconv.ParseUint(string(entryResult.Value), 10, 64)
			if err != nil || lastCommit != txnTs {
				break loop
			}
			lastCommit = 0
			validEndoffset = read.recordOffset
			// 一个entries对应一个valuepointer
			// 任务是将key写到跳表中
			for i, e := range entries {
				vp := vptrs[i]
				if err := fn(*e, vp); err != nil {
					if err == errStop {
						break
					}
					return 0, errFile(err, vLogFile.path, "Iteration function")
				}
			}
			entries = entries[:0]
			vptrs = vptrs[:0]
		default:
			if lastCommit != 0 {
				break loop
			}
			validEndoffset = read.recordOffset
			if err := fn(*entryResult, vp); err != nil {
				if err == errStop {
					break
				}
				return 0, errFile(err, vLogFile.path, "Iteration function")
			}
		}
	}
	//最后返回偏移量
	return validEndoffset, nil
}

var errStop = errors.New("Stop iteration")

func (vLogFile *writeAheadLog) bootstrap() error {
	var err error
	var dk *pb.DataKey

	// 读取最新的DataKey
	if dk, err = vLogFile.registry.LatestDataKey(); err != nil {
		return y.Wrapf(err, "Error while retrieving datakey in logFile.bootstarp")
	}
	vLogFile.dataKey = dk
	buf := make([]byte, vlogHeaderSize)
	binary.BigEndian.PutUint64(buf[:8], vlogHeaderSize)
	binary.BigEndian.PutUint64(buf[:8], vLogFile.keyID())
	if _, err := cryptorand.Read(buf[:8]); err != nil {
		return y.Wrapf(err, "Error while creating base IV, while creating logfile")
	}

	vLogFile.baseIV = buf[8:]
	y.AssertTrue(len(vLogFile.baseIV) == 12)
	y.AssertTrue(vlogHeaderSize == copy(vLogFile.Data[0:], buf))

	vLogFile.zeroNextEntry()
	return nil
}

func (vLogFile *writeAheadLog) read(p valuePointer) (buf []byte, err error) {
	return nil, nil
}
