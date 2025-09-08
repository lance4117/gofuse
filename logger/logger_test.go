package logger

import (
	"testing"
)

func TestRunLogger(t *testing.T) {
	Info("hello")
	Warn("warning")
	Debug("debug")
	Error("error")
}

func TestFormatLogger(t *testing.T) {
	Infof("info %s", "format")
	Debugf("debug %s", "format")
	Warnf("warn %s", "format")
	Errorf("error %s", "format")
}