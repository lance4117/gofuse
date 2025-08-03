package test

import (
	"gofuse/config"
	"testing"
)

func TestGetConfig(t *testing.T) {
	cfg := config.GetConfig("./config.yaml")
	t.Log(cfg.AllKeys())
	t.Log(cfg.All())
}
