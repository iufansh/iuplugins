package text

import (
	"fmt"
	"github.com/iufansh/iuplugins/baidu"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

type CensorParam struct {
	ApiKey      string // 必填
	SecretKey   string // 必填
	AccessToken string // 非必填，填写时，优先使用，如果过期，再使用key获取新值
	Text        string // 必填
}

type CensorResult struct {
	LogId          int64  `json:"log_id,omitempty"`
	ErrorCode      int64  `json:"error_code,omitempty"`
	ErrorMsg       string `json:"error_msg,omitempty"`
	Conclusion     string `json:"conclusion,omitempty"`
	ConclusionType int    `json:"conclusionType,omitempty"`
}

func BaiduCensor(p CensorParam) (bool, string, error) {
	var accessToken string
	var err error
	if p.ApiKey == "" || p.SecretKey == "" || p.Text == "" {
		return false, "", errors.New("Param empty")
	}
	var result CensorResult
	fmt.Println(p.Text)
	if p.AccessToken != "" {
		accessToken = p.AccessToken
		if result, err = postCheck(accessToken, p.Text); err != nil {
			return false, accessToken, err
		}
		// token 有效时，判断结果，如果无效，则跳过继续下面的步骤
		if result.ErrorCode != 110 && result.ErrorCode != 111 { // 110 Access Token失效; 111 Access token过期
			if result.ErrorCode != 0 && result.ErrorMsg != "" {
				return false, accessToken, errors.New(fmt.Sprintf("errCode:%d; errMsg:%s", result.ErrorCode, result.ErrorMsg))
			}
			if result.ConclusionType == 1 {
				return true, accessToken, nil
			}
			return false, accessToken, nil
		}
	}
	// 获取新的token
	if at, err := baidu.GetAccessToken(p.ApiKey, p.SecretKey); err != nil {
		return false, "", err
	} else {
		accessToken = at.AccessToken
	}
	// 审查内容
	if result, err = postCheck(accessToken, p.Text); err != nil {
		return false, accessToken, err
	}
	if result.ErrorCode != 0 && result.ErrorMsg != "" {
		return false, accessToken, errors.New(fmt.Sprintf("errCode:%d; errMsg:%s", result.ErrorCode, result.ErrorMsg))
	}
	if result.ConclusionType == 1 {
		return true, accessToken, nil
	}
	return false, accessToken, nil
}

func postCheck(accessToken string, text string) (CensorResult, error) {
	req := gorequest.New().Post("https://aip.baidubce.com/rest/2.0/solution/v1/text_censor/v2/user_defined")
	req = req.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.Param("access_token", accessToken)
	req = req.Param("text", text)
	var result CensorResult
	if _, _, errs := req.EndStruct(&result); len(errs) > 0 {
		return result, errs[0]
	}
	return result, nil
}
