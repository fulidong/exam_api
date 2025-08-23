package model

import (
	"gorm.io/gorm"
	"time"
)

type ExamineeSalesPaperAssociation struct {
	ID             string         `json:"id"`               // 主键（Guid）
	SalesPaperID   string         `json:"sales_paper_id"`   // SalesPaper表的外键
	SalesPaperName string         `json:"sales_paper_name"` // 考试名称
	ExamineeID     string         `json:"examinee_id"`      // 关联考生ID
	ExamineeName   string         `json:"examinee_name"`    // 关联考生名称
	ExamineeEmail  string         `json:"examinee_email"`   // 关联考生邮箱
	ExamineePhone  string         `json:"examinee_phone"`   // 关联考生电话
	EmailStatus    int32          `json:"email_status"`     // 邮件状态：1.未发送，2.已发送，3.发送失败
	ReportPath     string         `json:"report_path"`      // 答题报告路径
	StageNumber    int32          `json:"stage_number"`     // 阶段编号（0~4）
	CreatedAt      time.Time      `json:"created_at"`       // 创建时间
	UpdatedAt      time.Time      `json:"updated_at"`       // 更新时间
	CreatedBy      string         `json:"created_by"`       // 创建人标识
	UpdatedBy      string         `json:"updated_by"`       // 更新人标识
	DeletedAt      gorm.DeletedAt `json:"deleted_at"`       // 逻辑删除时间
}
