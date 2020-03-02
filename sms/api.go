package sms

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	utils "github.com/iufansh/iutils"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type SmsParam struct {
	Api      string
	Uid      string
	Key      string
	SignName string
	Mobile   string
	Text     string
}

func SendSms(p SmsParam) (num int64, err error) {
	switch p.Api {
	case "1":
		return SendWebChinese(p)
	default:
		return SendAliyunSms(p)
	}
	return
}

func SendAliyunSms(sender SmsParam) (num int64, err error) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", sender.Uid, sender.Key)

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = sender.Mobile
	request.SignName = sender.SignName
	request.TemplateCode = sender.Api
	request.TemplateParam = `{"code":"` + sender.Text + `"}`

	response, err := client.SendSms(request)
	if err != nil {
		return 0, err
	}
	if response.Code == "OK" {
		num = 1
		return
	}
	return 0, errors.New(response.Message)
}

func SendWebChinese(sender SmsParam) (num int64, err error) {
	if sender.Uid == "" || sender.Key == "" || sender.Mobile == "" || sender.Text == "" {
		err = errors.New("Params empty")
		return
	}

	smsText := fmt.Sprintf("【%s】验证码：%s，请勿泄露于他人，如非本人操作，请忽略本短信。", sender.SignName, sender.Text)
	reqUrl := fmt.Sprintf("http://utf8.api.smschinese.cn/?Uid=%s&KeyMD5=%s&smsMob=%s&smsText=%s", sender.Uid, strings.ToUpper(utils.Md5(sender.Key)), sender.Mobile, smsText)
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
