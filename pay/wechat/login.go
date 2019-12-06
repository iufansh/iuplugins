package wechat

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"errors"
	"encoding/json"
	"net/http"
	"net/url"
)

// 微信授权请求地址
// Scope默认 snsapi_base
// 文档地址：https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140842
func GetOauthUrl(vo *WxOauth) string {
	baseUrl := "https://open.weixin.qq.com/connect/oauth2/authorize"
	if vo.ResponseType == "" {
		vo.ResponseType = "code"
	}
	if vo.Scope == "" {
		vo.Scope = "snsapi_base"
	}
	baseUrl = fmt.Sprintf("%s?appid=%s&redirect_uri=%s&response_type=%s&scope=%s&state=%s#wechat_redirect", baseUrl, vo.Appid, url.QueryEscape(vo.RedirectUri), vo.ResponseType, vo.Scope, vo.State)
	fmt.Println(baseUrl)
	return baseUrl
}

// 解析授权返回 code 、 state
func GetOauthResult(r *http.Request) (string, string) {
	f := r.Form
	return f.Get("code"), f.Get("state")
}

func GetAccessToken(vo *WxAccessToken) (*WxAccessTokenResp, error) {
	baseUrl := "https://api.weixin.qq.com/sns/oauth2/access_token"
	urlstr := fmt.Sprintf("%s?appid=%s&secret=%s&code=%s&grant_type=authorization_code", baseUrl, vo.Appid, vo.Secret, vo.Code)
	_, body, errs := gorequest.New().Get(urlstr).End()
	if errs != nil {
		return nil, errors.New("Wechat get access_token request err")
	}
	var wxResp WxAccessTokenResp
	err := json.Unmarshal([]byte(body), &wxResp)
	if err == nil {
		return &wxResp, nil
	}
	var resultVo WxLoginResp
	err = json.Unmarshal([]byte(body), &resultVo)
	if err != nil {
		return nil, errors.New("Wechat get access_token json unmarshal err")
	}
	return nil, errors.New(fmt.Sprintf("%s:%s", resultVo.Errcode, resultVo.Errmsg))
}
