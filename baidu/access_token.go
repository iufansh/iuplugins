package baidu

import (
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

type BaiduAccessToken struct {
	ExpiresIn        int    `json:"expires_in,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`  // 忽略
	SessionKey       string `json:"session_key,omitempty"`    // 忽略
	Scope            string `json:"scope,omitempty"`          // 忽略
	SessionSecret    string `json:"session_secret,omitempty"` // 忽略
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func GetAccessToken(clientId, clientSecret string) (accessToken BaiduAccessToken, err error) {
	req := gorequest.New().Post("https://aip.baidubce.com/oauth/2.0/token")
	req = req.Param("grant_type", "client_credentials")
	req = req.Param("client_id", clientId).Param("client_secret", clientSecret)
	if _, _, errs := req.EndStruct(&accessToken); len(errs) > 0 {
		err = errs[0]
		return
	}
	if accessToken.Error != "" {
		return accessToken, errors.New("errCode:" + accessToken.Error + ";desc:" + accessToken.ErrorDescription)
	}
	return
}
