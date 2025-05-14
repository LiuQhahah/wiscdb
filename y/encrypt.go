package y

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

// 随机生成16个字节
// 使用硬件随机生成token等,安全
// https://medium.com/@smafjal/understanding-crypto-rand-in-go-hardware-to-software-51798d3ebcbd
func GenerateIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	_, err := rand.Read(iv)
	return iv, err
}

// dst: 目标字节切片，用于存储加密/解密结果
// src: 源字节切片，包含要处理的数据
// key: AES 加密密钥
// iv: 初始化向量 (Initialization Vector)
// 一个使用 AES 加密算法在 CTR (计数器) 模式下进行 XOR 操作的函数
func XORBlock(dst, src, key, iv []byte) error {
	//创建 AES 加密块
	//使用提供的密钥创建 AES 加密块
	//如果密钥无效（长度不是 16、24 或 32 字节），会返回错误
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	//使用 AES 块和初始化向量创建 CTR 模式的流
	//CTR 模式将块密码转换为流密码
	stream := cipher.NewCTR(block, iv)
	//对源数据进行 XOR 操作，结果存入目标缓冲区
	//在 CTR 模式下，相同的操作既可加密也可解密
	stream.XORKeyStream(dst, src)
	return nil
}
