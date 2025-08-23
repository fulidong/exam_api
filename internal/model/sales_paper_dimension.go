package model

import (
	"time"
)

// 试卷评语模型
type SalesPaperComment struct {
	ID                int64     `json:"id"`
	SalesPaperID      int64     `json:"salesPaperId"`
	Content           string    `json:"content"`
	UpScore           float64   `json:"upScore"`
	LowScore          float64   `json:"lowScore"`
	CommentCategoryID int       `json:"commentCategoryId"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	CreatedBy         string    `json:"createdBy"`
	UpdatedBy         string    `json:"updatedBy"`
}
