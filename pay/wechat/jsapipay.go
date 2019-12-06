package wechat

import (
	"strconv"
	"time"
	"iufan/common/utils"
	"encoding/json"
	"errors"
)

// 根据统一下单返回数据，获取jsapi参数
func GetJsApiParams(md5Key string, unifiedOrderResp *WxUnifiedOrderResp) (*WxJsApiParams, error)  {
	params := WxJsApiParams{}
	params.Appid = unifiedOrderResp.Appid
	params.TimeStamp = strconv.FormatInt(time.Now().Unix(), 10)
	params.NonceStr = utils.RandString(5)
	params.Package = "prepay_id=" + unifiedOrderResp.PrepayId
	params.SignType = "MD5"
	var m map[string]interface{}
	b, err := json.Marshal(&params)
	if err != nil {
		return nil, errors.New("GetJsApiParams marshal json error")
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, errors.New("WxUnifiedOrder Unmarshal json error")
	}
	sign ,err := WechatGenSign(md5Key,m)
	if err != nil {
		return nil, err
	}
	params.PaySign = sign
	return &params, nil
}
