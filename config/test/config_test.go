package test

import (
	"testing"

	"github.com/lance4117/gofuse/config"
	"github.com/lance4117/gofuse/logger"
)

type app struct {
	Name    string  `json:"name"`
	Version float64 `json:"version"`
}

func TestGetConfig(t *testing.T) {
	config.Init("./config.yaml")

	logger.Info(config.All())

	logger.Info(config.GetString("app.name"))

	key, err := config.LoadKey[app]("app")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(key)
}
