package ihash

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// 密码哈希函数
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashed), nil
}
