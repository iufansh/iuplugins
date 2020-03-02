package alipay

import (
	"fmt"
	"sort"
	"strings"
	"crypto"
	"net/url"
	"encoding/base64"
	"crypto/sha256"
	"crypto/rand"
	utils "github.com/iufansh/iutils"
	"github.com/astaxie/beego"
)

// GenSign 产生签名
func GenRsaSha256Sign(priKey string, m map[string]string) string {
	var data []string
	for k, v := range m {
		if k != "sign" || v != "" {
			data = append(data, fmt.Sprintf(`%s=%s`, k, v))
		}
	}
	sort.Strings(data)
	signData := strings.Join(data, "&")
	beego.Info(signData)
	s := sha256.New()
	_, err := s.Write([]byte(signData))
	if err != nil {
		return ""
	}
	hashByte := s.Sum(nil)
	privateKey := utils.RsaParsePrivateKey(priKey)
	signByte, err := privateKey.Sign(rand.Reader, hashByte, crypto.SHA256)
	if err != nil {
		return ""
	}
	return url.QueryEscape(base64.StdEncoding.EncodeToString(signByte))
}
// GenSign 产生签名
func GenStrRsaSha256Sign(priKey string, signData string) string {
	beego.Info(signData)
	s := sha256.New()
	_, err := s.Write([]byte(signData))
	if err != nil {
		return ""
	}
	hashByte := s.Sum(nil)
	privateKey := utils.RsaParsePrivateKey(priKey)
	signByte, err := privateKey.Sign(rand.Reader, hashByte, crypto.SHA256)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(signByte)
}

// ToURL
func ToURL(payUrl string, m map[string]string) string {
	var buf []string
	for k, v := range m {
		buf = append(buf, fmt.Sprintf("%s=%s", k, v))
	}
	return fmt.Sprintf("%s?%s", payUrl, strings.Join(buf, "&"))
}
