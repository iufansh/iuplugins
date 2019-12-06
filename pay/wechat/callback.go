package wechat

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/xml"
	"iufan/common/utils"
)

func CallbackUnifiedOrder(w http.ResponseWriter, r *http.Request, getWxMd5Key GetWxMd5Key, busCallback BusinessCallback) error {
	var returnCode = "FAIL"
	var returnMsg = ""
	var reXML WxUnifiedOrderCallback
	defer func() {
		// 验证通过后才进行业务回调
		if returnCode == "SUCCESS" {
			if isOk, busMsg := busCallback(&reXML); !isOk {
				returnCode = "FAIL"
				if busMsg != "" {
					returnMsg = busMsg
				}
			}
		}
		formatStr := `<xml><return_code><![CDATA[%s]]></return_code><return_msg>![CDATA[%s]]</return_msg></xml>`
		returnBody := fmt.Sprintf(formatStr, returnCode, returnMsg)
		w.Write([]byte(returnBody))
	}()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		returnCode = "FAIL"
		returnMsg = "Bodyerror"
		return errors.New("CallbackUnifiedOrder read body error")
	}
	err = xml.Unmarshal(body, &reXML)
	if err != nil {
		returnMsg = "Paramserror"
		returnCode = "FAIL"
		return errors.New("CallbackUnifiedOrder xml unmarshal error")
	}

	if reXML.ReturnCode != "SUCCESS" {
		returnCode = "FAIL"
		return errors.New(reXML.ReturnMsg)
	}
	m := utils.XmlToMap(body)

	mySign, err := WechatGenSign(getWxMd5Key(reXML.OutTradeNo), m)
	if err != nil {
		return err
	}

	if mySign != m["sign"] {
		return errors.New("CallbackUnifiedOrder check sign false error")
	}

	returnCode = "SUCCESS"
	return nil
}

// 获取MD5秘钥
type GetWxMd5Key func(orderNo string) string
// 业务状态回调
type BusinessCallback func(unifiedOrderCallback *WxUnifiedOrderCallback) (bool, string)
