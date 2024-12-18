package googleauth

import (
	"encoding/base32"
	"fmt"
	"strconv"
	"time"

	utils "github.com/iufansh/iutils"
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
