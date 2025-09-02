package test

import (
	"gofuse/fxconfig"
	"gofuse/fxlogger"
	"testing"
)

func TestGetConfig(t *testing.T) {
	fxconfig.InitConfig("")

	fxlogger.Info(fxconfig.All())

	fxlogger.Info(fxconfig.GetString("app.name"))
}
