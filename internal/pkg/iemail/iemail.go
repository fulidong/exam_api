package iemail

import (
	"bytes"
	_const "exam_api/internal/const"
	"text/template"
)

type EmailData struct {
	Name         string
	CompanyName  string
	ExamName     string
	ExamURL      string
	Username     string
	Password     string
	Duration     string
	ContactName  string
	ContactEmail string
	ContactPhone string
	SendDate     string
}

func RenderEmail(data EmailData) (string, error) {
	// 解析模板
	tmpl, err := template.New("email").Parse(_const.EmailTemplate)
	if err != nil {
		return "", err
	}

	// 渲染模板
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
