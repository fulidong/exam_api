package model

import (
	"time"
)

// 登录日志模型
type LoginRecord struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"userId"`
	LoginPlatform int       `json:"loginPlatform"`
	CreatedAt     time.Time `json:"createdAt"`
}
