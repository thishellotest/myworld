package biz

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

func IsValidUSAPhoneNumber(phone string) bool {
	// 正则表达式：匹配 3-3-4 的数字格式
	re := regexp.MustCompile(`^\d{3}-\d{3}-\d{4}$`)
	return re.MatchString(phone)
}

// FormatPhoneNumber +13109719619 转为300-971-9619 、(300) 971-9619、3009719619
func FormatPhoneNumber(phone string) (string, string, string, error) {
	// 正则提取数字部分
	re := regexp.MustCompile(`\d+`)
	digits := re.FindAllString(phone, -1)
	fullNumber := ""
	for _, d := range digits {
		fullNumber += d
	}

	// 确保号码长度符合格式要求
	if len(fullNumber) != 11 || fullNumber[:1] != "1" {
		return "", "", "", errors.New("Phone Format is wrong")
	}

	// 去掉国家码
	number := fullNumber[1:]

	// 生成两种格式
	format1 := fmt.Sprintf("%s-%s-%s", number[:3], number[3:6], number[6:])
	format2 := fmt.Sprintf("(%s) %s-%s", number[:3], number[3:6], number[6:])

	return format1, format2, number, nil
}

func FormatUSAPhoneHandle(phone string) (string, error) {
	phone, err := USAPhoneHandle(phone)
	if err != nil {
		return "", err
	}
	return "+1" + phone, nil
}

// (402) 215-6064： 转为：4022156064
// 904-415-7090 转为：9044157090
// USAPhoneHandle 美国手机号码格式转换
func USAPhoneHandle(phone string) (string, error) {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, "-", "")
	if len(phone) == 10 {
		return phone, nil
	}
	return phone, errors.New(phone + " format is wrong")
}

func IsUSAPhone(phone string) bool {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, "-", "")
	if len(phone) == 10 {
		return true
	}
	return false
}
