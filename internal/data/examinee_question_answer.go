package data

import (
	"context"
	"exam_api/internal/biz"
	"exam_api/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ExamineeQuestionAnswerRepo struct {
	data *Data
	log  *log.Helper
}

func NewExamineeQuestionAnswerRepo(data *Data, logger log.Logger) biz.ExamineeQuestionAnswerRepo {
	return &ExamineeQuestionAnswerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *ExamineeQuestionAnswerRepo) GetByExamineeAnswerId(ctx context.Context, examineeAnswerId string) (list []*entity.ExamineeAnswerQuestionAnswer, err error) {
	err = r.data.db.WithContext(ctx).Model(&entity.ExamineeAnswerQuestionAnswer{}).Where(" examinee_answer_id = ? ", examineeAnswerId).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ExamineeQuestionAnswerRepo) SaveAnswer(ctx context.Context, answers []*entity.ExamineeAnswerQuestionAnswer) error {
	return r.data.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "examinee_answer_id"}, {Name: "question_id"}}, // 唯一索引字段
		DoUpdates: clause.Assignments(map[string]interface{}{
			"question_answer": clause.Column{Table: "", Name: "question_answer"},
			"updated_at":      gorm.Expr("VALUES(`updated_at`)"),
		}),
	}).Create(&answers).Error
}
