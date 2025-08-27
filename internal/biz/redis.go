// internal/biz/redis_repo.go

package biz

import (
	"context"
	"time"
)

// RedisRepository 是缓存操作的抽象接口
// 由 data 层实现，biz 层只依赖此接口
type RedisRepository interface {
	// String 操作
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)

	// 数值操作
	Incr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, increment int64) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	DecrBy(ctx context.Context, key string, decrement int64) (int64, error)

	// Hash 操作
	HSet(ctx context.Context, key, field, value string) error
	HGet(ctx context.Context, key, field string) (string, error)
	HDel(ctx context.Context, key string, fields ...string) error
	HExists(ctx context.Context, key, field string) (bool, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)

	// 通用操作
	Exists(ctx context.Context, key string) (bool, error)
	Keys(ctx context.Context, pattern string) ([]string, error)

	// 事务/批处理（可选）
	Pipeline(ctx context.Context) RedisPipeline
}

// RedisPipeline 表示一个 Redis Pipeline（管道）
type RedisPipeline interface {
	// String
	Set(key string, value string, expiration time.Duration)
	Get(key string)
	Del(keys ...string)

	// Hash
	HSet(key, field, value string)
	HGet(key, field string)
	HDel(key string, fields ...string)

	// 执行所有累积的命令，返回结果
	Exec(ctx context.Context) ([]RedisCmdResult, error)
}

// RedisCmdResult 表示 pipeline 中每条命令的结果
type RedisCmdResult struct {
	// 可以根据需要记录命令类型、key、value、error 等
	Cmd   string        // 命令名，如 "SET", "GET"
	Args  []interface{} // 命令参数
	Value interface{}   // 成功时的返回值
	Err   error         // 错误
}
