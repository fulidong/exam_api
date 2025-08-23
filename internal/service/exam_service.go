package service

import (
	v1 "exam_api/api/exam_api/v1"
	"exam_api/internal/biz"
)

type ExamService struct {
	v1.UnimplementedExamServiceServer
	loginUc                         *biz.LoginUseCase
	examineeSalesPaperAssociationUc *biz.ExamineeSalesPaperAssociationUseCase
	questionUc                      *biz.QuestionUseCase
	salesPaperUseCase               *biz.SalesPaperUseCase
	examineeAnswerUseCase           *biz.ExamineeAnswerUseCase
}

func NewExamService(loginUc *biz.LoginUseCase,
	examineeSalesPaperAssociationUc *biz.ExamineeSalesPaperAssociationUseCase,
	questionUc *biz.QuestionUseCase,
	salesPaperUseCase *biz.SalesPaperUseCase,
	examineeAnswerUseCase *biz.ExamineeAnswerUseCase) *ExamService {
	return &ExamService{
		loginUc:                         loginUc,
		examineeSalesPaperAssociationUc: examineeSalesPaperAssociationUc,
		questionUc:                      questionUc,
		salesPaperUseCase:               salesPaperUseCase,
		examineeAnswerUseCase:           examineeAnswerUseCase,
	}
}
