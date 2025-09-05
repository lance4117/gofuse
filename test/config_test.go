package test

import (
	"testing"

	"gitee.com/lance4117/GoFuse/fsconfig"
	"gitee.com/lance4117/GoFuse/fslogger"
)

func TestGetConfig(t *testing.T) {
	fsconfig.InitConfig("")

	fslogger.Info(fsconfig.All())

	fslogger.Info(fsconfig.GetString("app.name"))
}
