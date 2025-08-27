package data

import (
	"context"
	"exam_api/internal/biz"
	"exam_api/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type QuestionRepo struct {
	data *Data
	log  *log.Helper
}

func NewQuestionRepo(data *Data, logger log.Logger) biz.QuestionRepo {
	return &QuestionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *QuestionRepo) GetList(ctx context.Context, salesPaperId string) (res []*entity.Question, err error) {
	err = r.data.db.WithContext(ctx).Model(&entity.Question{}).
		Where(" sales_paper_id = ?", salesPaperId).
		Order(" `order` asc").
		Find(&res).Error
	if err != nil {
		return
	}

	return
}

func (r *QuestionRepo) GetPageListBySalesPaperId(ctx context.Context, salesPaperId string) (res []*entity.Question, err error) {

	session := r.data.db.WithContext(ctx)
	session = session.Table((&entity.Question{}).TableName())
	q, v := r.buildConditions(salesPaperId)
	if q != "" {
		session.Where(q, v...)
	}
	err = session.
		Order(" `order` asc").
		Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *QuestionRepo) GetListBySalesPaperId(ctx context.Context, salesPaperId string) (res []*entity.Question, err error) {
	err = r.data.db.WithContext(ctx).Model(&entity.Question{}).
		Where(" sales_paper_id = ? ", salesPaperId).
		Order(" `order` asc").
		Find(&res).Error
	if err != nil {
		return
	}

	return
}

func (r *QuestionRepo) GetOptionList(ctx context.Context, questionId string) (res []*entity.QuestionOption, err error) {
	err = r.data.db.WithContext(ctx).Model(&entity.QuestionOption{}).
		Where(" question_id = ?", questionId).
		Order(" `order` asc").
		Find(&res).Error
	if err != nil {
		return
	}

	return
}

func (r *QuestionRepo) GetOptionListByQuestionIds(ctx context.Context, questionIds []string) (res map[string][]*entity.QuestionOption, err error) {
	qOptions := make([]*entity.QuestionOption, 0, 10)
	res = make(map[string][]*entity.QuestionOption)
	err = r.data.db.WithContext(ctx).Model(&entity.QuestionOption{}).
		Where(" question_id in ?", questionIds).
		Order(" `order` asc").
		Find(&qOptions).Error
	if err != nil {
		return
	}
	for _, option := range qOptions {
		value, ok := res[option.QuestionID]
		if !ok {
			res[option.QuestionID] = []*entity.QuestionOption{option}
		} else {
			res[option.QuestionID] = append(value, option)
		}
	}
	return
}
func (r *QuestionRepo) GetById(ctx context.Context, questionId string) (qEntity *entity.Question, qOptionsEntities []*entity.QuestionOption, err error) {
	d := r.data.db.WithContext(ctx)
	qEntity, err = getSingleRecordByScope[entity.Question](
		d.Model(qEntity).Where(" id = ? ", questionId),
	)
	err = d.Model(qOptionsEntities).Where(" question_id = ? ", questionId).Find(&qOptionsEntities).Error
	if err != nil {
		return
	}
	return
}

func (r *QuestionRepo) buildConditions(salesPaperId string) (string, []interface{}) {
	var (
		query strings.Builder
		value []interface{}
		q     string
	)
	if salesPaperId != "" {
		query.WriteString(" sales_paper_id = ? ")
		value = append(value, salesPaperId)
		query.WriteString(" AND")
	}
	// 过滤已删除
	query.WriteString(" deleted_at is NULL ")
	query.WriteString(" AND")

	if query.String() != "" {
		q = strings.TrimSuffix(query.String(), "AND")
	} else {
		q = query.String()
	}

	return q, value
}
