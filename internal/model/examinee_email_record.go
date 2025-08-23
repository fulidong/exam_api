package model

import (
	"time"
)

// 邮件记录模型
type ExamineeEmailRecord struct {
	ID                int64     `json:"id"`
	SalesPaperID      int64     `json:"salesPaperId"`
	ExamineeID        int64     `json:"examineeId"`
	Title             string    `json:"title"`
	Content           string    `json:"content"`
	ReceiverEmail     string    `json:"receiverEmail"`
	SendTime          time.Time `json:"sendTime"`
	IsSended          bool      `json:"isSended"`
	SenderEmail       string    `json:"senderEmail"`
	CopyReceiverEmail string    `json:"copyReceiverEmail"`
	Attachment        string    `json:"attachment"`
	IsFalseAddress    bool      `json:"isFalseAddress"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	CreatedBy         string    `json:"createdBy"`
	UpdatedBy         string    `json:"updatedBy"`
}
