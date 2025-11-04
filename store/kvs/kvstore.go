package kvs

import (
	"context"

	"github.com/cockroachdb/pebble"
	"github.com/lance4117/gofuse/store/kvs/pebblekv"
	"github.com/lance4117/gofuse/store/kvs/rediskv"
	"github.com/redis/go-redis/v9"
)

// KVStore 定义了键值存储接口，提供基本的增删改查操作
type KVStore interface {
	// Put 将指定的键值对存储到数据库中，异步写入
	Put(key string, val []byte) error
	// Get 根据指定的键从数据库中获取对应的值
	Get(key string) ([]byte, error)
	// Del 根据指定的键从数据库中删除对应的键值对
	Del(key string) error
	// Has 检查指定的键是否存在于数据库中
	Has(key string) (bool, error)
	// Close 关闭数据库连接并释放相关资源
	Close() error
}

// NewPebbleKV 根据给定的配置创建一个新的 PebbleKV 实例。
func NewPebbleKV(config PebbleConfig) (KVStore, error) {
	// 防止空值
	opts := config.Options
	if opts == nil {
		opts = &pebble.Options{}
		opts.EnsureDefaults()
	} else {
		opts.EnsureDefaults()
	}
	// 获取实例
	kv, err := pebble.Open(config.DirName, config.Options)
	if err != nil {
		return nil, err
	}
	return &pebblekv.PebbleKV{PebbleDB: kv}, nil
}

func NewRedisKV(cfg RedisConfig) (KVStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 尝试连接
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &rediskv.RedisStore{RedisCli: client}, nil
}
