package biz

import (
	"context"
	"exam_api/internal/data/entity"
	"github.com/go-kratos/kratos/v2/log"
)

type ExamineeQuestionAnswerRepo interface {
	GetByExamineeAnswerId(ctx context.Context, examineeAnswerId string) (list []*entity.ExamineeAnswerQuestionAnswer, err error)
	SaveAnswer(ctx context.Context, answers []*entity.ExamineeAnswerQuestionAnswer) error
}

type ExamineeQuestionAnswerUseCase struct {
	repo ExamineeQuestionAnswerRepo
	log  *log.Helper
}

func NewExamineeQuestionAnswerUseCase(repo ExamineeQuestionAnswerRepo,
	logger log.Logger) *ExamineeQuestionAnswerUseCase {
	return &ExamineeQuestionAnswerUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
func (uc *ExamineeQuestionAnswerUseCase) GetByExamineeAnswerId(ctx context.Context, examineeAnswerId string) (list []*entity.ExamineeAnswerQuestionAnswer, err error) {
	return uc.repo.GetByExamineeAnswerId(ctx, examineeAnswerId)
}

func (uc *ExamineeQuestionAnswerUseCase) SaveAnswer(ctx context.Context, answers []*entity.ExamineeAnswerQuestionAnswer) error {
	return uc.repo.SaveAnswer(ctx, answers)
}
