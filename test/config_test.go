package test

import (
	"testing"

	"github.com/lance4117/gofuse/config"
	"github.com/lance4117/gofuse/logger"
)

func TestGetConfig(t *testing.T) {
	config.Init("12312312")

	logger.Info(config.All())

	logger.Info(config.GetString("app.name"))
}
