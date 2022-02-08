package encryption

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"
	"ubx/common/utils"
)

func TestPlugins(t *testing.T) {
	var info SignatureInfo
	info.Random = string(utils.GenerateRandomString(32))
	info.Sha256 = "c58e8fb8e2b46a7a7501e5dff87db33b2e88076a3fa7505f4cc744cde4a6cb1e"
	info.TimeStamp = fmt.Sprintf("%d", time.Now().Unix())

	sign, err := Sha256Signature(info)
	if err != nil {
		fmt.Printf("Sha256Signature error %v\n", err)
	}

	fmt.Printf("sign = %v\n", sign)

	strEnctypt, err := PriKeyENCTYPT(pirvatekey, []byte(sign))
	if nil != err {
		fmt.Printf("PriKeyENCTYPT error %v\n", err)
	}
	fmt.Printf("strEnctypt = %v\n", strEnctypt)

	signature := base64.StdEncoding.EncodeToString(strEnctypt)
	fmt.Printf("signature = %v\n", signature)

}
