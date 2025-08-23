package iregexp

import (
	"regexp"
	"strings"
)

// IsValidEmail 验证邮箱格式是否合法
func IsValidEmail(email string) bool {
	// RFC 5322 官方标准正则（简化版）
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// GetEmailPrefix 截取 @ 之前的部分
func GetEmailPrefix(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[0]
	}
	return ""
}

func IsValidPhoneNumberWithCountryCode(phone string) bool {
	regex := `^1[345789]\d{9}$`
	matched, _ := regexp.MatchString(regex, phone)
	return matched
}
