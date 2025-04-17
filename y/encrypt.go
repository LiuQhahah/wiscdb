package y

import (
	"crypto/aes"
	"crypto/cipher"
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
