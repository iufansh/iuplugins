package sms

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/sms"
	"github.com/baidubce/bce-sdk-go/services/sms/api"
	utils "github.com/iufansh/iutils"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type SmsParam struct {
	RegionId string
	Api      string
	Uid      string
	Key      string
	SignName string
	Mobile   string
	Text     string
}

func SendSms(p SmsParam) (num int64, err error) {
	if strings.HasPrefix(p.Api, "SMS_") {
		return SendAliyunSms(p)
	} else if strings.HasPrefix(p.Api, "sms-tmpl-") {
		return SendBaiduSms(p)
	} else if p.Api == "1" {
		return SendWebChinese(p)
	}
	return 0, errors.New("No match api")
}

func SendBaiduSms(sender SmsParam) (num int64, err error) {
	AK, SK := sender.Uid, sender.Key
	ENDPOINT := "https://smsv3.bj.baidubce.com"
	client, _ := sms.NewClient(AK, SK, ENDPOINT)

	// 配置不进行重试，默认为Back Off重试
	client.Config.Retry = bce.NewNoRetryPolicy()

	// 配置连接超时时间为30秒
	client.Config.ConnectionTimeoutInMillis = 30 * 1000

	contentMap := make(map[string]interface{})
	contentMap["code"] = sender.Text
	sendSmsArgs := &api.SendSmsArgs{
		Mobile:      sender.Mobile,
		Template:    sender.Api,
		SignatureId: sender.SignName,
		ContentVar:  contentMap,
		//ClientToken:
	}
	result, err := client.SendSms(sendSmsArgs)
	if err != nil {
		return 0, err
	}
	if result.Code == "1000" {
		num = 1
		return
	}
	return 0, errors.New(result.Message)
}

func SendAliyunSms(sender SmsParam) (num int64, err error) {
	if sender.RegionId == "" {
		sender.RegionId = "cn-hangzhou" // 默认杭州
	}
	client, err := dysmsapi.NewClientWithAccessKey(sender.RegionId, sender.Uid, sender.Key)

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
