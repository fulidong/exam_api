package model

import "time"

// 管理员用户模型
type Administrator struct {
	ID           int64     `json:"id"`
	UserName     string    `json:"userName"`
	LoginAccount string    `json:"loginAccount"`
	Password     string    `json:"-"`
	Status       int       `json:"status"`
	Email        string    `json:"email"`
	UserType     int       `json:"userType"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
