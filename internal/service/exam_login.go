package service

import (
	"context"
	v1 "exam_api/api/exam_api/v1"
)

func (s *ExamService) ExamLogin(ctx context.Context, in *v1.ExamLoginRequest) (*v1.ExamLoginResponse, error) {
	return s.loginUc.ExamLogin(ctx, in)
}
