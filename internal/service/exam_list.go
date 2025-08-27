package service

import (
	"context"
	v1 "exam_api/api/exam_api/v1"
)

func (s *ExamService) GetExamPageList(ctx context.Context, in *v1.GetExamPageListRequest) (*v1.GetExamPageListResponse, error) {
	return s.examineeSalesPaperAssociationUc.GetExamPageList(ctx, in)
}

func (s *ExamService) ExamQuestion(ctx context.Context, in *v1.ExamQuestionRequest) (*v1.ExamQuestionResponse, error) {
	return s.examineeSalesPaperAssociationUc.ExamQuestion(ctx, in)
}

func (s *ExamService) StartExam(ctx context.Context, in *v1.StartExamRequest) (*v1.StartExamResponse, error) {
	return s.examineeAnswerUseCase.StartExam(ctx, in)
}
