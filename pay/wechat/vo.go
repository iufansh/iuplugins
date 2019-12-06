package wechat

import (
	"encoding/xml"
)

// cdata支持不好，当用这种方式时，转化为map时，会被解析为map
type CDATA struct {
	Cd string `xml:",cdata"`
}

// 返回：状态和信息
type WxLoginResp struct {
	Errcode string `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// 参数：认证接口参数
type WxOauth struct {
	Appid        string
	RedirectUri  string
	ResponseType string // 可空 默认为code
	Scope        string // 可空 可选择snsapi_base、snsapi_userinfo，默认为snsapi_base
	State        string // 可空 穿透参数
}

// 参数：获取access_token接口参数
type WxAccessToken struct {
	Appid  string
	Secret string
	Code   string
}

// 返回：获取access_token结果
type WxAccessTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

// 参数：统一下单接口
type WxUnifiedOrder struct {
	XMLName        xml.Name `xml:"xml" json:"-"`
	Appid          string   `xml:"appid" json:"appid"`
	MchId          string   `xml:"mch_id" json:"mch_id"`
	NonceStr       string   `xml:"nonce_str" json:"nonce_str"`               //随机字符串，长度要求在32位以内
	Sign           string   `xml:"sign" json:"sign"`                         // 传参时不需要，由支付接口赋值
	Body           string    `xml:"body" json:"body"`                         // 商品简单描述
	OutTradeNo     string   `xml:"out_trade_no" json:"out_trade_no"`         //商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
	TotalFee       int      `xml:"total_fee" json:"total_fee"`               //订单总金额，单位为分
	SpbillCreateIp string   `xml:"spbill_create_ip" json:"spbill_create_ip"` //APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP
	NotifyUrl      string   `xml:"notify_url" json:"notify_url"`             //异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数
	TradeType      string   `xml:"trade_type" json:"trade_type"`             //JSAPI:公众号支付 NATIVE:扫码支付 APP:APP支付
	// 以下参数根据条件必填
	Openid     string `xml:"openid" json:"openid"`  // trade_type=JSAPI时（即公众号支付），此参数必传
	ProductId  string `xml:"product_id" json:"product_id"` // trade_type=NATIVE时（即扫码支付），此参数必传。此参数为二维码中包含的商品ID
	// 以下参数非必填
	DeviceInfo string `xml:"device_info" json:"device_info"` // PC网页或公众号内支付可以传"WEB"
	SignType   string `xml:"sign_type" json:"sign_type"`
	Detail     string  `xml:"detail" json:"detail"`
	Attach     string  `xml:"attach" json:"attach"`
	FeeType    string `xml:"fee_type" json:"fee_type"`
	TimeStart  string `xml:"time_start" json:"time_start"`
	TimeExpire string `xml:"time_expire" json:"time_expire"`
	GoodsTag   string `xml:"goods_tag" json:"goods_tag"`
	LimitPay   string `xml:"limit_pay" json:"limit_pay"`
	SceneInfo  string  `xml:"scene_info" json:"scene_info"`  // h5时使用，格式{"h5_info": {"type":"Wap","wap_url": "https://pay.qq.com","wap_name": "腾讯充值"}}
}

// 返回：微信返回基本信息
type WxBaseResp struct {
	ReturnCode string `xml:"return_code"` // SUCCESS/FAIL
	ReturnMsg  string `xml:"return_msg"`
}

// 返回：以下信息在WxBaseResp的return_code为SUCCESS的时候有返回
type WxResultResp struct {
	Appid      string `xml:"appid,omitempty"`
	MchId      string `xml:"mch_id,omitempty"`
	DeviceInfo string `xml:"device_info,omitempty"`
	NonceStr   string `xml:"nonce_str,omitempty"`
	Sign       string `xml:"sign,omitempty"`
	ResultCode string `xml:"result_code,omitempty"` // SUCCESS/FAIL
	ErrCode    string `xml:"err_code,omitempty"`
	ErrCodeDes string `xml:"err_code_des,omitempty"`
}

// 返回：以下字段在WxBaseResp的return_code 和WxResultResp的result_code都为SUCCESS的时候有返回
type WxSuccessResp struct {
	TradeType string `xml:"trade_type,omitempty"`
	PrepayId  string `xml:"prepay_id,omitempty"`
	CodeURL   string `xml:"code_url,omitempty"` // 公众号支付时返回
	MwebUrl   string `xml:"mweb_url,omitempty"` // h5支付时返回
}

// 返回：微信统一下单接口返回
type WxUnifiedOrderResp struct {
	WxBaseResp
	WxResultResp
	WxSuccessResp
}

// 返回：微信jsapi参数
type WxJsApiParams struct {
	Appid     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

// 返回：微信统一下单接口回调参数
type WxUnifiedOrderCallback struct {
	WxBaseResp
	// 以下信息在WxBaseResp的return_code为SUCCESS的时候有返回
	Appid              string `xml:"appid,omitempty"`
	MchId              string `xml:"mch_id,omitempty"`
	DeviceInfo         string `xml:"device_info,omitempty"`
	NonceStr           string `xml:"nonce_str,omitempty"`
	Sign               string `xml:"sign,omitempty"`
	SignType           string `xml:"sign_type,omitempty"`
	ResultCode         string `xml:"result_code,omitempty"` // SUCCESS/FAIL
	ErrCode            string `xml:"err_code,omitempty"`
	ErrCodeDes         string `xml:"err_code_des,omitempty"`
	Openid             string `xml:"openid,omitempty"`
	IsSubscribe        string `xml:"is_subscribe,omitempty"`
	TradeType          string `xml:"trade_type,omitempty"`
	BankType           string `xml:"bank_type,omitempty"`
	TotalFee           int    `xml:"total_fee,omitempty"`
	SettlementTotalFee int    `xml:"settlement_total_fee,omitempty"`
	FeeType            string `xml:"fee_type,omitempty"`
	CashFee            int    `xml:"cash_fee,omitempty"`
	CashFeeType        string `xml:"cash_fee_type,omitempty"`
	CouponFee          int    `xml:"coupon_fee,omitempty"`
	CouponCount        int    `xml:"coupon_count,omitempty"`
	//ErrCodeDes string `xml:"coupon_type_$n,omitempty"`
	//ErrCodeDes string `xml:"coupon_id_$n,omitempty"`
	//ErrCodeDes int `xml:"coupon_fee_$n,omitempty"`
	TransactionId string `xml:"transaction_id,omitempty"`
	OutTradeNo    string `xml:"out_trade_no,omitempty"`
	Attach        string `xml:"attach,omitempty"`
	TimeEnd       string `xml:"time_end,omitempty"`
}
