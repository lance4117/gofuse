package fslogger

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

func Info(args ...interface{}) {
	sugaredLogger.Info(args...)
}

func Debug(args ...interface{}) {
	sugaredLogger.Debug(args...)
}

func Warn(args ...interface{}) {
	sugaredLogger.Warn(args...)
}

func Error(args ...interface{}) {
	sugaredLogger.Error(args...)
}

func Panic(args ...interface{}) {
	sugaredLogger.Panic(args...)
}

func Fatal(args ...interface{}) {
	sugaredLogger.Fatal(args...)
}
