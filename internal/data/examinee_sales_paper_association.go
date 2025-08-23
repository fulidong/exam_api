package data

import (
	"context"
	v1 "exam_api/api/exam_api/v1"
	"exam_api/internal/biz"
	"exam_api/internal/data/entity"
	"exam_api/internal/model"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type ExamineeSalesPaperAssociationRepo struct {
	data *Data
	log  *log.Helper
}

func NewExamineeSalesPaperAssociationRepo(data *Data, logger log.Logger) biz.ExamineeSalesPaperAssociationRepo {
	return &ExamineeSalesPaperAssociationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 获取发放试卷列表
func (r *ExamineeSalesPaperAssociationRepo) GetPageListByExamineeId(ctx context.Context, in *v1.GetExamPageListRequest, examineeId string) (res []*model.ExamineeSalesPaperAssociation, total int64, err error) {
	session := r.data.db.WithContext(ctx)
	session = session.Table((&entity.ExamineeSalesPaperAssociation{}).TableName())
	q, v := r.buildConditions(examineeId)
	if q != "" {
		session.Where(q, v...)
	}
	session.Count(&total)
	err = session.
		Order("created_at desc").
		Offset(int((in.PageIndex - 1) * in.PageSize)).
		Limit(int(in.PageSize)).
		Find(&res).Error
	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

func (r *ExamineeSalesPaperAssociationRepo) GetBySalesPaperIds(ctx context.Context, salesPaperIds []string) (list []*entity.ExamineeSalesPaperAssociation, err error) {
	err = r.data.db.WithContext(ctx).Model(&entity.ExamineeSalesPaperAssociation{}).Where(" sales_paper_id in ? ", salesPaperIds).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ExamineeSalesPaperAssociationRepo) GetByExamineeIds(ctx context.Context, examineeIds []string) (list []*entity.ExamineeSalesPaperAssociation, err error) {
	err = r.data.db.WithContext(ctx).Model(&entity.ExamineeSalesPaperAssociation{}).Where(" examinee_id in ? ", examineeIds).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ExamineeSalesPaperAssociationRepo) GetById(ctx context.Context, id string) (resEntity *entity.ExamineeSalesPaperAssociation, err error) {
	resEntity, err = getSingleRecordByScope[entity.ExamineeSalesPaperAssociation](
		r.data.db.WithContext(ctx).Model(resEntity).Where(" id = ? ", id),
	)
	if err != nil {
		return nil, err
	}
	return resEntity, nil
}

// 更新进度
func (r *ExamineeSalesPaperAssociationRepo) UpdateStageNumber(ctx context.Context, examineeSalesPaperAssociationId string, stageNumber v1.StageNumber) (err error) {
	return r.data.db.WithContext(ctx).Model(&entity.ExamineeSalesPaperAssociation{}).
		Where(" id = ? ", examineeSalesPaperAssociationId).
		Updates(map[string]interface{}{
			"stage_number": stageNumber,
		}).Error
}

func (r *ExamineeSalesPaperAssociationRepo) buildConditions(examineeId string) (string, []interface{}) {
	var (
		query strings.Builder
		value []interface{}
		q     string
	)
	if examineeId != "" {
		query.WriteString(" examinee_id = ? ")
		value = append(value, examineeId)
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
