package config

import (
	"time"

	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/logger"
	"github.com/lance4117/gofuse/once"
	"github.com/spf13/viper"
)

var cfg *viper.Viper

// Init 是一个初始化配置的函数，使用 once.DoWithParam 确保只执行一次。
// 参数 path 指定配置文件路径，默认为 "./config.yaml"。
// 如果读取配置失败会调用 logger.Fatal 记录致命错误并终止程序。
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

// Has 判断指定 key 是否存在于配置中。
// 参数 key 表示要查询的配置项名称。
// 返回值表示该配置项是否存在。
func Has(key string) bool {
	return cfg != nil && cfg.IsSet(key)
}

// LoadKey 将指定 key 对应的配置解析到传入的结构体中。
// 参数 key 表示要加载的配置项名称。
// 参数 configStru 是用于接收配置数据的结构体指针。
// 返回值是反序列化过程中可能发生的错误。
func LoadKey(key string, configStru interface{}) error {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), configStru)
		return errs.ErrConfigLoad(key)
	}
	if err := cfg.UnmarshalKey(key, configStru); err != nil {
		return err
	}
	return nil
}

// GetString 获取字符串类型的配置值。
// 参数 key 表示要获取的配置项名称。
// 返回值是对应的字符串值；如果配置未初始化或不存在则记录错误日志并返回空字符串。
func GetString(key string) string {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), key)
		return ""
	}
	return cfg.GetString(key)
}

// GetStringOr 获取字符串类型的配置值，若不存在则返回默认值。
// 参数 key 表示要获取的配置项名称。
// 参数 def 表示当配置项不存在时返回的默认值。
// 返回值是配置值或默认值。
func GetStringOr(key, def string) string {
	if cfg == nil {
		return def
	}
	if v := cfg.GetString(key); v != "" {
		return v
	}
	return def
}

// GetInt 获取整数类型的配置值。
// 参数 key 表示要获取的配置项名称。
// 返回值是对应的整数值；如果配置未初始化或出错则记录错误日志并返回 0。
func GetInt(key string) int {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), key)
		return 0
	}
	return cfg.GetInt(key)
}

// GetIntOr 获取整数类型的配置值，若不存在则返回默认值。
// 参数 key 表示要获取的配置项名称。
// 参数 def 表示当配置项不存在时返回的默认值。
// 返回值是配置值或默认值。
func GetIntOr(key string, def int) int {
	if cfg == nil {
		return def
	}
	if !cfg.IsSet(key) {
		return def
	}
	return cfg.GetInt(key)
}

// GetInt64 获取 64 位整数类型的配置值。
// 参数 key 表示要获取的配置项名称。
// 返回值是对应的 int64 值；如果配置未初始化或出错则记录错误日志并返回 0。
func GetInt64(key string) int64 {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), key)
		return 0
	}
	return cfg.GetInt64(key)
}

// GetInt64Or 获取 64 位整数类型的配置值，若不存在则返回默认值。
// 参数 key 表示要获取的配置项名称。
// 参数 def 表示当配置项不存在时返回的默认值。
// 返回值是配置值或默认值。
func GetInt64Or(key string, def int64) int64 {
	if cfg == nil {
		return def
	}
	if !cfg.IsSet(key) {
		return def
	}
	return cfg.GetInt64(key)
}

// GetUint64 获取无符号 64 位整数类型的配置值。
// 参数 key 表示要获取的配置项名称。
// 返回值是对应的 uint64 值；如果配置未初始化或出错则记录错误日志并返回 0。
func GetUint64(key string) uint64 {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), key)
		return 0
	}
	return cfg.GetUint64(key)
}

// GetUint64Or 获取无符号 64 位整数类型的配置值，若不存在则返回默认值。
// 参数 key 表示要获取的配置项名称。
// 参数 def 表示当配置项不存在时返回的默认值。
// 返回值是配置值或默认值。
func GetUint64Or(key string, def uint64) uint64 {
	if cfg == nil {
		return def
	}
	if !cfg.IsSet(key) {
		return def
	}
	return cfg.GetUint64(key)
}

// GetFloat64 获取浮点数类型的配置值。
// 参数 key 表示要获取的配置项名称。
// 返回值是对应的 float64 值；如果配置未初始化或出错则记录错误日志并返回 0。
func GetFloat64(key string) float64 {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), key)
		return 0
	}
	return cfg.GetFloat64(key)
}

// GetFloat64Or 获取浮点数类型的配置值，若不存在则返回默认值。
// 参数 key 表示要获取的配置项名称。
// 参数 def 表示当配置项不存在时返回的默认值。
// 返回值是配置值或默认值。
func GetFloat64Or(key string, def float64) float64 {
	if cfg == nil {
		return def
	}
	if !cfg.IsSet(key) {
		return def
	}
	return cfg.GetFloat64(key)
}

// GetDuration 获取时间间隔类型的配置值。
// 参数 key 表示要获取的配置项名称。
// 返回值是对应的时间间隔值；如果配置未初始化则记录错误日志并返回 0。
func GetDuration(key string) time.Duration {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key))
		return 0
	}
	return cfg.GetDuration(key)
}

// GetDurationOr 获取时间间隔类型的配置值，若不存在或为零值则返回默认值。
// 参数 key 表示要获取的配置项名称。
// 参数 def 表示当配置项不存在或为零值时返回的默认值。
// 返回值是配置值或默认值。
func GetDurationOr(key string, def time.Duration) time.Duration {
	if cfg == nil || !cfg.IsSet(key) {
		return def
	}
	d := cfg.GetDuration(key)
	if d == 0 {
		return def
	}
	return d
}

// GetBool 获取布尔类型的配置值。
// 参数 key 表示要获取的配置项名称。
// 返回值是对应的布尔值；如果配置未初始化则记录错误日志并返回 false。
func GetBool(key string) bool {
	if cfg == nil {
		logger.Error(errs.ErrConfigLoad(key), key)
		return false
	}
	return cfg.GetBool(key)
}

// All 获取所有配置项。
// 返回值是一个包含所有配置键值对的映射。
func All() map[string]any {
	return cfg.AllSettings()
}
