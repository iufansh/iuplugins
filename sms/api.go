package sms

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"strconv"
	"github.com/pkg/errors"
	"iufan/common/utils"
	"strings"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type SmsParam struct {
	Api    string
	Uid    string
	Key    string
	Mobile string
	Text   string
}

func SendSms(p SmsParam) (num int64, err error) {
	switch p.Api {
	case "1":
		return SendWebChinese(p.Uid, p.Key, p.Mobile, p.Text)
	default:
		return 0, errors.New("No match api")
	}
	return
}

func SendWebChinese(uid string, key string, smsMob string, smsText string) (num int64, err error) {
	if uid == "" || key == "" || smsMob == "" || smsText == "" {
		err = errors.New("Params empty")
		return
	}
	reqUrl := fmt.Sprintf("http://utf8.api.smschinese.cn/?Uid=%s&KeyMD5=%s&smsMob=%s&smsText=%s", uid, strings.ToUpper(utils.Md5(key)), smsMob, smsText)
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return
	}
	req.Close = true

	var resp *http.Response
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return
	} else {
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		body := string(bodyBytes)
		fmt.Println(body)
		if num, err = strconv.ParseInt(body, 10, 64); err != nil {
			return
		} else if num <= 0 {
			err = errors.New("发送失败，接口返回：" + body)
		}
	}
	return
}

func SendAliyunSms(accessKeyId string, accessSecret string, smsMob string, smsText string) (num int64, err error) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessSecret)

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = smsMob
	request.SignName = "xx公司"
	request.TemplateCode = "tempid1"

	response, err := client.SendSms(request)
	if err != nil {
		return
	}
	if response.Code == "OK" {
		num = 1
		return
	}
	return 0, errors.New(response.Message)
}
