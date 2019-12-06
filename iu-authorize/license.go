package license

import (
	"fmt"
	"strconv"
	"net"
	"iufan/common/utils"
	"strings"
	"math"
	"time"
)

const defaultSalt  = "licgo"

func GetMachineData(salts ...string) []string {
	var salt = defaultSalt
	if len(salt) > 0 {
		for _, v := range salts {
			salt += "'" + v
		}
	}
	var data = make([]string, 0)
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return data
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		s := fmt.Sprintf("%d%s%s", netInterface.Index, salt, macAddr)
		data = append(data, utils.Md5(s, utils.Pubsalt) + "(" + netInterface.Name + ")")
	}
	return data
}

func CheckLicense(lic string, expTime time.Time, netTime bool, salts ...string) (ok bool, msg string) {
	var currentTime time.Time
	var err error
	if netTime {
		currentTime, err = utils.GetNetTime()
		if err != nil {
			msg = "请确认网络连接是否正常"
			return
		}
	} else {
		currentTime = time.Now()
	}
	if currentTime.After(expTime) {
		msg = "激活已过期，请重新激活"
		return
	}

	var salt = defaultSalt
	if len(salt) > 0 {
		for _, v := range salts {
			salt += "'" + v
		}
	}

	for _, v := range GetMachineData(salts...) {
		ss := strings.Split(v, "(")
		if len(ss) < 2 {
			continue
		}
		activeCode := genActiveCode(ss[0], salt)
		//fmt.Println(activeCode, expTime.Format("20060102150405"), salt)
		licTmp := utils.Md5(activeCode, expTime.Format("20060102150405"), salt)

		if licTmp == lic {
			ok = true
			msg = "验证成功"
			return
		}
	}
	msg = "验证失败"
	return
}

func GenLicense(activeCode string, netTime bool, salts ...string) (ok bool, lic string, expTime time.Time) {
	var salt = defaultSalt
	if len(salt) > 0 {
		for _, v := range salts {
			salt += "'" + v
		}
	}

	daySign := activeCode[32:]
	var total int64
	for _, v := range []rune(activeCode[:32]) {
		total += int64(v)
	}
	if total > 9999 {
		total = 9999
	}
	//fmt.Println(total)
	num, err := strconv.ParseInt(daySign[:1], 10, 64)
	if err != nil {
		return false, "生成失败（1）", expTime
	}
	//fmt.Println(num)
	var k = num * 2 + 2
	var dayStr string
	for {
		//fmt.Println(k)
		dayStr += daySign[k-1:k]

		k -= 2
		num--
		if num <= 0 {
			break
		}
	}
	//fmt.Println(dayStr)
	days, err := strconv.ParseInt(dayStr, 10, 64)
	if err != nil {
		return false, "生成失败（2）", expTime
	}
	//fmt.Println(days)
	day := math.Floor(float64(days * 35) / float64(12 * total) + 0.5)
	//fmt.Println(day)
	var currentTime time.Time
	if netTime {
		currentTime, err = utils.GetNetTime()
		if err != nil {
			return false, "请确认网络连接是否正常", expTime
		}
	} else {
		currentTime = time.Now()
	}
	expTime = currentTime.AddDate(0, 0, int(day))

	//fmt.Println(activeCode[:32], expTime.Format("20060102150405"), salt)
	lic = utils.Md5(activeCode[:32], expTime.Format("20060102150405"), salt)

	ok = true
	return
}

func GenActiveCode(acKey string, days int64, salts ...string) string {
	var salt = defaultSalt
	if len(salt) > 0 {
		for _, v := range salts {
			salt += "'" + v
		}
	}
	activeCode := genActiveCode(acKey, salt)
	var total int64
	var seed int64
	for i, v := range []rune(activeCode) {
		total += int64(v)
		seed = seed + int64(v) * int64(i)
	}
	if total > 9999 {
		total = 9999
	}
	//fmt.Println(total)

	randStr := utils.RandStrLowerWithNumSeed(18, seed)
	//fmt.Println(randStr)

	//fmt.Println(days)
	days = days * total / 5  + days * total / 7
	//fmt.Println(days)

	var daySign string
	var k int
	for {
		k += 2
		a2 := days % 10
		days = days / 10

		randStr = randStr[:k] + strconv.FormatInt(a2, 10) + randStr[k:]

		if days <= 0 {
			break
		}
	}
	randStr = fmt.Sprintf("%d%s", k/2, randStr)
	daySign = randStr[:18]
	//fmt.Println(daySign)

	ac := activeCode + daySign
	//fmt.Println(ac)

	return ac
}

func genActiveCode(acKey string, salt string) string {
	runes := []rune(acKey)
	var a string
	for i, v := range runes {
		x := utils.AnyToDecimal(string(v), 36)
		a = a + strconv.Itoa(x+int(math.Pow(float64(x), float64(i%7))))
	}
	return utils.Md5(a+salt)
}
