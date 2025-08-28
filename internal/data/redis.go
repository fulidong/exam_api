// internal/data/redis_client.go

package data

import (
	"context"
	"exam_api/internal/biz"
	"exam_api/internal/conf"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisClient 包装 *redis.Client
type RedisClient struct {
	client *redis.Client
	log    *log.Helper
}

func NewRedisClient(c *conf.Data, logger log.Logger) *RedisClient {
	logHelper := log.NewHelper(logger)

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Database),
		PoolSize:     int(c.Redis.PoolSize),
		MinIdleConns: int(c.Redis.MinIdleConns),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})

	// 连接测试
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		logHelper.Fatalf("Redis connect failed: %v", err)
	}
	logHelper.Info("Connected to Redis")

	return &RedisClient{client: client, log: logHelper}
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// internal/data/redis_repo.go

var _ biz.RedisRepository = (*redisRepository)(nil) // 确保实现接口

type redisRepository struct {
	client *redis.Client
	log    *log.Helper
}

func NewRedisRepository(client *RedisClient, logger log.Logger) biz.RedisRepository {
	return &redisRepository{
		client: client.GetClient(),
		log:    log.NewHelper(logger),
	}
}

// === String 操作 ===

func (r *redisRepository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisRepository) Del(ctx context.Context, keys ...string) error {
	_, err := r.client.Del(ctx, keys...).Result()
	return err
}

func (r *redisRepository) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

func (r *redisRepository) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *redisRepository) SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

// === 数值操作 ===

func (r *redisRepository) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *redisRepository) IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	return r.client.IncrBy(ctx, key, increment).Result()
}

func (r *redisRepository) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

func (r *redisRepository) DecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	return r.client.DecrBy(ctx, key, decrement).Result()
}

// === Hash 操作 ===

func (r *redisRepository) HSet(ctx context.Context, key, field, value string) error {
	return r.client.HSet(ctx, key, field, value).Err()
}

func (r *redisRepository) HGet(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

func (r *redisRepository) HDel(ctx context.Context, key string, fields ...string) error {
	_, err := r.client.HDel(ctx, key, fields...).Result()
	return err
}

func (r *redisRepository) HExists(ctx context.Context, key, field string) (bool, error) {
	return r.client.HExists(ctx, key, field).Result()
}

func (r *redisRepository) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// === 通用操作 ===

func (r *redisRepository) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *redisRepository) Keys(ctx context.Context, pattern string) ([]string, error) {

	return r.client.Keys(ctx, pattern).Result()
}

func (r *redisRepository) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(ctx, script, keys, args).Result()
}

// === Pipeline 实现 ===

type redisPipeline struct {
	pipe redis.Pipeliner // 持有 go-redis 的 Pipeliner
}

// Set 实现
func (p *redisPipeline) Set(key string, value string, expiration time.Duration) {
	p.pipe.Set(context.Background(), key, value, expiration)
}

// Get 实现
func (p *redisPipeline) Get(key string) {
	p.pipe.Get(context.Background(), key)
}

// Del 实现
func (p *redisPipeline) Del(keys ...string) {
	p.pipe.Del(context.Background(), keys...)
}

// HSet 实现
func (p *redisPipeline) HSet(key, field, value string) {
	p.pipe.HSet(context.Background(), key, field, value)
}

// HGet 实现
func (p *redisPipeline) HGet(key, field string) {
	p.pipe.HGet(context.Background(), key, field)
}

// HDel 实现
func (p *redisPipeline) HDel(key string, fields ...string) {
	p.pipe.HDel(context.Background(), key, fields...)
}

// Exec 执行所有命令并返回结果
func (p *redisPipeline) Exec(ctx context.Context) ([]biz.RedisCmdResult, error) {
	// 在传入的 ctx 下执行
	cmds, err := p.pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	results := make([]biz.RedisCmdResult, len(cmds))
	for i, cmd := range cmds {
		var value interface{}
		var cmdErr error

		// 根据命令类型提取结果
		switch c := cmd.(type) {
		case *redis.StringCmd:
			value, cmdErr = c.Result()
		case *redis.IntCmd:
			value, cmdErr = c.Result()
		case *redis.BoolCmd:
			value, cmdErr = c.Result()
		case *redis.SliceCmd:
			value, cmdErr = c.Result()
		case *redis.StatusCmd:
			value, cmdErr = c.Result()
		default:
			value = nil
			cmdErr = fmt.Errorf("unsupported command type: %T", c)
		}

		results[i] = biz.RedisCmdResult{
			Cmd:   cmd.Name(),
			Args:  cmd.Args(),
			Value: value,
			Err:   cmdErr,
		}
	}

	return results, nil
}

// === RedisRepository 接口方法 ===

func (r *redisRepository) Pipeline(ctx context.Context) biz.RedisPipeline {
	// 使用原生 Pipeline
	pipe := r.client.Pipeline()
	return &redisPipeline{pipe: pipe}
}
