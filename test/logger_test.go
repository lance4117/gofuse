package test

import (
	"gofuse/logger"
	"testing"
)

func TestRunLogger(t *testing.T) {
	logger.Info("hello")
	logger.Warn("warning")
}
