package rediskv

import (
	"context"
	"errors"
	"time"

	"github.com/lance4117/gofuse/errs"
	"github.com/redis/go-redis/v9"
)

// RedisStore 封装了 redis.Client 实例，提供键值存储的基本操作接口。
type RedisStore struct {
	RedisCli *redis.Client
}

// Put 写入键值（持久存储）
func (r *RedisStore) Put(key string, val []byte) error {
	return r.RedisCli.Set(context.Background(), key, val, 0).Err()
}

// PutWithTTL 写入带有效期的键值
func (r *RedisStore) PutWithTTL(key string, val []byte, ttl time.Duration) error {
	return r.RedisCli.Set(context.Background(), key, val, ttl).Err()
}

// Get 获取键值，如果 key 不存在返回 ErrKeyNotFound
func (r *RedisStore) Get(key string) ([]byte, error) {
	get := r.RedisCli.Get(context.Background(), key)
	err := get.Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errs.ErrKeyNotFound
		}
		return nil, err
	}
	return []byte(get.Val()), nil
}

// Del 删除键
func (r *RedisStore) Del(key string) error {
	return r.RedisCli.Del(context.Background(), key).Err()
}

// Has 判断键是否存在
func (r *RedisStore) Has(key string) (bool, error) {
	exists, err := r.RedisCli.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// Close 关闭客户端连接
func (r *RedisStore) Close() error {
	return r.RedisCli.Close()
}
