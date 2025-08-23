package middleware

import (
	"reflect"
	"time"

	"github.com/airunny/copier"
	"github.com/airunny/wiki-go-tools/env"
	"github.com/go-kratos/kratos/v2/errors"
)

type BizResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message,omitempty"`
	Reason  string      `json:"reason,omitempty"`
	Time    int64       `json:"time,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var (
	ErrInvalidArgs = &BizResponse{
		Code:    400,
		Message: "INVALID ARGS",
		Reason:  "INVALID_ARGS",
	}
	ErrInternalServer = &BizResponse{
		Code:    500,
		Message: "INTERNAL SERVER ERROR ",
		Reason:  "INTERNAL_SERVER_ERROR",
	}
)

func (b *BizResponse) WithMessage(msg string) *BizResponse {
	return &BizResponse{
		Code:    b.Code,
		Message: msg,
		Reason:  b.Reason,
		Time:    time.Now().Unix(),
	}
}

func ResponseWithData(data interface{}) *BizResponse {
	typ := reflect.TypeOf(data)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	out := reflect.New(typ)
	err := copier.CopyWithOptional(out.Interface(), data,
		copier.WithDeepCopyOption(),
		copier.WithInitNilSlice())
	if err != nil {
		return ResponseWithError(errors.FromError(err))
	}

	return &BizResponse{
		Code: 200,
		Time: time.Now().Unix(),
		Data: out.Interface(),
	}
}

func ResponseWithError(err *errors.Error) *BizResponse {
	if err == nil {
		return ResponseWithData(nil)
	}

	if err.Code == 400 && err.Reason == "CODEC" {
		return &BizResponse{
			Code:    ErrInvalidArgs.Code,
			Message: err.Error(),
			Reason:  ErrInvalidArgs.Reason,
			Time:    time.Now().Unix(),
		}
	}

	if env.Environment() == env.ProdMode && err.Code == 500 {
		return &BizResponse{
			Code:    ErrInternalServer.Code,
			Message: ErrInternalServer.Message,
			Reason:  ErrInternalServer.Reason,
			Time:    time.Now().Unix(),
		}
	}

	return &BizResponse{
		Code:    err.Code,
		Message: err.Message,
		Reason:  err.Reason,
		Time:    time.Now().Unix(),
	}
}
