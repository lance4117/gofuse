package kvs

import "github.com/cockroachdb/pebble"

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	PoolSize int
}

type PebbleConfig struct {
	DirName string
	Options *pebble.Options
}

// NewPebbleConfig 创建一个默认配置的 Pebble 配置
func NewPebbleConfig(dirname string) PebbleConfig {
	opts := &pebble.Options{}
	opts.EnsureDefaults()
	return PebbleConfig{
		DirName: dirname,
		Options: opts,
	}
}

// NewRedisConfig  创建一个 Redis 配置
func NewRedisConfig(addr, password string, db, poolSize int) RedisConfig {
	return RedisConfig{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	}
}
