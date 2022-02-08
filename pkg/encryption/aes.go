package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"time"
)

/*AesCFBEncrypt Encrypt data use AES-CFB mode
The key argument should be the AES key,
either 16, 24, or 32 bytes to select
AES-128, AES-192, or AES-256.
*/
func AesCFBEncrypt(rawData, key []byte) ([]byte, error) {
	var cipherDatas []byte
	var err error

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherDatas = make([]byte, block.BlockSize()+len(rawData))
	//Use timestamp as IV,timestamp is 8byte length,but AES block
	//is 16byte,so just left other 8bytes to zero, that's ok, the
	//IV is unique
	iv := cipherDatas[:block.BlockSize()]
	binary.LittleEndian.PutUint64(iv, uint64(time.Now().UnixNano()))

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherDatas[aes.BlockSize:], rawData)

	return cipherDatas, err
}

/*AesCFBDecrypt Decrypt data use AES-CFB mode
The key argument should be the AES key,
either 16, 24, or 32 bytes to select
AES-128, AES-192, or AES-256.
*/
func AesCFBDecrypt(ciphers, key []byte) ([]byte, error) {
	//cipher data should equal,big than aes.BlockSize
	//it's never less than aes.BlockSize
	if len(ciphers) < aes.BlockSize {
		return nil, errors.New("cipher text too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//Get IV from cipher's head
	iv := ciphers[:aes.BlockSize]
	ciphers = ciphers[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphers, ciphers)

	return ciphers, nil
}

//GenAesKey Generate a Aes key with give size
//keySize:must be 16,24,32bytes.
//return 32 bytes size if keysize is invalid
func GenAesKey(keySize int) []byte {
	var key []byte

	switch keySize {
	case 16:
		key = make([]byte, 16)
	case 24:
		key = make([]byte, 24)
	case 32:
		fallthrough
	default:
		key = make([]byte, 32)
	}

	//Timestamp 8byte part
	binary.LittleEndian.PutUint64(key, uint64(time.Now().UnixNano()))

	//Fill left with random data
	randPart := key[8:]
	_, err := rand.Read(randPart)

	//Timestamp is unique & enough,if read faild, it's doen't matter.
	//Just fill any data is ok
	if nil != err {
		binary.LittleEndian.PutUint64(key, uint64(0xDeadDeadDeadDead))
	}

	return key
}
