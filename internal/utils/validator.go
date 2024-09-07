package utils

import "regexp"

func isValidURL(url string) bool {
	// 正则表达式模式，用于匹配 URL
	pattern := `^(https?://)[^\s/$.?#].[^\s]*$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(url)
}

func IsURL(url string) bool {
	return isValidURL(url)
}

func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

func IsPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(phone)
}

func IsPassword(password string) bool {
	pattern := `^\S{6,20}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(password)
}

func IsUsername(username string) bool {
	pattern := `^[a-zA-Z0-9_-]{4,16}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(username)
}

func IsNickname(nickname string) bool {
	pattern := `^[a-zA-Z0-9\u4e00-\u9fa5]{2,8}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(nickname)
}

func IsIDCard(idCard string) bool {
	pattern := `^(\d{15}$|^\d{18}$|^\d{17}(\d|X|x))$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(idCard)
}

func IsBankCard(bankCard string) bool {
	pattern := `^([1-9]{1})(\d{14}|\d{15}|\d{16}|\d{17}|\d{18}|\d{19})$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(bankCard)
}

func IsZipCode(zipCode string) bool {
	pattern := `^(\d{6})$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(zipCode)
}

func IsQQ(qq string) bool {
	pattern := `^[1-9][0-9]{4,9}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(qq)
}

func IsWechat(wechat string) bool {
	pattern := `^[a-zA-Z]{1}[-_a-zA-Z0-9]{5,19}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(wechat)
}

func IsAlipay(alipay string) bool {
	pattern := `^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+((\.[a-zA-Z0-9_-]{2,3}){1,2})$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(alipay)
}

func IsIDNumber(idNumber string) bool {
	pattern := `^(\d{15}$|^\d{18}$|^\d{17}(\d|X|x))$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(idNumber)
}

func IsChinese(chinese string) bool {
	pattern := `^[\u4e00-\u9fa5]+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(chinese)
}
