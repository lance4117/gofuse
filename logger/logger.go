package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once          sync.Once
	sugaredLogger *zap.SugaredLogger
	initialized   bool
)

// Config 日志配置
type Config struct {
	Production bool          // 是否为生产环境
	Level      zapcore.Level // 日志级别
	OutputPath string        // 输出路径，空则输出到 stdout
	CallerSkip int           // 调用栈跳过层数，默认为 1
}

// DefaultConfig 返回默认配置（开发模式）
func DefaultConfig() Config {
	return Config{
		Production: false,
		Level:      zapcore.DebugLevel,
		OutputPath: "",
		CallerSkip: 1,
	}
}

// Init 初始化日志记录器
// 如果已经初始化过，再次调用不会重新初始化
func Init(config Config) error {
	var err error
	once.Do(func() {
		var logger *zap.Logger

		// 根据环境选择不同的配置
		if config.Production {
			zapConfig := zap.NewProductionConfig()
			zapConfig.Level = zap.NewAtomicLevelAt(config.Level)
			if config.OutputPath != "" {
				zapConfig.OutputPaths = []string{config.OutputPath}
			}
			logger, err = zapConfig.Build(
				zap.AddCaller(),
				zap.AddCallerSkip(config.CallerSkip),
			)
		} else {
			zapConfig := zap.NewDevelopmentConfig()
			zapConfig.Level = zap.NewAtomicLevelAt(config.Level)
			if config.OutputPath != "" {
				zapConfig.OutputPaths = []string{config.OutputPath}
			}
			logger, err = zapConfig.Build(
				zap.AddCaller(),
				zap.AddCallerSkip(config.CallerSkip),
			)
		}

		if err != nil {
			return
		}

		sugaredLogger = logger.Sugar()
		initialized = true
	})
	return err
}

// MustInit 初始化日志记录器，如果失败则 panic
func MustInit(config Config) {
	if err := Init(config); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

// ensureInitialized 确���日志已初始化，如果未初始化则使用默认配置
func ensureInitialized() {
	if !initialized {
		_ = Init(DefaultConfig())
	}
}

// Info 记录INFO级别的日志信息
func Info(args ...interface{}) {
	ensureInitialized()
	sugaredLogger.Info(args...)
}

// Infof 记录INFO级别的日志信息
func Infof(template string, args ...interface{}) {
	ensureInitialized()
	sugaredLogger.Infof(template, args...)
}

// Debug 记录DEBUG级别的日志信息
func Debug(args ...interface{}) {
	ensureInitialized()
	sugaredLogger.Debug(args...)
}

// Debugf 记录DEBUG级别的日志信息
func Debugf(template string, args ...interface{}) {
	ensureInitialized()
	sugaredLogger.Debugf(template, args...)
}

// Warn 记录WARN级别的日志信息
func Warn(args ...interface{}) {
	ensureInitialized()
	sugaredLogger.Warn(args...)
}

// Warnf 记录WARN级别的日志信息
func Warnf(template string, args ...interface{}) {
	ensureInitialized()
	sugaredLogger.Warnf(template, args...)
}

// Error 记录ERROR级别的日志信息
func Error(args ...interface{}) {
	ensureInitialized()
	sugaredLogger.Error(args...)
}

// Errorf 记录ERROR级别的日志信息
func Errorf(template string, args ...interface{}) {
	ensureInitialized()
	sugaredLogger.Errorf(template, args...)
}

// Panic 记录PANIC级别的日志信息并引发panic
func Panic(args ...interface{}) {
	ensureInitialized()
	sugaredLogger.Panic(args...)
}

// Panicf 记录PANIC级别的日志信息并引发panic
func Panicf(template string, args ...interface{}) {
	ensureInitialized()
	sugaredLogger.Panicf(template, args...)
}

// Fatal 记录FATAL级别的日志信息并退出程序
func Fatal(args ...interface{}) {
	ensureInitialized()
	sugaredLogger.Fatal(args...)
}

// Fatalf 记录FATAL级别的日志信息并退出程序
func Fatalf(template string, args ...interface{}) {
	ensureInitialized()
	sugaredLogger.Fatalf(template, args...)
}
