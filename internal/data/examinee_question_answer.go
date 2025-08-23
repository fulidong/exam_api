package data

import (
	"context"
	"exam_api/internal/biz"
	"exam_api/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
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

func (r *ExamineeQuestionAnswerRepo) GetByAssociationId(ctx context.Context, associationId string) (resEntity *entity.ExamineeAnswerQuestionAnswer, err error) {
	resEntity, err = getSingleRecordByScope[entity.ExamineeAnswerQuestionAnswer](
		r.data.db.WithContext(ctx).Model(resEntity).Where(" examinee_sales_paper_association_id = ? ", associationId),
	)
	if err != nil {
		return nil, err
	}
	return resEntity, nil
}

func (r *ExamineeQuestionAnswerRepo) GetByIDs(ctx context.Context, examineeId string) (list []*entity.ExamineeAnswerQuestionAnswer, err error) {
	err = r.data.db.WithContext(ctx).Model(&entity.ExamineeAnswerQuestionAnswer{}).Where(" examinee_id = ? ", examineeId).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ExamineeQuestionAnswerRepo) Create(ctx context.Context, examineeAnswer *entity.ExamineeAnswer) error {
	return r.data.db.WithContext(ctx).Create(examineeAnswer).Error
}

func (r *ExamineeQuestionAnswerRepo) SaveAnswer(ctx context.Context, answers []*entity.ExamineeAnswerQuestionAnswer) error {
	mysql.ClauseOnConflict{}
	r.data.db.WithContext(ctx).Clauses()
	return r.data.db.WithContext(ctx).Create(examineeAnswer).Error
}
