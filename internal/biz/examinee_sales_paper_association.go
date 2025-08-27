package biz

import (
	"context"
	"errors"
	v1 "exam_api/api/exam_api/v1"
	"exam_api/internal/data/entity"
	"exam_api/internal/model"
	"exam_api/internal/pkg/icontext"
	innErr "exam_api/internal/pkg/ierrors"
	"github.com/go-kratos/kratos/v2/log"
)

type ExamineeSalesPaperAssociationRepo interface {
	GetPageListByExamineeId(ctx context.Context, in *v1.GetExamPageListRequest, examineeId string) (res []*model.ExamineeSalesPaperAssociation, total int64, err error)
	GetBySalesPaperIds(ctx context.Context, salesPaperIds []string) (list []*entity.ExamineeSalesPaperAssociation, err error)
	GetByExamineeIds(ctx context.Context, examineeIds []string) (list []*entity.ExamineeSalesPaperAssociation, err error)
	GetById(ctx context.Context, id string) (resEntity *entity.ExamineeSalesPaperAssociation, err error)
	UpdateStageNumber(ctx context.Context, examineeSalesPaperAssociationId string, stageNumber v1.StageNumber) (err error)
}

type ExamineeSalesPaperAssociationUseCase struct {
	repo            ExamineeSalesPaperAssociationRepo
	salesPaperCase  *SalesPaperUseCase
	questionUseCase *QuestionUseCase
	log             *log.Helper
}

func NewExamineeSalesPaperAssociationUseCase(repo ExamineeSalesPaperAssociationRepo,
	salesPaperCase *SalesPaperUseCase,
	questionUseCase *QuestionUseCase,
	logger log.Logger) *ExamineeSalesPaperAssociationUseCase {
	return &ExamineeSalesPaperAssociationUseCase{
		repo:            repo,
		salesPaperCase:  salesPaperCase,
		questionUseCase: questionUseCase,
		log:             log.NewHelper(logger),
	}
}

func (uc *ExamineeSalesPaperAssociationUseCase) GetExamPageList(ctx context.Context, req *v1.GetExamPageListRequest) (resp *v1.GetExamPageListResponse, err error) {
	resp = &v1.GetExamPageListResponse{ExamList: make([]*v1.ExamData, 0, 10)}
	if req.PageIndex == 0 {
		req.PageIndex = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	l := uc.log.WithContext(ctx)
	userId, _ := icontext.UserIdFrom(ctx)

	res, total, err := uc.repo.GetPageListByExamineeId(ctx, req, userId)
	if err != nil {
		l.Errorf("GetExamPageList.repo.GetPageListByExamineeId Failed, req:%v, userId:%v, err:%v", req, userId, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	resp.Total = total
	for _, re := range res {
		var status int32 = 1
		if re.StageNumber >= int32(v1.StageNumber_Submit) {
			status = 2
		}
		cur := &v1.ExamData{
			ExamineeAssociationId: re.ID,
			SalesPaperName:        re.SalesPaperName,
			ExamStatus:            status,
		}
		resp.ExamList = append(resp.ExamList, cur)
	}
	return
}

func (uc *ExamineeSalesPaperAssociationUseCase) ExamQuestion(ctx context.Context, req *v1.ExamQuestionRequest) (resp *v1.ExamQuestionResponse, err error) {
	l := uc.log.WithContext(ctx)
	associationId, _ := icontext.AssociationIdFrom(ctx)
	// 根据关联id获取试卷id
	association, err := uc.repo.GetById(ctx, associationId)
	if err != nil {
		l.Errorf("ExamQuestion.repo.GetById Failed, associationId:%v, err:%v", associationId, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	if association == nil {
		err = errors.New("试卷不存在")
		return
	}
	//获取试卷题目
	resp, err = uc.questionUseCase.ExamQuestion(ctx, association.SalesPaperID)
	if err != nil {
		l.Errorf("ExamQuestion.questionUseCase.ExamQuestion Failed, req:%v, salesPaperId:%v, err:%v", req, association.SalesPaperID, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	return
}

func (uc *ExamineeSalesPaperAssociationUseCase) GetById(ctx context.Context, id string) (resEntity *entity.ExamineeSalesPaperAssociation, err error) {
	return uc.repo.GetById(ctx, id)
}

func (uc *ExamineeSalesPaperAssociationUseCase) UpdateStageNumber(ctx context.Context, examineeSalesPaperAssociationId string, stageNumber v1.StageNumber) (err error) {
	l := uc.log.WithContext(ctx)
	err = uc.repo.UpdateStageNumber(ctx, examineeSalesPaperAssociationId, stageNumber)
	if err != nil {
		l.Errorf("UpdateStageNumber.repo.UpdateStageNumber Failed, examineeSalesPaperAssociationId:%v, stageNumber:%v, err:%v", examineeSalesPaperAssociationId, stageNumber, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	return
}
