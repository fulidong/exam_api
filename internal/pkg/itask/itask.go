package itask

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

func TaskWithContext(ctx context.Context, fn func(), log *log.Helper) {
	select {
	case <-ctx.Done():
		// 如果超时，记录错误，返回nil
		log.Error("[taskWithContext] err : " + ctx.Err().Error())
	default:
		fn()
	}
}
