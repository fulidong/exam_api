package data

import (
	"context"
	"exam_api/internal/biz"
	"exam_api/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
)

type ExamineeRepo struct {
	data *Data
	log  *log.Helper
}

func NewExamineeRepo(data *Data, logger log.Logger) biz.ExamineeRepo {
	return &ExamineeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *ExamineeRepo) GetByEmail(ctx context.Context, email string) (resEntity *entity.Examinee, err error) {
	resEntity, err = getSingleRecordByScope[entity.Examinee](
		r.data.db.WithContext(ctx).Model(resEntity).Where(" email = ? ", email),
	)
	if err != nil {
		return nil, err
	}
	return resEntity, nil
}

func (r *ExamineeRepo) GetByID(ctx context.Context, examineeId string) (resEntity *entity.Examinee, err error) {
	resEntity, err = getSingleRecordByScope[entity.Examinee](
		r.data.db.WithContext(ctx).Model(resEntity).Where(" id = ? ", examineeId),
	)
	if err != nil {
		return nil, err
	}
	return resEntity, nil
}

func (r *ExamineeRepo) GetByIDs(ctx context.Context, examineeIds []string) (list []*entity.Examinee, err error) {
	err = r.data.db.WithContext(ctx).Model(&entity.Examinee{}).Where(" id in ? ", examineeIds).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
