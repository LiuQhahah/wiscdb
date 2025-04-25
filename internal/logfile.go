package internal

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto/v2/z"
	"hash/crc32"
	"io"
	"strconv"
	"sync"
	"sync/atomic"
	"wiscdb/pb"
	"wiscdb/y"
)

type valueLogFile struct {
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
func (vLogFile *valueLogFile) Truncate(end int64) error {
	if fi, err := vLogFile.Fd.Stat(); err != nil {
		return fmt.Errorf("while file.stat on file: %s,error: %v\n", vLogFile.Fd.Name(), err)
	} else if fi.Size() == end {
		return nil
	}
	y.AssertTrue(!vLogFile.opt.ReadOnly)
	vLogFile.size.Store(uint32(end))
	return vLogFile.MmapFile.Truncate(end)
}

func (vLogFile *valueLogFile) encodeEntry(buf *bytes.Buffer, e *Entry, offset uint32) (int, error) {
	return 0, nil
}

func (vLogFile *valueLogFile) writeEntry(buf *bytes.Buffer, e *Entry, opt Options) error {
	return nil
}

func (vLogFile *valueLogFile) encryptionEnabled() bool {
	return vLogFile.dataKey != nil
}

func (vLogFile *valueLogFile) zeroNextEntry() {

}

// IV指的是 Initialization vector
// 16个字节数组中存储内容是，前12个字节存储的baseIV的值，后四个字节存储的偏移量
func (vLogFile *valueLogFile) generateIV(offset uint32) []byte {
	iv := make([]byte, aes.BlockSize)
	//前12个字节用来存储baseIV
	y.AssertTrue(12 == copy(iv[:12], vLogFile.baseIV))
	//后四个字节存储偏移量
	binary.BigEndian.PutUint32(iv[12:], offset)
	return iv
}

// 对KV进行解密
func (vLogFile *valueLogFile) decryptKV(buf []byte, offset uint32) ([]byte, error) {
	return y.XORBlockAllocate(buf, vLogFile.dataKey.Data, vLogFile.generateIV(offset))
}

func (vLogFile *valueLogFile) keyID() uint64 {
	return 0
}

func (vLogFile *valueLogFile) doneWriting(offset uint32) error {
	return nil
}

func (vLogFile *valueLogFile) open(path string, flags int, fSize int64) error {
	return nil
}

var errTruncate = errors.New("Do truncate")

// 指的是vlog file存储的是value中的实际内容
// 输入 offset指的是从哪个位置开始迭代、可以从头开始迭代，logEntry是个函数用于执行具体迭代工作
// 输出
func (vLogFile *valueLogFile) iterate(readOnly bool, offset uint32, fn logEntry) (uint32, error) {
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

		switch {
		//如果entry的mete值 第6位有值
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
		// TODO: 暂时没有找到设置meta的code
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

func (vLogFile *valueLogFile) bootstrap() error {
	return nil
}

func (vLogFile *valueLogFile) read(p valuePointer) (buf []byte, err error) {
	return nil, nil
}
