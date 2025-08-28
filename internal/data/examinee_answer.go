package data

import (
	"context"
	"exam_api/internal/biz"
	"exam_api/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type ExamineeAnswerRepo struct {
	data *Data
	log  *log.Helper
}

func NewExamineeAnswerRepo(data *Data, logger log.Logger) biz.ExamineeAnswerRepo {
	return &ExamineeAnswerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *ExamineeAnswerRepo) GetByAssociationId(ctx context.Context, associationId string) (resEntity *entity.ExamineeAnswer, err error) {
	resEntity, err = getSingleRecordByScope[entity.ExamineeAnswer](
		r.data.db.WithContext(ctx).Model(resEntity).Where(" examinee_sales_paper_association_id = ? ", associationId),
	)
	if err != nil {
		return nil, err
	}
	return resEntity, nil
}

func (r *ExamineeAnswerRepo) GetByIDs(ctx context.Context, examineeId string) (list []*entity.ExamineeAnswer, err error) {
	err = r.data.db.WithContext(ctx).Model(&entity.ExamineeAnswer{}).Where(" examinee_id = ? ", examineeId).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ExamineeAnswerRepo) Create(ctx context.Context, examineeAnswer *entity.ExamineeAnswer) error {
	return r.data.db.WithContext(ctx).Create(examineeAnswer).Error
}

// 更新最新动作
func (r *ExamineeAnswerRepo) UpdateAction(ctx context.Context, examineeAnswerId string, lastActionTime, lastActionTime2 time.Time, remaining int32, completeQuestionNum int32) (int64, error) {
	// 准备更新字段
	updates := map[string]interface{}{
		"last_action_time":      lastActionTime,
		"remaining_timelimit":   remaining,
		"complete_question_num": completeQuestionNum,
		"updated_by":            "service",
	}
	// 执行更新
	res := r.data.db.WithContext(ctx).Model(&entity.ExamineeAnswer{}).
		Where(" id = ? and last_action_time = ? ", examineeAnswerId, lastActionTime2.Format(time.DateTime)).
		Updates(updates)

	return res.RowsAffected, res.Error
}

// 更新结果
func (r *ExamineeAnswerRepo) UpdateResult(ctx context.Context, examineeAnswerId string, score float64, comparability, usability int32) error {
	// 准备更新字段
	updates := map[string]interface{}{
		"score":         score,
		"comparability": comparability,
		"usability":     usability,
		"updated_by":    "service",
	}
	// 执行更新
	err := r.data.db.WithContext(ctx).Model(&entity.ExamineeAnswer{}).
		Where(" id = ? ", examineeAnswerId).
		Updates(updates).Error

	return err
}

// 提交试卷
func (r *ExamineeAnswerRepo) SubmitResult(ctx context.Context, examineeAnswerId string, score float64, comparability, usability int32) error {
	// 准备更新字段
	updates := map[string]interface{}{
		"submit_time":           score,
		"complete_question_num": comparability,
		"last_action_time":      usability,
		"updated_by":            "service",
	}
	// 执行更新
	err := r.data.db.WithContext(ctx).Model(&entity.ExamineeAnswer{}).
		Where(" id = ? ", examineeAnswerId).
		Updates(updates).Error

	return err
}
