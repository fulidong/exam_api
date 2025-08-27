package middleware

import (
	"context"
	_const "exam_api/internal/const"
	"exam_api/internal/pkg/iclient"
	"exam_api/internal/pkg/icontext"
	"exam_api/internal/pkg/iheader"
	"exam_api/internal/pkg/ijwt"
	"github.com/airunny/wiki-go-tools/reqid"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	stdHttp "net/http"
	"strings"
	"time"
)

var JWT *ijwt.SecureJWT

var (
	allowOrigins = []string{"*"}
	allowHeaders = []string{"X-Token", "Authorization", "Content-Type", "X-User-Id"}
	allowMethods = []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}
)

func CORS() http.FilterFunc {
	return handlers.CORS(
		handlers.AllowedOrigins(allowOrigins),
		handlers.AllowedHeaders(allowHeaders),
		handlers.AllowedMethods(allowMethods),
		handlers.OptionStatusCode(204),
	)
}

func RequestIdWithHeader(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		tr, ok := transport.FromServerContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		var (
			requestId string
		)

		requestId = iheader.GetRequestId(tr.RequestHeader())
		if requestId == "" {
			requestId = reqid.GenRequestID()
		}

		tr.ReplyHeader().Set(iheader.RequestIdKey, requestId)
		ctx = icontext.WithRequestId(ctx, requestId)

		return handler(ctx, req)
	}
}

func TryParseHeader(opts ...Option) middleware.Middleware {
	o := Options{}
	for _, opt := range opts {
		opt(o)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			header := tr.RequestHeader()
			// 获取当前请求的 FullMethod
			method := tr.Operation()

			// 判断是否是需要跳过的接口
			if method == "/exam_api.v1.ExamService/ExamLogin" {
				return handler(ctx, req)
			}

			jwt := iheader.GetToken(header)
			if jwt == "" {
				// 返回 Unauthorized 错误
				return nil, errors.Unauthorized("Unauthorized", "missing or invalid token")
			}
			claims, err := JWT.VerifyAccessToken(jwt)
			if err != nil {
				// JWT 解析失败
				return nil, errors.Unauthorized("Unauthorized", "invalid token: "+err.Error())
			}
			if claims.Exp < time.Now().Unix() {
				// JWT 国旗
				return nil, errors.Unauthorized("Unauthorized", " token expiration ")
			}
			// 用户ID
			ctx = icontext.WithUserIdKey(ctx, claims.UserID)
			// 用户名
			ctx = icontext.WithUserNameKey(ctx, claims.Username)
			// 用户角色
			ctx = icontext.WithUserRuleKey(ctx, claims.Role)
			// accessToken
			ctx = icontext.WithUserTokenKey(ctx, jwt)
			// user agent
			ctx = icontext.WithUserAgentKey(ctx, iheader.GetUserAgent(header))
			return handler(ctx, req)
		}
	}
}

func AuthExamTokenMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}
			header := tr.RequestHeader()
			clientInfo := &iclient.ClientInfo{}
			if ok {
				clientInfo = &iclient.ClientInfo{
					IP:         getClientIP(tr),
					UserAgent:  header.Get("User-Agent"),
					Accept:     header.Get("Accept"),
					AcceptLang: header.Get("Accept-Language"),
					Platform:   header.Get("Sec-CH-UA-Platform"),
					CanvasFP:   header.Get("X-Canvas-FP"),
					WebGLFP:    header.Get("X-WebGL-FP"),
					FontsFP:    header.Get("X-Fonts-FP"),
				}
				// 将 clientInfo 放入 context
				ctx = icontext.WithUserClientKey(ctx, clientInfo)
			}
			// 获取当前请求的 FullMethod
			method := tr.Operation()
			jwt, _ := icontext.UserTokenFrom(ctx)
			//判断是否需要验证考试token
			if _, ok := _const.VerifyExamTokenMethod[method]; !ok {
				return handler(ctx, req)
			}
			examJwt := iheader.GetExamToken(header)
			if examJwt == "" {
				// 返回 Unauthorized 错误
				return nil, errors.Unauthorized("Unauthorized", "missing or invalid exam token")
			}
			claims, err := JWT.VerifyExamToken(jwt, examJwt, clientInfo)
			if err != nil || claims == nil || claims.AssociationId == "" || claims.UserID == "" {
				// JWT 解析失败
				return nil, errors.Unauthorized("Unauthorized", "invalid token: "+err.Error())
			}
			if claims.ExpiresAt < time.Now().Unix() {
				return nil, errors.Unauthorized("Unauthorized", " exam token expiration ")
			}
			ctx = icontext.WithAssociationIdKey(ctx, claims.AssociationId)
			ctx = icontext.WithSessionIdKey(ctx, claims.SessionID)
			ctx = icontext.WithExamTokenKey(ctx, examJwt)
			return handler(ctx, req)
		}
	}
}

func ResponseEncoder(w http.ResponseWriter, r *stdHttp.Request, v interface{}) error {
	if v == nil {
		return nil
	}

	if rd, ok := v.(http.Redirector); ok {
		url, code := rd.Redirect()
		stdHttp.Redirect(w, r, url, code)
		return nil
	}

	if res, ok := v.(TextPlainReply); ok {
		w.Header().Set("Content-Type", iheader.ResponseContentTextType)
		_, err := w.Write([]byte(res.StringReply()))
		if err != nil {
			w.WriteHeader(stdHttp.StatusInternalServerError)
		}
		return nil
	}

	WriteResponse(w, r, ResponseWithData(v))
	return nil
}

func ErrorEncoder(w http.ResponseWriter, r *stdHttp.Request, err error) {
	WriteResponse(w, r, ResponseWithError(errors.FromError(err)))
}

func WriteResponse(w http.ResponseWriter, _ *stdHttp.Request, body interface{}) {
	codec := encoding.GetCodec(json.Name)
	data, err := codec.Marshal(body)
	if err != nil {
		w.WriteHeader(stdHttp.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", iheader.ResponseContentJsonType)
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(stdHttp.StatusInternalServerError)
	}
}

// 获取客户端IP
func getClientIP(tr transport.Transporter) string {
	// 检查代理头
	host := ""
	if ip := tr.RequestHeader().Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := tr.RequestHeader().Get("X-Real-Ip"); ip != "" {
		return ip
	}

	// 使用远程地址
	if ht, ok := tr.(*http.Transport); ok {
		// 获取 *http.Request
		httpRequest := ht.Request()
		if httpRequest != nil {
			host, _, _ = strings.Cut(httpRequest.RemoteAddr, ":")
			// 你可以选择将 remoteAddr 存入 ctx 供后续 handler 使用
			// ctx = context.WithValue(ctx, "remote_addr", remoteAddr)
		}
	}

	return host
}
