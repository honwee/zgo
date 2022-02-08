package encryption

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strings"
)

type SignatureInfo struct {
	Sha256    string `json:"sha256"`
	TimeStamp string `json:"timestamp"`
	Random    string `json:"random"`
}

/*
描述 ：获取签名值（key升序排序后，根据url规则生成字符窜，进行sha256取值）
参数 ：
	obj:源数据
返回值：
	签名结果及错误吗
*/
func Sha256Signature(obj interface{}) (string, error) {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]string)
	for i := 0; i < obj1.NumField(); i++ {

		// fmt.Println(strings.ToLower(obj1.Field(i).Name))
		data[strings.ToLower(obj1.Field(i).Name)] = obj2.Field(i).String()
	}

	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var p = url.Values{}

	for _, k := range keys {
		p.Add(k, data[k])
	}

	str := p.Encode()

	fmt.Printf("str=%v\n", str)
	retStr, err := HashSHA256String(str)

	return retStr, err
}
