package test

import (
	"gofuse/config"
	"gofuse/logger"
	"testing"
)

func TestGetConfig(t *testing.T) {
	config.InitConfig("")

	logger.Info(config.All())

	logger.Info(config.GetString("app.name"))
}
