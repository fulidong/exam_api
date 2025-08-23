package model

import (
	"time"
)

// 维度模型
type Dimension struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	AverageMark      float64   `json:"averageMark"`
	StandardMark     float64   `json:"standardMark"`
	Description      string    `json:"description"`
	CreateUserID     int64     `json:"createUserId"`
	CreateTime       time.Time `json:"createTime"`
	MaxScore         int       `json:"maxScore"`
	MinScore         int       `json:"minScore"`
	IsLeaf           bool      `json:"isLeaf"`
	QuestionUIModeID int       `json:"questionUiModeId"`
	Type             int       `json:"type"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	CreatedBy        string    `json:"createdBy"`
	UpdatedBy        string    `json:"updatedBy"`
}
