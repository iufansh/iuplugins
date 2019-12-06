package alipay

import (
	"github.com/astaxie/beego"
	"github.com/parnurzeal/gorequest"
	"encoding/json"
	"net/http"
	"net/url"
	"fmt"
	"sort"
	"iufan/common/utils"
	"encoding/base64"
	"crypto/sha256"
	"strings"
	"encoding/xml"
	"crypto/rsa"
	"crypto"
	"bytes"
	"io"
	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
)

func CheckService(w http.ResponseWriter, r *http.Request, getAliPubKeyAndDepPriKey func(appId string) (string, string), getDepPubKey func(appId string) string) error {
	var m = make(map[string]string)
	r.ParseForm()
	var signSlice []string
	for k, v := range r.Form {
		// k不会有多个值的情况
		m[k] = v[0]
		if k == "sign" {
			continue
		}
		decValue, _ := url.QueryUnescape(v[0])
		signSlice = append(signSlice, fmt.Sprintf("%s=%s", k, decValue))
	}
	fmt.Println(m)
	var bizContent CheckServiceBizContent
	decoder := xml.NewDecoder(bytes.NewReader([]byte(m["biz_content"])))
	decoder.CharsetReader = func(c string, i io.Reader) (io.Reader, error) {
		return charset.NewReaderLabel(c, i)
	}
	decoder.Decode(&bizContent)
	if err := decoder.Decode(&bizContent); bizContent.AppId == "" && err != nil {
		return errors.New("biz_content 解析失败"+err.Error())
	}

	sort.Strings(signSlice)
	if m["sign_type"] != "RSA2" {
		return errors.New("sign_type not RSA2")
	}
	signByte, err := base64.StdEncoding.DecodeString(m["sign"])
	if err != nil {
		return errors.New("sign base64 decode err:"+err.Error())
	}
	s := sha256.New()
	signData := strings.Join(signSlice, "&")
	fmt.Println(signData)
	_, err = s.Write([]byte(signData))
	if err != nil {
		return errors.New("sha256 err:"+err.Error())
	}
	hash := s.Sum(nil)
	aliPubKey, depPriKey := getAliPubKeyAndDepPriKey(bizContent.AppId)
	publicKey := utils.RsaParsePublicKey(aliPubKey)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, signByte)
	if err != nil {
		return errors.New("sign verify err:"+err.Error())
	}
	// 验签通过，返回消息
	response := fmt.Sprintf(`<biz_content>%s</biz_content><success>true</success>`, getDepPubKey(bizContent.AppId))
	sign := GenStrRsaSha256Sign(depPriKey, response)

	res := `<?xml version="1.0" encoding="GBK"?>
			<alipay>
				<response>
					%s
				</response>
				<sign>%s</sign>
				<sign_type>RSA2</sign_type>
			</alipay>`
	res = fmt.Sprintf(res, response, sign)
	fmt.Println(res)
	w.Write([]byte(res))
	return nil
}

func OrderSettle(baseParam BaseParam, settleParam OrderSettleParam, priKey string) (code int, msg string, tradeNo string) {
	sp, _ := json.Marshal(settleParam)
	baseParam.BizContent = string(sp)
	bp, _ := json.Marshal(baseParam)
	var bpMap map[string]string
	json.Unmarshal(bp, &bpMap)

	bpMap["sign"] = GenRsaSha256Sign(priKey, bpMap)

	payUrl := ToURL("https://openapi.alipay.com/gateway.do", bpMap)
	beego.Info("alipay OrderSettle url：", payUrl)
	_, body, errs := gorequest.New().Get(payUrl).End()
	if errs != nil && len(errs) > 0 {
		msg = "分账请求异常1"
		beego.Error("alipay OrderSettle err:", errs)
		return
	}

	beego.Info("alipay OrderSettle response:", body)
	var res OrderSettleResp
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		msg = "分账请求异常2"
		beego.Info("alipay unmarshal json err:", err)
		return
	}
	if res.Data.Code == "10000" && res.Data.TradeNo != "" {
		code = 1
		msg = "分账成功"
		tradeNo = res.Data.TradeNo
		return
	}
	msg = res.Data.Msg + "-" + res.Data.SubMsg
	return
}

