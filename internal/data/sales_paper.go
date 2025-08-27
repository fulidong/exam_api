package data

import (
	"context"
	"exam_api/internal/biz"
	"exam_api/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
)

type SalesPaperRepo struct {
	data *Data
	log  *log.Helper
}

func NewSalesPaperRepo(data *Data, logger log.Logger) biz.SalesPaperRepo {
	return &SalesPaperRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *SalesPaperRepo) GetByID(ctx context.Context, salesPaperId string) (resEntity *entity.SalesPaper, err error) {
	resEntity, err = getSingleRecordByScope[entity.SalesPaper](
		r.data.db.WithContext(ctx).Model(resEntity).Where(" id = ? and is_enabled = 1", salesPaperId),
	)
	if err != nil {
		return nil, err
	}
	return resEntity, nil
}
