package model

import (
	"time"
)

// 维度得分模型
type ExamineeAnswerDimensionScore struct {
	ID                     int64     `json:"id"`
	ExamineeAnswerID       int64     `json:"examineeAnswerId"`
	DimensionID            int64     `json:"dimensionId"`
	DimensionRawScore      float64   `json:"dimensionRawScore"`
	DimensionStandardScore float64   `json:"dimensionStandardScore"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
	CreatedBy              string    `json:"createdBy"`
	UpdatedBy              string    `json:"updatedBy"`
}
