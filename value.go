package main

import (
	"hash"
	"hash/crc32"
	"io"
	"math"
	"wiscdb/y"
)

var maxVlogFileSize uint32 = math.MaxUint32

const (
	BitDelete       byte  = 1 << 0
	bitValuePointer byte  = 1 << 1 //表示该key对应的是真正的value还是value对应的指针
	bitTxn          byte  = 1 << 6 //如果条目是txn的一部分，则设置。
	bitFinTxn       byte  = 1 << 7 //设置该条目是否在值日志中指示txn的结束。
	mi              int64 = 1 << 20
	vlogHeaderSize        = 20
)

type logEntry func(e Entry, vp valuePointer) error
type hashReader struct {
	r         io.Reader
	h         hash.Hash32
	bytesRead int
}

// new hash reader除了使用buffer io reader之外，还是使用了校验和
func newHashReader(r io.Reader) *hashReader {
	hashCrc := crc32.New(y.CastTagNoLiCrcTable)
	return &hashReader{
		r: r,
		h: hashCrc,
	}
}

// 在将read写到[]byte中还更新已读自己的数量以及hash值
// Read调用的是bufio将数据读取到p中
func (t *hashReader) Read(p []byte) (int, error) {
	n, err := t.r.Read(p)
	if err != nil {
		return n, err
	}
	t.bytesRead += n
	//将读取的value进行crc32校验和编码，并将写好的值写到h中.
	return t.h.Write(p[:n])
}

// 只读取一个字节大小
func (t *hashReader) ReadByte() (byte, error) {
	b := make([]byte, 1)
	_, err := t.Read(b)
	return b[0], err
}

// 返回sum32的和
func (t *hashReader) Sum32() uint32 {
	return t.h.Sum32()
}

type safeRead struct {
	k            []byte
	v            []byte
	recordOffset uint32
	vLogFile     *writeAheadLog // safeRead包含value log的file，是真正意义上的存储数据的地方
}

// return entry
// 输入 传入读取的reader对象
// 输出 返回entry，按照Entry的结构题来读取数据，验证头、长度等信息， 只返回一个entry的值
func (r *safeRead) Entry(reader io.Reader) (*Entry, error) {
	//利用buffer io的reader struct来创建hashReader
	hashReader1 := newHashReader(reader)
	var h header
	//获取到header的长度，包含meta，usermeta,key的长度，value的长度，expiredAt的长度
	hlen, err := h.DecodeFrom(hashReader1)
	if err != nil {
		return nil, err
	}

	//1<<16,1<<6等于64KB,如果k的长度大于了64KB就报error
	if h.kLen > uint32(1<<16) {
		return nil, errTruncate
	}
	kl := int(h.kLen)
	//这里面使用的数据，一开始10个字节，如果没有key的长度长，则进行扩容
	if cap(r.k) < kl {
		r.k = make([]byte, 2*kl)
	}
	vl := int(h.vLen)
	//这里面使用的数据，一开始10个字节，如果没有key的长度长，则进行扩容
	if cap(r.v) < vl {
		r.v = make([]byte, 2*vl)
	}
	e := &Entry{}
	e.offset = r.recordOffset
	e.hlen = hlen
	buf := make([]byte, h.kLen+h.vLen)

	if _, err := io.ReadFull(hashReader1, buf[:]); err != nil {
		if err == io.EOF {
			err = errTruncate
		}
		return nil, err
	}
	// TODO:  如果没有enable解密，等同于不会解密valuefile中的data
	if r.vLogFile.encryptionEnabled() {
		if buf, err = r.vLogFile.decryptKV(buf[:], r.recordOffset); err != nil {
			return nil, err
		}
	}

	//得到Entry的key,key和value都存储在buf中.
	e.Key = buf[:h.kLen]
	//得到Entry的Value
	e.Value = buf[h.kLen:]
	var crcBuf [crc32.Size]byte
	if _, err := io.ReadFull(reader, crcBuf[:]); err != nil {
		if err == io.EOF {
			err = errTruncate
		}
		return nil, err
	}

	crc := y.BytesToU32(crcBuf[:])
	if crc != hashReader1.Sum32() {
		return nil, errTruncate
	}
	e.meta = h.meta
	e.UserMeta = h.userMeta
	e.ExpiresAt = h.expiresAt
	return e, nil
}
