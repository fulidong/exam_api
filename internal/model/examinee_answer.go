package model

import (
	"time"
)

// 考生答题模型
type ExamineeAnswer struct {
	ID                  int64     `json:"id"`
	SalesPaperID        int64     `json:"salesPaperId"`
	ExamineeID          int64     `json:"examineeId"`
	Score               float64   `json:"score"`
	BeginTestTime       time.Time `json:"beginTestTime"`
	SubmitTime          time.Time `json:"submitTime"`
	CompleteQuestionNum int       `json:"completeQuestionNum"`
	LastActionTime      time.Time `json:"lastActionTime"`
	Comparability       int       `json:"comparability"`
	ReportPath          string    `json:"reportPath"`
	Deadline            time.Time `json:"deadline"`
	Usability           int8      `json:"usability"`
	IsCompleted         bool      `json:"isCompleted"`
	IsReaded            bool      `json:"isReaded"`
	RemainingTimelimit  int       `json:"remainingTimelimit"`
	HasPDF              bool      `json:"hasPdf"`
	StageNumber         int8      `json:"stageNumber"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	CreatedBy           string    `json:"createdBy"`
	UpdatedBy           string    `json:"updatedBy"`
}
