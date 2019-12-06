# 微信、支付宝支付接口

依赖：
go get github.com/shopspring/decimal
go get github.com/parnurzeal/gorequest

package main

import (
	"fmt"
	"github.com/guidao/gopay"
	"github.com/guidao/gopay/client"
	"github.com/guidao/gopay/common"
	"github.com/guidao/gopay/constant"
	"net/http"
)

//支付宝举例
func main() {
	//设置支付宝账号信息
	initClient()
	//设置回调函数
	initHandle()

	//支付
	charge := new(common.Charge)
	charge.PayMethod = constant.WECHAT                              //支付方式
	charge.MoneyFee = 1                                             // 支付钱单位分
	charge.Describe = "test pay"                                    //支付描述
	charge.TradeNum = "1111111111"                                  //交易号
	charge.CallbackURL = "http://127.0.0.1/callback/aliappcallback" //回调地址必须跟下面一样

	fdata, err := gopay.Pay(charge)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fdata)
}

func initClient() {
	client.InitAliAppClient(&client.AliAppClient{
		PartnerID:  "xxx",
		SellerID:   "xxxx",
		AppID:      "xxx",
		PrivateKey: nil,
		PublicKey:  nil,
	})
}

func initHandle() {
	http.HandleFunc("callback/aliappcallback", func(w http.ResponseWriter, r *http.Request) {
		//返回支付结果
		aliResult, err := gopay.AliAppCallback(w, r)
		if err != nil {
			fmt.Println(err)
			//log.xxx
			return
		}
		//接下来处理自己的逻辑
		fmt.Println(aliResult)
	})
}
