package biz

import (
	"context"
	v1 "exam_api/api/exam_api/v1"
	"exam_api/internal/data/entity"
	innErr "exam_api/internal/pkg/ierrors"
	"exam_api/internal/pkg/iutils"
	"github.com/go-kratos/kratos/v2/log"
)

type QuestionRepo interface {
	GetList(ctx context.Context, salesPaperId string) (res []*entity.Question, err error)
	GetPageListBySalesPaperId(ctx context.Context, in *v1.ExamQuestionRequest, salesPaperId string) (res []*entity.Question, total int64, err error)
	GetListBySalesPaperId(ctx context.Context, salesPaperId string) (res []*entity.Question, err error)
	GetOptionList(ctx context.Context, questionId string) (res []*entity.QuestionOption, err error)
	GetOptionListByQuestionIds(ctx context.Context, questionIds []string) (res map[string][]*entity.QuestionOption, err error)
	GetById(ctx context.Context, questionId string) (qEntity *entity.Question, qOptionsEntities []*entity.QuestionOption, err error)
}

type QuestionUseCase struct {
	repo QuestionRepo
	log  *log.Helper
}

func NewQuestionUseCase(repo QuestionRepo, logger log.Logger) *QuestionUseCase {
	return &QuestionUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *QuestionUseCase) ExamQuestion(ctx context.Context, req *v1.ExamQuestionRequest, salesPaperId string) (resp *v1.ExamQuestionResponse, err error) {
	if req.PageIndex == 0 {
		req.PageIndex = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 8
	}
	resp = &v1.ExamQuestionResponse{QuestionData: make([]*v1.QuestionData, 0, req.PageSize)}
	l := uc.log.WithContext(ctx)

	res, total, err := uc.repo.GetPageListBySalesPaperId(ctx, req, salesPaperId)
	if err != nil {
		l.Errorf("GetQuestionBySalesPaperId.repo.GetPageListBySalesPaperId Failed, req:%v, err:%v", req, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	resp.Total = total
	questionIds := make([]string, 0, len(res))
	for _, re := range res {
		questionIds = append(questionIds, re.ID)
	}
	mQuestionOptions, err := uc.repo.GetOptionListByQuestionIds(ctx, questionIds)
	if err != nil {
		l.Errorf("GetQuestionBySalesPaperId.repo.GetOptionListByQuestionIds Failed, questionIds:%v, err:%v", questionIds, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	for _, re := range res {
		cur := &v1.QuestionData{
			QuestionId:     re.ID,
			Title:          re.Title,
			QuestionTypeId: v1.QuestionType(re.QuestionTypeID),
			Order:          re.Order_,
		}
		if v, ok := mQuestionOptions[cur.QuestionId]; ok {
			for _, option := range v {
				cur.QuestionOptionsData = append(cur.QuestionOptionsData, &v1.QuestionOptionData{
					QuestionOptionId: option.ID,
					Description:      option.Description,
					SerialNumber:     iutils.OrderToLetter(option.Order_),
				})
			}
		}
		resp.QuestionData = append(resp.QuestionData, cur)
	}
	return
}
