package errors

import (
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrBadRequest          = errors.New(400, "INVALID_ARGS", "INVALID ARGS")
	ErrLogin               = errors.New(401, "UNAUTHORIZED", "请授权登录")
	ErrRefreshTokenExpired = errors.New(400, "REFRESH_TOKEN_EXPIRED", "请重新登录")
	ErrAccessTokenExpired  = errors.New(400, "ACCESS_TOKEN_EXPIRED", "token已经过期")
	ErrResourceNotFound    = errors.New(404, "RESOURCE_NOT_FOUND", "资源不存在")
	ErrInternalServer      = errors.New(500, "INTERNAL_SERVER_ERROR", "服务内部错误")
)

func WithReason(e *errors.Error, in string) *errors.Error {
	return errors.New(int(e.Code), in, e.Message)
}

func WithMessage(e *errors.Error, msg string) *errors.Error {
	return errors.New(int(e.Code), e.Reason, msg)
}
