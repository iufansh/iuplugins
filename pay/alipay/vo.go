package alipay

type BaseParam struct {
	AppId        string `json:"app_id"`
	Method       string `json:"method"`
	Format       string `json:"format,omitempty"`
	Charset      string `json:"charset"`
	SignType     string `json:"sign_type"`
	Sign         string `json:"sign"`
	Timestamp    string `json:"timestamp"`
	Version      string `json:"version"`
	AppAuthToken string `json:"app_auth_token,omitempty"`
	BizContent   string `json:"biz_content"`
}

type OrderSettleParam struct {
	OutRequestNo      string                       `json:"out_request_no"`
	TradeNo           string                       `json:"trade_no"`
	RoyaltyParameters []OrderSettleParamDetailInfo `json:"royalty_parameters"`
	OperatorId        string                       `json:"operator_id,omitempty"`
}

type OrderSettleParamDetailInfo struct {
	TransOut         string  `json:"trans_out,omitempty"`
	TransIn          string  `json:"trans_in,omitempty"`
	Amount           float64 `json:"amount,omitempty"`
	AmountPercentage int     `json:"amount_percentage,omitempty"`
	Desc             string  `json:"desc,omitempty"`
}

type OrderSettleResp struct {
	Data struct {
		Code    string `json:"code"`
		Msg     string `json:"msg"`
		SubCode string `json:"sub_code"`
		SubMsg  string `json:"sub_msg"`
		TradeNo string `json:"trade_no"`
	} `json:"alipay_trade_order_settle_response"`
	Sign string `json:"sign"`
}

// 支付宝服务窗验证参数：alipay.service.check
type CheckServiceBizContent struct {
	AppId       string `xml:"AppId,omitempty"`
	FromUserId  string `xml:"FromUserId,omitempty"`
	CreateTime  string `xml:"CreateTime,omitempty"`
	MsgType     string `xml:"MsgType,omitempty"`
	EventType   string `xml:"EventType,omitempty"`
	ActionParam string `xml:"ActionParam,omitempty"`
	AgreementId string `xml:"AgreementId,omitempty"`
	AccountNo   string `xml:"AccountNo,omitempty"`
}

type AlipayNoticeParam struct {
	NotifyTime     string
	NotifyType     string
	NotifyId       string
	SignType       string
	Sign           string
	TradeNo        string
	AppId          string
	OutTradeNo     string
	OutBizNo       string
	BuyerId        string
	BuyerLogonId   string
	SellerId       string
	SellerEmail    string
	TradeStatus    string
	TotalAmount    float32
	ReceiptAmount  float32
	InvoiceAmount  float32
	BuyerPayAmount float32
	PointAmount    float32
	RefundFee      float32
	SendBackFee    float32
	Subject        string
	Body           string
	GmtCreate      string
	GmtPayment     string
	GmtRefund      string
	GmtClose       string
	FundBillList   string
}
