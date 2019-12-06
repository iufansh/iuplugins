package common

type BusCallBackParam struct {
	Code          int // 0:未知,1:成功,2:失败,3:等待付款
	Msg           string
	OutTradeNo    string
	TransactionId string
	PayAmount     int
	Remark 		  string
}

type BusCallBack func(param *BusCallBackParam) (bool, string)
