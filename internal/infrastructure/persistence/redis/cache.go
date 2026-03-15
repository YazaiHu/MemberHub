package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/YazaiHu/MemberHub/internal/pkg/config"
)

var Client *redis.Client
var ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// 测试连接
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("connect to redis failed: %w", err)
	}

	log.Println("Redis connected successfully")
	return nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return Client
}

// Close 关闭Redis连接
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// --- 常用操作封装 ---

// Set 设置键值
func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func Get(key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// Del 删除键
func Del(keys ...string) error {
	return Client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(keys ...string) (int64, error) {
	return Client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	return Client.Expire(ctx, key, expiration).Err()
}

// Incr 自增
func Incr(key string) (int64, error) {
	return Client.Incr(ctx, key).Result()
}

// Decr 自减
func Decr(key string) (int64, error) {
	return Client.Decr(ctx, key).Result()
}

// IncrBy 增加指定值
func IncrBy(key string, value int64) (int64, error) {
	return Client.IncrBy(ctx, key, value).Result()
}

// DecrBy 减少指定值
func DecrBy(key string, value int64) (int64, error) {
	return Client.DecrBy(ctx, key, value).Result()
}

// SetNX 设置键值（仅当键不存在时）
func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return Client.SetNX(ctx, key, value, expiration).Result()
}

// HSet 设置哈希字段
func HSet(key string, field string, value interface{}) error {
	return Client.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希字段
func HGet(key, field string) (string, error) {
	return Client.HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希所有字段
func HGetAll(key string) (map[string]string, error) {
	return Client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func HDel(key string, fields ...string) error {
	return Client.HDel(ctx, key, fields...).Err()
}

// HIncrBy 哈希字段增加
func HIncrBy(key, field string, incr int64) (int64, error) {
	return Client.HIncrBy(ctx, key, field, incr).Result()
}

// ZAdd 添加有序集合成员
func ZAdd(key string, score float64, member interface{}) error {
	return Client.ZAdd(ctx, key, &redis.Z{Score: score, Member: member}).Err()
}

// ZRangeByScore 按分数范围获取有序集合成员
func ZRangeByScore(key string, min, max string) ([]string, error) {
	return Client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
}

// ZRem 删除有序集合成员
func ZRem(key string, members ...interface{}) error {
	return Client.ZRem(ctx, key, members...).Err()
}

// --- 分布式锁 ---

// Lock 分布式锁
type Lock struct {
	key        string
	value      string
	expiration time.Duration
}

// AcquireLock 获取锁
func AcquireLock(key string, expiration time.Duration) (*Lock, bool, error) {
	value := fmt.Sprintf("%d", time.Now().UnixNano())
	ok, err := SetNX(key, value, expiration)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	return &Lock{
		key:        key,
		value:      value,
		expiration: expiration,
	}, true, nil
}

// ReleaseLock 释放锁
func (l *Lock) Release() error {
	// Lua脚本保证原子性：只有持有锁的客户端才能释放
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	return Client.Eval(ctx, script, []string{l.key}, l.value).Err()
}

// --- 限流 ---

// RateLimiter 限流器（令牌桶算法）
func RateLimiter(key string, limit int, window time.Duration) (bool, error) {
	script := `
		local key = KEYS[1]
		local limit = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])
		local current = redis.call('incr', key)
		if current == 1 then
			redis.call('expire', key, window)
		end
		if current > limit then
			return 0
		else
			return 1
		end
	`
	result, err := Client.Eval(ctx, script, []string{key}, limit, int(window.Seconds())).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}
