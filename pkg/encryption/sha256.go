package encryption

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

/*
描述 ：获取字符串的sha256值
参数 ：
	str:目标数据
返回值：
	sha256值及错误码
*/
func HashSHA256String(str string) (string, error) {
	var retStr string
	hash := sha256.New()
	_, err := hash.Write([]byte(str))
	if err != nil {
		return retStr, err
	}
	hashValue := hash.Sum(nil)
	retStr = hex.EncodeToString(hashValue)

	return retStr, nil
}

/*
描述 ：获取文件的sha256值
参数 ：
	filePath:文件路径
返回值：
	sha256值及错误码
*/
func HashSHA256File(filePath string) (string, error) {
	var retStr string
	file, err := os.Open(filePath)
	if err != nil {
		return retStr, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return retStr, err
	}
	hashValue := hash.Sum(nil)
	retStr = hex.EncodeToString(hashValue)
	return retStr, nil
}
