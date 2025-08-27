package biz

import (
	"context"
	"errors"
	v1 "exam_api/api/exam_api/v1"
	_const "exam_api/internal/const"
	"exam_api/internal/data/entity"
	"exam_api/internal/middleware"
	"exam_api/internal/pkg/isecurity"
	"exam_api/internal/pkg/isnowflake"
	"github.com/go-kratos/kratos/v2/log"
)

type SysLoginRepo interface {
	Create(ctx context.Context, entity *entity.SysLoginRecord) error
}

type LoginUseCase struct {
	repo     ExamineeRepo
	sysLogin SysLoginRepo
	log      *log.Helper
}

func NewLoginUseCase(repo ExamineeRepo, sysLogin SysLoginRepo, logger log.Logger) *LoginUseCase {
	return &LoginUseCase{repo: repo, sysLogin: sysLogin, log: log.NewHelper(logger)}
}

func (uc *LoginUseCase) ExamLogin(ctx context.Context, req *v1.ExamLoginRequest) (resp *v1.ExamLoginResponse, err error) {
	resp = &v1.ExamLoginResponse{}
	l := uc.log.WithContext(ctx)
	user, err := uc.repo.GetByEmail(ctx, req.LoginAccount)
	if err != nil {
		l.Errorf("Login.repo.GetByEmail Failed, req:%v, err:%v", req, err.Error())
		return nil, err
	}
	if user == nil {
		return resp, errors.New("用户不存在")
	}
	// 验证密码
	if req.PassWord != "" {
		ok := isecurity.CheckPassword(req.PassWord, user.HashPassword)
		if !ok {
			err = errors.New("密码错误")
			return
		}
	}
	if user.Status != int32(v1.ExamineeStatus_ExamineeActive) {
		err = errors.New("用户未激活")
		return
	}
	id, err := isnowflake.SnowFlake.NextID(_const.SysLoginRecordPrefix)
	if err != nil {
		l.Errorf("Login.isnowflake.SnowFlake.NextID Failed, req:%v, err:%v", req, err.Error())
		return
	}
	uc.sysLogin.Create(ctx, &entity.SysLoginRecord{
		ID:            id,
		UserID:        user.ID,
		LoginPlatform: int32(v1.LoginPlatform_Exam),
	})
	// 生成jwt
	accessJWT, err := middleware.JWT.GenerateAccessToken(user.ID, user.UserName, "")
	if err != nil {
		// 处理错误
		err = errors.New("登录失败")
		return
	}
	resp.UserName = user.UserName
	resp.Token = accessJWT
	return resp, nil
}
