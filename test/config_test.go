package test

import (
	"gofuse/fsconfig"
	"gofuse/fslogger"
	"testing"
)

func TestGetConfig(t *testing.T) {
	fsconfig.InitConfig("")

	fslogger.Info(fsconfig.All())

	fslogger.Info(fsconfig.GetString("app.name"))
}
