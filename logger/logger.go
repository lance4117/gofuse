package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

var (
	once          sync.Once
	sugaredLogger *zap.SugaredLogger
)

// Init 初始化日志记录器
func init() {
	once.Do(func() {
		logger, err := zap.NewDevelopment()
		if err != nil {
			log.Fatal(err, "Init Logger")
			return
		}
		sugaredLogger = logger.Sugar()
	})
}

// Info 记录INFO级别的日志信息
func Info(args ...interface{}) {
	sugaredLogger.Info(args...)
}

// Debug 记录DEBUG级别的日志信息
func Debug(args ...interface{}) {
	sugaredLogger.Debug(args...)
}

// Warn 记录WARN级别的日志信息
func Warn(args ...interface{}) {
	sugaredLogger.Warn(args...)
}

// Error 记录ERROR级别的日志信息
func Error(args ...interface{}) {
	sugaredLogger.Error(args...)
}

// Panic 记录PANIC级别的日志信息并引发panic
func Panic(args ...interface{}) {
	sugaredLogger.Panic(args...)
}

// Fatal 记录FATAL级别的日志信息并退出程序
func Fatal(args ...interface{}) {
	sugaredLogger.Fatal(args...)
}
