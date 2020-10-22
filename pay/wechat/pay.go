package wechat

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

// Pay 支付
func PayUnifiedOrder(order *WxUnifiedOrder, md5Key string) (*WxUnifiedOrderResp, error) {
	var m map[string]interface{}
	b, err := json.Marshal(&order)
	if err != nil {
		return nil, errors.New("WxUnifiedOrder marshal json error")
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, errors.New("WxUnifiedOrder Unmarshal json error")
	}

	sign, err := WechatGenSign(md5Key, m)
	if err != nil {
		return nil, err
	}
	order.Sign = sign
	b, err = xml.Marshal(&order)
	fmt.Println("order parmas:", string(b))
	if err != nil {
		return nil, errors.New("WxUnifiedOrder marshal xml error")
	}

	_, body, errs := gorequest.New().Post("https://api.mch.weixin.qq.com/pay/unifiedorder").Type("xml").Send(string(b)).End()
	if errs != nil && len(errs) > 0 {
		return nil, errors.New("WxUnifiedOrder post resp error")
	}
	var xmlResp WxUnifiedOrderResp
	err = xml.Unmarshal([]byte(body), &xmlResp)
	if err != nil {
		return nil, errors.New("WxUnifiedOrder post resp unmarshal error")
	}

	// 验签
	var m2 map[string]interface{}
	b2, err := json.Marshal(&xmlResp)
	if err != nil {
		return nil, errors.New("WxUnifiedOrder marshal resp json error2")
	}
	err = json.Unmarshal(b2, &m2)
	if err != nil {
		return nil, errors.New("WxUnifiedOrder Unmarshal resp json error2")
	}

	sign2, err := WechatGenSign(md5Key, m2)
	if err != nil {
		return nil, err
	}
	if sign2 != xmlResp.Sign {
		return nil, errors.New("WxUnifiedOrder resp sign not match")
	}

	return &xmlResp, nil
}
