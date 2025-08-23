package model

import (
	"time"
)

// 试题选项模型
type QuestionOption struct {
	ID          int64     `json:"id"`
	QuestionID  int64     `json:"questionId"`
	Score       float64   `json:"score"`
	Description string    `json:"description"`
	DimensionID int64     `json:"dimensionId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedBy   string    `json:"createdBy"`
	UpdatedBy   string    `json:"updatedBy"`
}
