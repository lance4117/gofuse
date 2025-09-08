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
		logger, err := zap.NewDevelopment(
			// 调用栈深度
			zap.AddCaller(),
			// 跳过封装的一层
			zap.AddCallerSkip(1),
		)
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

// Infof 记录INFO级别的日志信息
func Infof(template string, args ...interface{}) {
	sugaredLogger.Infof(template, args...)
}

// Debug 记录DEBUG级别的日志信息
func Debug(args ...interface{}) {
	sugaredLogger.Debug(args...)
}

// Debugf 记录DEBUG级别的日志信息
func Debugf(template string, args ...interface{}) {
	sugaredLogger.Debugf(template, args...)
}

// Warn 记录WARN级别的日志信息
func Warn(args ...interface{}) {
	sugaredLogger.Warn(args...)
}

// Warnf 记录WARN级别的日志信息
func Warnf(template string, args ...interface{}) {
	sugaredLogger.Warnf(template, args...)
}

// Error 记录ERROR级别的日志信息
func Error(args ...interface{}) {
	sugaredLogger.Error(args...)
}

// Errorf 记录ERROR级别的日志信息
func Errorf(template string, args ...interface{}) {
	sugaredLogger.Errorf(template, args...)
}

// Panic 记录PANIC级别的日志信息并引发panic
func Panic(args ...interface{}) {
	sugaredLogger.Panic(args...)
}

// Panicf 记录PANIC级别的日志信息并引发panic
func Panicf(template string, args ...interface{}) {
	sugaredLogger.Panicf(template, args...)
}

// Fatal 记录FATAL级别的日志信息并退出程序
func Fatal(args ...interface{}) {
	sugaredLogger.Fatal(args...)
}

// Fatalf 记录FATAL级别的日志信息并退出程序
func Fatalf(template string, args ...interface{}) {
	sugaredLogger.Fatalf(template, args...)
}