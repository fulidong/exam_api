package biz

import (
	"context"
	"errors"
	"exam_api/internal/data/entity"
	innErr "exam_api/internal/pkg/ierrors"
	"github.com/go-kratos/kratos/v2/log"
)

type SalesPaperRepo interface {
	GetByID(ctx context.Context, salesPaperId string) (resEntity *entity.SalesPaper, err error)
}

type SalesPaperUseCase struct {
	repo SalesPaperRepo
	log  *log.Helper
}

func NewSalesPaperUseCase(repo SalesPaperRepo, logger log.Logger) *SalesPaperUseCase {
	return &SalesPaperUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *SalesPaperUseCase) GetSalesPaperDetail(ctx context.Context, salesPaperId string) (resp *entity.SalesPaper, err error) {
	l := uc.log.WithContext(ctx)
	res, err := uc.repo.GetByID(ctx, salesPaperId)
	if err != nil {
		l.Errorf("GetSalesPaperDetail.repo.GetByID Failed, salesPaperId:%v, err:%v", salesPaperId, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	if res == nil {
		err = errors.New("试卷不存在")
		return
	}
	return
}

func (uc *SalesPaperUseCase) CheckSalesPaper(ctx context.Context, iSalesPaperId string, l *log.Helper) (err error) {
	salesPaper, err := uc.repo.GetByID(ctx, iSalesPaperId)
	if err != nil {
		l.Errorf("CheckSalesPaper.repo.GetByID Failed, req:%v, err:%v", err, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	if salesPaper == nil {
		err = errors.New("试卷不存在")
		return
	}
	return
}
