package alipay

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/iufansh/iuplugins/pay/common"
	utils "github.com/iufansh/iutils"
)

func CallbackAlipay(w http.ResponseWriter, r *http.Request, getKey GetAlipayPubKey, busCallback common.BusCallBack) error {
	var m = make(map[string]string)
	var res = common.BusCallBackParam{Code: 0}
	defer func() {
		// 业务回调
		if isOk, _ := busCallback(&res); isOk {
			w.Write([]byte("success"))
		} else {
			w.Write([]byte("fail"))
		}
	}()
	r.ParseForm()
	var signSlice []string
	for k, v := range r.Form {
		// k不会有多个值的情况
		m[k] = v[0]
		if k == "sign" || k == "sign_type" {
			continue
		}
		decValue, _ := url.QueryUnescape(v[0])
		signSlice = append(signSlice, fmt.Sprintf("%s=%s", k, decValue))
	}
	fmt.Println(m)
	res.OutTradeNo = m["out_trade_no"]
	sort.Strings(signSlice)
	if m["sign_type"] != "RSA2" {
		res.Code = 2
		res.Msg = "加密类型不是RSA2"
		return errors.New("alipay notice sign type not rsa2")
	}
	signByte, err := base64.StdEncoding.DecodeString(m["sign"])
	if err != nil {
		res.Code = 2
		res.Msg = "签名base64转化失败"
		return err
	}
	s := sha256.New()
	signData := strings.Join(signSlice, "&")
	_, err = s.Write([]byte(signData))
	if err != nil {
		res.Code = 2
		res.Msg = "签名sha256转化失败"
		return err
	}
	hash := s.Sum(nil)
	publicKey := utils.RsaParsePublicKey(getKey(m["out_trade_no"]))
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, signByte)
	if err != nil {
		res.Code = 2
		res.Msg = "验签失败"
		return err
	}

	switch m["trade_status"] {
	case "TRADE_SUCCESS":
		res.Code = 1
		res.Msg = "交易成功"
		res.TransactionId = m["trade_no"]
		res.Remark = m["buyer_id"] + "," + m["seller_id"]
	case "WAIT_BUYER_PAY":
		res.Code = 3
		res.Msg = "等待付款"
	default:
		res.Code = 2
		res.Msg = "交易失败"
	}
	return nil
}

// 获取公钥
type GetAlipayPubKey func(orderNo string) string
