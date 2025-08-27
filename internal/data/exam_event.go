package data

import (
	"context"
	"exam_api/internal/biz"
	"exam_api/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
)

type ExamEventRepo struct {
	data *Data
	log  *log.Helper
}

func NewExamEventRepo(data *Data, logger log.Logger) biz.ExamEventRepo {
	return &ExamEventRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *ExamEventRepo) ExamEvent(ctx context.Context, examEvent *entity.ExamEvent) error {
	return r.data.db.WithContext(ctx).Create(examEvent).Error
}
