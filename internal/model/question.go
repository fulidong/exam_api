package model

import (
	"time"
)

// 试题模型
type Question struct {
	ID             int64     `json:"id"`
	DimensionID    int64     `json:"dimensionId"`
	Title          string    `json:"title"`
	Remark         string    `json:"remark"`
	Status         bool      `json:"status"`
	QuestionTypeID int       `json:"questionTypeId"`
	Order          int       `json:"order"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedBy      string    `json:"createdBy"`
	UpdatedBy      string    `json:"updatedBy"`
}
