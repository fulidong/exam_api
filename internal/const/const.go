package _const

const (
	AdministratorPrefix                     = "A"
	SysLoginRecordPrefix                    = "SLR"
	SalesPaperPrefix                        = "SPP"
	SalesPaperCommentPrefix                 = "SPCP"
	SalesPaperDimensionPrefix               = "SPDP"
	SalesPaperDimensionCommentPrefix        = "SPDCP"
	SalesPaperDimensionQuestionPrefix       = "SPDQP"
	SalesPaperDimensionQuestionOptionPrefix = "SPDQOP"
	ExamineePrefix                          = "EP"
	ExamineeSalesPaperAssociationPrefix     = "ESPA"
	ExamineeEmailRecordPrefix               = "EERP"
	ExamineeAnswerPrefix                    = "EAP"
	ExamineeAnswerDimensionScorePrefix      = "EADSP"
	ExamineeAnswerQuestionAnswerPrefix      = "EAQAP"
	ExamiEventPrefix                        = "EEP"
)

var AllowedVars = map[string]interface{}{
	"raw_score":     0.0,
	"average_mark":  0.0,
	"standard_mark": 0.0,
}

var VerifyExamTokenMethod = map[string]struct{}{
	"/exam_api.v1.ExamService/ExamQuestion":       struct{}{},
	"/exam_api.v1.ExamService/ExamQuestionRecord": struct{}{},
	"/exam_api.v1.ExamService/HeartbeatAndSave":   struct{}{},
	"/exam_api.v1.ExamService/SubmitExam":         struct{}{},
}

// 邮件模板
const EmailTemplate = `
尊敬的{{.Name}}您好：
感谢您应聘{{.CompanyName}}的职位！我们已为您安排在线能力测评，请您在规定时间内完成测试，以便我们继续推进后续招聘流程。
以下是您的考试信息，请妥善保存：
考试名称 ：{{.ExamName}}
考试网址 ：{{.ExamURL}}
登录账号 ：{{.Username}}
登录密码 ：{{.Password}}
📌 考试说明：
- 请您尽快完成考试。
- 考试时长 {{.Duration}} 分钟，一旦开始请尽量一次性完成。
- 请使用 Chrome 或 Edge 浏览器打开链接，确保网络稳定。
- 考试期间请保持独立作答，避免切换页面或使用其他设备。
👉 立即开始考试 ：{{.ExamURL}}
如有任何技术问题或疑问，请联系：{{.ContactName}}（邮箱：{{.ContactEmail}}，电话：{{.ContactPhone}}）。
期待您的顺利完成，祝您考试顺利！

此致
敬礼

{{.CompanyName}} 人力资源部
{{.SendDate}}
`

type ExamEventType string

const (
	ExamEventHeartbeat    ExamEventType = "heartbeat"     // 心跳
	ExamEventLongInactive ExamEventType = "long_inactive" // 长时间无心跳
	ExamEventSubmit       ExamEventType = "submit"        // 提交
	ExamEventTimeUp       ExamEventType = "time_up"       // 时间到
)
