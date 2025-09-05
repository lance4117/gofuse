package fsconfig

import (
	"gitee.com/lance4117/GoFuse/fserror"
	"gitee.com/lance4117/GoFuse/fslogger"
	"gitee.com/lance4117/GoFuse/fsonce"
	"github.com/spf13/viper"
)

var cfg *viper.Viper

var InitConfig = fsonce.DoWithParam(func(path string) struct{} {
	if path == "" {
		path = "./config.yaml"
	}
	cfg = viper.New()
	cfg.SetConfigFile(path)
	if err := cfg.ReadInConfig(); err != nil {
		fslogger.Panic(fserror.ErrConfigRead, err)
	}
	return struct{}{}
})

// LoadKey 通过key获取配置结构体
func LoadKey(key string, configStru interface{}) error {
	if cfg == nil {
		fslogger.Error(fserror.ErrConfigLoad, configStru)
		return fserror.ErrConfigLoad
	}
	if err := cfg.UnmarshalKey(key, configStru); err != nil {
		return err
	}
	return nil
}

// GetString 通过key访问配置
func GetString(key string) string {
	if cfg == nil {
		fslogger.Error(fserror.ErrConfigLoad, key)
		return ""
	}
	return cfg.GetString(key)
}

// GetInt 通过key访问配置
func GetInt(key string) int {
	if cfg == nil {
		fslogger.Error(fserror.ErrConfigLoad, key)
		return 0
	}
	return cfg.GetInt(key)
}

// GetBool 通过key访问配置
func GetBool(key string) bool {
	if cfg == nil {
		fslogger.Error(fserror.ErrConfigLoad, key)
		return false
	}
	return cfg.GetBool(key)
}

// All 获取全部配置
func All() map[string]interface{} {
	return cfg.AllSettings()
}
