package wechat

import (
	"fmt"
	"sort"
	"strings"
	"crypto/md5"
	"errors"
)

func WechatGenSign(key string, m map[string]interface{}) (string, error) {
	var signData []string
	for k, v := range m {
		if v != nil && v != "" && k != "sign" && k != "key" {
			signData = append(signData, fmt.Sprintf("%s=%v", k, v))
		}
	}

	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + key
	//fmt.Println("Sign str: ", signStr)
	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return "", errors.New("WechatGenSign md5.Write: " + err.Error())
	}
	signByte := c.Sum(nil)
	if err != nil {
		return "", errors.New("WechatGenSign md5.Sum: " + err.Error())
	}
	return strings.ToUpper(fmt.Sprintf("%x", signByte)), nil
}
