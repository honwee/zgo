/**
 * Copyright (C) 2021 UnionTech Software Technology Co., Ltd. All rights reserved.
 * @author 陈弘唯
 * @Email  : chenhongwei@uniontech.com
 * @date 2021/12/30 上午10:19
 */

package logs

import (
	"bufio"
	"crypto/md5" //#nosec
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/axgle/mahonia"

	"zgo/pkg/encryption"
)

type deLog struct {
	key []byte
	mu  sync.Mutex
}

// DecryptLog 日志解密接口
func DecryptLog(filePath string) error {
	/* #nosec */
	var wg sync.WaitGroup
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(file)

	reader := bufio.NewReader(file)
	d := &deLog{}

	for {
		wg.Add(1)
		//逐行读取
		lineBytes, _, err := reader.ReadLine()

		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			wg.Add(-1)
			continue
		} else if err != nil {
			fmt.Println(err)
			break
		}

		go func() {
			d.mu.Lock()
			defer d.mu.Unlock()
			//输出当前行内容
			gbkStr := string(lineBytes)
			utfStr, _ := convertEncoding(gbkStr, "GBK", "UTF-8")
			//fmt.Println(utfStr)
			err = d.generateKey(utfStr)
			if err != nil {
				fmt.Println(err)
				return
			}
			strArr := strings.Split(utfStr, ">")

			if len(strArr) > 3 {

				decrypt, err := d.decrypt(strArr[1])
				if err != nil {
					fmt.Println(err)
					return
				}

				utfStr = strings.Replace(utfStr, strings.TrimSpace(strArr[1]), string(decrypt), -1)

				bytes, err := d.decrypt(strArr[2])
				if err != nil {
					fmt.Println(err)
					return
				}

				errInfo := strings.Split(strArr[3], ":")
				if len(errInfo) > 1 {

					cfbDecrypt, err := d.decrypt(errInfo[1])
					if err != nil {
						fmt.Println(err)
						return
					}

					utfStr = strings.Replace(utfStr, strings.TrimSpace(errInfo[1]), string(cfbDecrypt), -1)
				}

				utfStr = strings.Replace(utfStr, strings.TrimSpace(strArr[2]), string(bytes), -1)

				fmt.Println(utfStr)
				wg.Done()
			}
		}()
		wg.Wait()
	}
	return nil
}

// DecryptLogToFile 日志解密接口，解密到文件
func DecryptLogToFile(filePath, target string) error {
	/* #nosec */
	var wg sync.WaitGroup
	//Open the encrypted original
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(file)

	reader := bufio.NewReader(file)

	filePointer, err := os.OpenFile(filepath.Clean(target), os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func(filePointer *os.File) {
		err := filePointer.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(filePointer)
	writer := bufio.NewWriter(filePointer)

	d := &deLog{}

	for {
		wg.Add(1)
		//逐行读取
		lineBytes, _, err := reader.ReadLine()

		if err != nil {
			if err == io.EOF {
				fmt.Println("解密完成")
				break
			}
			fmt.Println(err)
			break
		}

		go func() {
			d.mu.Lock()
			defer d.mu.Unlock()

			gbkStr := string(lineBytes)
			utfStr, _ := convertEncoding(gbkStr, "GBK", "UTF-8")

			err = d.generateKey(utfStr)
			if err != nil {
				fmt.Println(err)
				return
			}
			strArr := strings.Split(utfStr, ">")

			if len(strArr) > 3 {

				decrypt, err := d.decrypt(strArr[1])
				if err != nil {
					fmt.Println(err)
					return
				}

				utfStr = strings.Replace(utfStr, strings.TrimSpace(strArr[1]), string(decrypt), -1)

				bytes, err := d.decrypt(strArr[2])
				if err != nil {
					fmt.Println(err)
					return
				}

				errInfo := strings.Split(strArr[3], ":")
				if len(errInfo) > 1 {

					cfbDecrypt, err := d.decrypt(errInfo[1])
					if err != nil {
						fmt.Println(err)
						return
					}

					utfStr = strings.Replace(utfStr, strings.TrimSpace(errInfo[1]), string(cfbDecrypt), -1)
				}

				utfStr = strings.Replace(utfStr, strings.TrimSpace(strArr[2]), string(bytes), -1)

				_, err = writer.WriteString(utfStr + "\n")
				if err != nil {
					fmt.Println(err)
					return
				}

				wg.Done()
			}
		}()
		wg.Wait()
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

//编码转换
func convertEncoding(srcStr string, srcEncoding string, dstEncoding string) (dstStr string, err error) {
	//创建指定字符集的解码器
	srcDecoder := mahonia.NewDecoder(srcEncoding)
	dstDecoder := mahonia.NewDecoder(dstEncoding)
	//将内容转换为UTF-8字符串
	utfStr := srcDecoder.ConvertString(srcStr)

	//将UTF-8 字节转换为目标字符集的字节
	_, dstBytes, err := dstDecoder.Translate([]byte(utfStr), true)
	if err != nil {
		return
	}

	//还原为字符串并返回
	dstStr = string(dstBytes)
	return
}

//生成加密key
func (d *deLog) generateKey(original string) error {
	/* #nosec */
	var aesKey [16]byte
	strArr := strings.Split(original, " ")
	if len(strArr) > 0 {
		sArr := strings.Split(strArr[0], ".")
		if len(sArr) > 1 && len(sArr[1]) > 3 {
			key := sArr[1][:3]
			/* #nosec */
			aesKey = md5.Sum([]byte(key + "ubx2022"))
			d.key = aesKey[:]
			return nil
		}

	}
	return errors.New("generateKey err")
}

//解密接口
func (d *deLog) decrypt(original string) ([]byte, error) {
	cipher, err := base64.StdEncoding.DecodeString(strings.TrimSpace(original))
	if err != nil {
		return nil, err
	}

	decrypt, err := encryption.AesCFBDecrypt(cipher, d.key)
	if err != nil {
		return nil, err
	}
	return decrypt, nil
}
