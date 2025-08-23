package model

import (
	"time"
)

// 试卷维度关系模型
type SalesPaperDimension struct {
	ID              int64     `json:"id"`
	SalesPaperID    int64     `json:"salesPaperId"`
	DimensionID     int64     `json:"dimensionId"`
	Weight          int       `json:"weight"`
	SelfDefineScore int       `json:"selfDefineScore"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	CreatedBy       string    `json:"createdBy"`
	UpdatedBy       string    `json:"updatedBy"`
}
