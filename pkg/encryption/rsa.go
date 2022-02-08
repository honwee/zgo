package encryption

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"log"
	"sync"
)

// 进程通信加密用
var pubkey = `-----BEGIN 公钥-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDz2fyPpk9d/0T6Tjm4SbetWjNu
9kbDQSCk1oCrwNSq+8QrYFtAMNQgsN5wvDOAj3OREoSL2c2CVsVGl7C63HJNCpJf
fn+kgPnDsK0LIypjFXhpeVluzPsFGIAfLtep0bXrU53YvO2jsnYKM+do/b4S11hT
/L+WxAlwAoGZuQMm2wIDAQAB
-----END 公钥-----
`

// 进程通信加密用
var pirvatekey = `-----BEGIN 私钥-----
MIICXAIBAAKBgQDz2fyPpk9d/0T6Tjm4SbetWjNu9kbDQSCk1oCrwNSq+8QrYFtA
MNQgsN5wvDOAj3OREoSL2c2CVsVGl7C63HJNCpJffn+kgPnDsK0LIypjFXhpeVlu
zPsFGIAfLtep0bXrU53YvO2jsnYKM+do/b4S11hT/L+WxAlwAoGZuQMm2wIDAQAB
AoGAVEfERfXqOoeu1IBS7MH1zOF/I1vVS0joOnC02if0mQAZZhCQmVgHCSF4UCiL
+GQcQkjPLPLjV6gb2PE2sO7eRdum8eyQCw4UE3ita8lNr3C0W/dJ3qE8V4bTpHVS
RT9vRtEd8sYOUuZ9JPPC/F8Nsn4hujnm1u8kw84i6Bp/UkkCQQD9Cv0sCKmTID/P
jJKjGzMu+DxJsOcAK713o4Zvk29oRmKcw6GyFWLGFRW80r4gl2OSbFbxM54eeZ60
T1bRMeF1AkEA9rOAEwno2x7yr7GDsTemolGDuqDyos/hqROIOWAo/NRbD7MT4cuH
M2wSMLcy8Eu56AGVGBo7Awg6m00lW5QNDwJAejrZqnCQwRHd4PqtRn54DeM48/uw
yeNXBTiHUtQsB3mgXssdCzHLYZWDx48g6gtWvL76jE57vYrP/5cnf6uRlQJAWIZu
9eX/bem8Ejmz1PrwS5zOlUC98JiCFGbS4ivUaW1WQ9rxznt3R4eHO33xxHKYAl3W
/3AiLuNcDHBxcFw/FwJBAJldqmNyNBG8SqxXUeJgBPaZd+tFOxPGTIKAi/Oy98Nk
gC8k1lhCLdY6ezs757iWin6FYL032pHglqUsnxDAJy4=
-----END 私钥-----
`

type rsaSecurity struct {
	pubStr string          //公钥字符串
	priStr string          //私钥字符串
	pubkey *rsa.PublicKey  //公钥
	prikey *rsa.PrivateKey //私钥
}

var rasForSigature = &rsaSecurity{}

var (
	once sync.Once        // 	来判断某个函数是否已经执行过一次，如果没有就执行，如果执行过就不会再执行
	rsas = &rsaSecurity{} //	实例
)

// 设置公钥
func (rsas *rsaSecurity) setPublicKey(pubStr string) (err error) {
	rsas.pubStr = pubStr
	rsas.pubkey, err = rsas.getPublickey()
	return err
}

// 设置私钥
func (rsas *rsaSecurity) setPrivateKey(priStr string) (err error) {
	rsas.priStr = priStr
	rsas.prikey, err = rsas.getPrivatekey()
	return err
}

// *rsa.PublicKey
func (rsas *rsaSecurity) getPrivatekey() (*rsa.PrivateKey, error) {
	return getPriKey([]byte(rsas.priStr))
}

// *rsa.PrivateKey
func (rsas *rsaSecurity) getPublickey() (*rsa.PublicKey, error) {
	return getPubKey([]byte(rsas.pubStr))
}

// 公钥加密
func (rsas *rsaSecurity) PubKeyENCTYPT(input []byte) ([]byte, error) {
	if rsas.pubkey == nil {
		return []byte(""), errors.New(`Please set the public key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(rsas.pubkey, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}

// 公钥解密
func (rsas *rsaSecurity) PubKeyDECRYPT(input []byte) ([]byte, error) {
	if rsas.pubkey == nil {
		return []byte(""), errors.New(`Please set the public key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(rsas.pubkey, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}

// 私钥加密
func (rsas *rsaSecurity) PriKeyENCTYPT(input []byte) ([]byte, error) {
	if rsas.prikey == nil {
		return []byte(""), errors.New(`Please set the private key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(rsas.prikey, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}

// 私钥解密
func (rsas *rsaSecurity) PriKeyDECRYPT(input []byte) ([]byte, error) {
	if rsas.prikey == nil {
		return []byte(""), errors.New(`Please set the private key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(rsas.prikey, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}

	return ioutil.ReadAll(output)
}

func (rsas *rsaSecurity) init() error {
	if err := rsas.setPublicKey(pubkey); err != nil {
		log.Fatalln(`set public key :`, err)

		return err
	}
	if err := rsas.setPrivateKey(pirvatekey); err != nil {
		log.Fatalln(`set private key :`, err)
		return err
	}

	return nil
}

func GetInstance() *rsaSecurity {
	once.Do(func() {
		rsas = &rsaSecurity{}
		rsas.init()
	})

	return rsas
}

/*
描述 ：公钥解密--主要用于插件签名加解密
参数 ：
	pubStr:公钥
	input:加密源数据
返回值：
	解密的目标数据、及错误码
*/
func PubKeyDECRYPT(pubStr string, input []byte) ([]byte, error) {

	rsa := rasForSigature
	rsa.setPublicKey(pubStr)
	return rsa.PubKeyDECRYPT(input)

}

/*
描述 ：私钥加密--主要用于插件签名加解密
参数 ：
	priStr:私钥
	input:要加密的源数据
返回值：
	加密的目标数据、及错误码
*/
func PriKeyENCTYPT(priStr string, input []byte) ([]byte, error) {

	rsa := rasForSigature
	rsa.setPrivateKey(priStr)
	return rsa.PriKeyENCTYPT(input)
}
