package config

import (
	"gitee.com/lance4117/GoFuse/errs"
	"gitee.com/lance4117/GoFuse/logger"
	"gitee.com/lance4117/GoFuse/once"
	"github.com/spf13/viper"
)

var cfg *viper.Viper

var InitConfig = once.DoWithParam(func(path string) struct{} {
	if path == "" {
		path = "./config.yaml"
	}
	cfg = viper.New()
	cfg.SetConfigFile(path)
	if err := cfg.ReadInConfig(); err != nil {
		logger.Panic(errs.ErrConfigRead, err)
	}
	return struct{}{}
})

// LoadKey 通过key获取配置结构体
func LoadKey(key string, configStru interface{}) error {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad, configStru)
		return errs.ErrConfigLoad
	}
	if err := cfg.UnmarshalKey(key, configStru); err != nil {
		return err
	}
	return nil
}

// GetString 通过key访问配置
func GetString(key string) string {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad, key)
		return ""
	}
	return cfg.GetString(key)
}

// GetInt 通过key访问配置
func GetInt(key string) int {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad, key)
		return 0
	}
	return cfg.GetInt(key)
}

// GetBool 通过key访问配置
func GetBool(key string) bool {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad, key)
		return false
	}
	return cfg.GetBool(key)
}

// All 获取全部配置
func All() map[string]interface{} {
	return cfg.AllSettings()
}
