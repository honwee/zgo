package encryption

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

func padding(src []byte, blocksize int) []byte {
	padnum := blocksize - len(src)%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(src, pad...)
}

func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	return src[:n-unpadnum]
}

func Encrypt3DES(src []byte, key []byte) []byte {
	block, _ := des.NewTripleDESCipher(key)
	src = padding(src, block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	blockmode.CryptBlocks(src, src)
	return src
}

func Decrypt3DES(src []byte, key []byte) []byte {
	block, _ := des.NewTripleDESCipher(key)
	blockmode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	blockmode.CryptBlocks(src, src)
	src = unpadding(src)
	return src
}
