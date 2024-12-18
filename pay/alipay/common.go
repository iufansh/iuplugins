package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"

	utils "github.com/iufansh/iutils"
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
