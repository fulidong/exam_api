package model

import (
	"time"
)

// 考生答案模型
type ExamineeAnswerQuestionAnswer struct {
	ID               int64     `json:"id"`
	ExamineeAnswerID int64     `json:"examineeAnswerId"`
	QuestionID       int64     `json:"questionId"`
	QuestionAnswer   string    `json:"questionAnswer"`
	Score            float64   `json:"score"`
	OptionSign       string    `json:"optionSign"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	CreatedBy        string    `json:"createdBy"`
	UpdatedBy        string    `json:"updatedBy"`
}
