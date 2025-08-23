package model

import (
	"time"
)

// 考生模型
type Examinee struct {
	ID           int64     `json:"id"`
	SalesPaperID int64     `json:"salesPaperId"`
	ExamName     string    `json:"examName"`
	ExamineeID   int64     `json:"examineeId"`
	EmailStatus  int       `json:"emailStatus"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedBy    string    `json:"createdBy"`
	UpdatedBy    string    `json:"updatedBy"`
}
