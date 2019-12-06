package googleauth

import (
	"strconv"
	"time"
	"encoding/base32"
	"fmt"
	"iufan/common/utils"
)

func GetGAuthQr(user string) (ok bool, secret string, qrCode string) {
	secret = utils.Md5(strconv.FormatInt(time.Now().Unix(), 10), utils.Pubsalt)
	secret = base32.StdEncoding.EncodeToString([]byte(secret))
	secret = string([]rune(secret)[0:16])
	otp := OTPConfig{
		Secret:      secret,
		HotpCounter: 0,
		WindowSize:  5,
	}
	qr := otp.ProvisionURI(user)
	fmt.Println(qr)
	qrCode, ok = utils.GenerateQrCode(&qr)
	return
}

func VerifyGAuth(secret string, authCode string) (ok bool, err error) {
	otp := OTPConfig{
		Secret:      secret,
		HotpCounter: 0,
		WindowSize:  5,
	}
	ok, err = otp.Authenticate(authCode)
	return
}
