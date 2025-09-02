package fxconfig

import (
	"gofuse/fxerror"
	"gofuse/fxlogger"
	"gofuse/fxonce"

	"github.com/spf13/viper"
)

var cfg *viper.Viper

var InitConfig = fxonce.DoWithParam(func(path string) struct{} {
	if path == "" {
		path = "./fxconfig.yaml"
	}
	cfg = viper.New()
	cfg.SetConfigFile(path)
	if err := cfg.ReadInConfig(); err != nil {
		fxlogger.Panic(fxerror.ErrConfigRead, err)
	}
	return struct{}{}
})

// LoadKey 通过key获取配置结构体
func LoadKey(key string, configStru interface{}) error {
	if cfg == nil {
		fxlogger.Error(fxerror.ErrConfigLoad, configStru)
		return fxerror.ErrConfigLoad
	}
	if err := cfg.UnmarshalKey(key, configStru); err != nil {
		return err
	}
	return nil
}

// GetString 通过key访问配置
func GetString(key string) string {
	if cfg == nil {
		fxlogger.Error(fxerror.ErrConfigLoad, key)
		return ""
	}
	return cfg.GetString(key)
}

// GetInt 通过key访问配置
func GetInt(key string) int {
	if cfg == nil {
		fxlogger.Error(fxerror.ErrConfigLoad, key)
		return 0
	}
	return cfg.GetInt(key)
}

// GetBool 通过key访问配置
func GetBool(key string) bool {
	if cfg == nil {
		fxlogger.Error(fxerror.ErrConfigLoad, key)
		return false
	}
	return cfg.GetBool(key)
}

// All 获取全部配置
func All() map[string]interface{} {
	return cfg.AllSettings()
}
