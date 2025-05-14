package main

import "encoding/binary"

// vlog中entry的header信息
// 包含key的长度，value的长度，过期时间，meta信息以及userMeta信息
type header struct {
	kLen      uint32
	vLen      uint32
	expiresAt uint64
	meta      byte
	userMeta  byte
}

const (
	maxHeaderSize = 22
)

// given out encode into header
// 第一个字节是meta，第二个字节是userMeta，
// 将kLen和vLen转化成uint64编成可变长度然后写到out中，
// 实现灵活变化，最后将expiredAt写到out中
// The encoded header looks like
// +------+----------+------------+--------------+-----------+
// | Meta | UserMeta | Key Length | Value Length | ExpiresAt |
// +------+----------+------------+--------------+-----------+
func (h header) Encode(out []byte) int {
	out[0] = h.meta
	out[1] = h.userMeta
	index := 2
	index += binary.PutUvarint(out[index:], uint64(h.kLen))
	index += binary.PutUvarint(out[index:], uint64(h.vLen))
	index += binary.PutUvarint(out[index:], h.expiresAt)
	return index
}

func (h *header) Decode(buf []byte) int {
	h.meta = buf[0]
	h.userMeta = buf[1]
	kLen, sz := binary.Uvarint(buf[2:])
	h.kLen = uint32(kLen)
	vLen, sz1 := binary.Uvarint(buf[2+sz:])
	h.vLen = uint32(vLen)
	expiresAt, sz3 := binary.Uvarint(buf[2+sz+sz1:])
	h.expiresAt = expiresAt
	return 2 + sz + sz1 + sz3
}

// 从hashReader解析出来hash的length
// 同时要将值写到struct中
func (h *header) DecodeFrom(reader *hashReader) (int, error) {
	var err error
	//只读取一个字节
	h.meta, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	h.userMeta, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	kLen, err := binary.ReadUvarint(reader)
	if err != nil {
		return 0, err
	}
	h.kLen = uint32(kLen)
	vLen, err := binary.ReadUvarint(reader)
	if err != nil {
		return 0, err
	}
	h.vLen = uint32(vLen)
	h.expiresAt, err = binary.ReadUvarint(reader)
	if err != nil {
		return 0, err
	}
	return reader.bytesRead, nil
}
