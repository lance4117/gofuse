package test

import (
	"testing"

	"gitee.com/lance4117/GoFuse/config"
	"gitee.com/lance4117/GoFuse/logger"
)

func TestGetConfig(t *testing.T) {
	config.Init("")

	logger.Info(config.All())

	logger.Info(config.GetString("app.name"))
}
