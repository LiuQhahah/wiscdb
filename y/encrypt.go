package y

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
)

// key,传入的是data key.
func XORBlockAllocate(src, key, iv []byte) ([]byte, error) {
	//对key采用AES加密
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//
	stream := cipher.NewCTR(block, iv)
	dst := make([]byte, len(src))
	//对stream进行加密，然后将值返回给dst
	stream.XORKeyStream(dst, src)
	return dst, nil
}

// 计算校验合
func XORBlockStream(w io.Writer, src, key, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)
	sw := cipher.StreamWriter{S: stream, W: w}
	_, err = io.Copy(sw, bytes.NewReader(src))
	return Wrapf(err, "XORBlockStream")
}
