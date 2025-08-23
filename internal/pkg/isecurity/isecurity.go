package isecurity

import (
	"golang.org/x/crypto/bcrypt"
)

// 简化版密码处理
const (
	DefaultCost = 12 // 合理的安全强度，可在生产环境使用
)

// HashPassword 生成安全的密码哈希
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPassword 验证密码与哈希是否匹配
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
