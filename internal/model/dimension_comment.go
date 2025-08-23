package model

import (
	"time"
)

// 维度评语模型
type DimensionComment struct {
	ID        int64     `json:"id"`
	LowScore  float64   `json:"lowScore"`
	UpScore   float64   `json:"upScore"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedBy string    `json:"updatedBy"`
}
