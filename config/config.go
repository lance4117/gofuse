package config

import (
	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/logger"
	"github.com/lance4117/gofuse/once"
	"github.com/spf13/viper"
)

var cfg *viper.Viper

var Init = once.DoWithParam(func(path string) struct{} {
	if path == "" {
		path = "./config.yaml"
	}
	cfg = viper.New()
	cfg.SetConfigFile(path)
	if err := cfg.ReadInConfig(); err != nil {
		logger.Fatal(errs.ErrConfigRead, err)
	}
	return struct{}{}
})

// LoadKey 通过key获取配置结构体
func LoadKey[T any](key string) (T, error) {
	var ret T
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key))
		return ret, errs.ErrConfigLoad(key)
	}
	err := cfg.UnmarshalKey(key, &ret)
	return ret, err
}

// GetString 通过key访问配置
func GetString(key string) string {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), key)
		return ""
	}
	ret := cfg.GetString(key)
	if ret == "" {
		logger.Warn(errs.WarnConfigLoadNil)
	}
	return ret
}

// GetInt 通过key访问配置
func GetInt(key string) int {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), key)
		return 0
	}
	ret := cfg.GetInt(key)
	if ret == 0 {
		logger.Warn(errs.WarnConfigLoadNil)
	}
	return ret
}

// GetInt64 通过key访问配置
func GetInt64(key string) int64 {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), key)
		return 0
	}
	ret := cfg.GetInt64(key)
	if ret == 0 {
		logger.Warn(errs.WarnConfigLoadNil)
	}
	return ret
}

// GetUint64 通过key访问配置
func GetUint64(key string) uint64 {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), key)
		return 0
	}
	ret := cfg.GetUint64(key)
	if ret == 0 {
		logger.Warn(errs.WarnConfigLoadNil)
	}
	return ret
}

// GetFloat64 通过key访问配置
func GetFloat64(key string) float64 {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), key)
		return 0
	}
	ret := cfg.GetFloat64(key)
	if ret == 0 {
		logger.Warn(errs.WarnConfigLoadNil)
	}
	return ret
}

// GetBool 通过key访问配置
func GetBool(key string) bool {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), key)
		return false
	}
	return cfg.GetBool(key)
}

// All 获取全部配置
func All() map[string]interface{} {
	return cfg.AllSettings()
}
