package biz

import (
	"context"
	"errors"
	"exam_api/internal/data/entity"
	innErr "exam_api/internal/pkg/ierrors"
	"github.com/go-kratos/kratos/v2/log"
)

type ExamineeRepo interface {
	GetByEmail(ctx context.Context, email string) (resEntity *entity.Examinee, err error)
	GetByID(ctx context.Context, examineeId string) (resEntity *entity.Examinee, err error)
	GetByIDs(ctx context.Context, examineeIds []string) (list []*entity.Examinee, err error)
}

type ExamineeUseCase struct {
	repo ExamineeRepo
	log  *log.Helper
}

func NewExamineeUseCase(repo ExamineeRepo, logger log.Logger) *ExamineeUseCase {
	return &ExamineeUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *ExamineeUseCase) GetExamineeDetail(ctx context.Context, examineeId string) (resp *entity.Examinee, err error) {
	l := uc.log.WithContext(ctx)
	res, err := uc.repo.GetByID(ctx, examineeId)
	if err != nil {
		l.Errorf("GetExamineeDetail.repo.GetByID Failed, examineeId:%v, err:%v", examineeId, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	if res == nil {
		err = errors.New("考生不存在")
		return
	}

	return
}
