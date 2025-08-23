package model

import (
	"time"
)

// 试卷模型
type SalesPaper struct {
	ID                         int64     `json:"id"`
	Name                       string    `json:"name"`
	RecommendTimeLim           int       `json:"recommendTimeLim"`
	MaxScore                   int       `json:"maxScore"`
	MinScore                   int       `json:"minScore"`
	SalesPaperReportTemplateID int       `json:"salesPaperReportTemplateId"`
	QuestionNumPerPage         int       `json:"questionNumPerPage"`
	IsEnabled                  bool      `json:"isEnabled"`
	IsUsed                     bool      `json:"isUsed"`
	CategoryNumber             int       `json:"categoryNumber"`
	Mark                       string    `json:"mark"`
	CreatedAt                  time.Time `json:"createdAt"`
	UpdatedAt                  time.Time `json:"updatedAt"`
	CreatedBy                  string    `json:"createdBy"`
	UpdatedBy                  string    `json:"updatedBy"`
}
