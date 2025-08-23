package ilog

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

var (
	logBufPool *sync.Pool
)

func LoggingGRPC() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				defer func(start time.Time) {
					reqStr, _ := json.Marshal(req)
					replyStr, _ := json.Marshal(reply)
					log.Context(ctx).Infow(
						"operation", tr.Operation(),
						"duration", time.Now().Sub(start),
						"request", string(reqStr),
						"reply", string(replyStr))
				}(time.Now())
			}
			return handler(ctx, req)
		}
	}
}

type Body struct {
	buf *bytes.Buffer
	req *http.Request
}

func (b Body) Read(p []byte) (n int, err error) {
	return b.buf.Read(p)
}

func (b Body) Close() error {
	return b.req.Body.Close()
}
