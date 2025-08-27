package icontext

import (
	"context"
	"encoding/json"
	"exam_api/internal/pkg/iclient"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"os"
)

const (
	requestIdKey     = "RequestId"        // req id
	userIdKey        = "userIdKey"        // 用户ID
	userNameKey      = "userNameKey"      // 用户名称
	userRuleKey      = "userRuleKey"      // 用户角色
	userTokenKey     = "userTokenKey"     // 登录token
	userClientKey    = "userClientKey"    //用户客户端信息
	associationIdKey = "associationIdKey" //关联id信息
	sessionIdKey     = "sessionIdKey"     //关联id信息
	examTokenKey     = "examTokenKey"     // exam token
	userAgentKey     = "userAgentKey"     // exam token
)

func withValue(ctx context.Context, key, value string) context.Context {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		md = metadata.Metadata{}
	}
	md.Set(key, value)
	return metadata.NewServerContext(ctx, md)
	//return metadata.AppendToClientContext(ctx, key, value)
}

func fromValue(ctx context.Context, key string) (string, bool) {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return "", false
	}

	out := md.Get(key)
	return out, out != ""
}

func WithUserIdKey(ctx context.Context, in string) context.Context {
	return withValue(ctx, userIdKey, in)
}

func UserIdFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, userIdKey)
}

func WithUserNameKey(ctx context.Context, in string) context.Context {
	return withValue(ctx, userNameKey, in)
}

func UserNameFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, userNameKey)
}

func WithUserRuleKey(ctx context.Context, in string) context.Context {
	return withValue(ctx, userRuleKey, in)
}

func UserTokenFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, userTokenKey)
}

func WithUserTokenKey(ctx context.Context, in string) context.Context {
	return withValue(ctx, userTokenKey, in)
}

func AssociationIdFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, associationIdKey)
}

func WithAssociationIdKey(ctx context.Context, in string) context.Context {
	return withValue(ctx, associationIdKey, in)
}

func SessionIdFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, sessionIdKey)
}

func WithSessionIdKey(ctx context.Context, in string) context.Context {
	return withValue(ctx, sessionIdKey, in)
}

func ExamTokenFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, examTokenKey)
}

func WithExamTokenKey(ctx context.Context, in string) context.Context {
	return withValue(ctx, examTokenKey, in)
}

func UserAgentFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, userAgentKey)
}

func WithUserAgentKey(ctx context.Context, in string) context.Context {
	return withValue(ctx, userAgentKey, in)
}

func UserClientFrom(ctx context.Context) (*iclient.ClientInfo, bool) {
	v, ok := fromValue(ctx, userClientKey)
	if !ok {
		return &iclient.ClientInfo{}, false
	}
	body := []byte(v)
	c := &iclient.ClientInfo{}
	err := json.Unmarshal(body, &c)
	return c, err == nil
}

func WithUserClientKey(ctx context.Context, in *iclient.ClientInfo) context.Context {
	if in == nil {
		in = &iclient.ClientInfo{}
	}
	body, _ := json.Marshal(in)
	return withValue(ctx, userClientKey, string(body))
}

func UserRuleFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, userRuleKey)
}

// request id

func WithRequestId(ctx context.Context, in string) context.Context {
	return withValue(ctx, requestIdKey, in)
}

func RequestIdFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, requestIdKey)
}

// context

func LoggerValues() []interface{} {
	return []interface{}{

		"client_info", log.Valuer(func(ctx context.Context) interface{} {
			clientIp, _ := UserClientFrom(ctx)
			return clientIp
		}),
		"user_id", log.Valuer(func(ctx context.Context) interface{} {
			userId, _ := UserIdFrom(ctx)
			return userId
		}),
		"request_id", log.Valuer(func(ctx context.Context) interface{} {
			reqId, _ := RequestIdFrom(ctx)
			return reqId
		}),
		"namespace", log.Valuer(func(ctx context.Context) interface{} {
			return os.Getenv("NAMESPACE")
		}),
	}
}
