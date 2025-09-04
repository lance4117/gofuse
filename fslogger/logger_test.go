package fslogger

import (
	"testing"
)

func TestRunLogger(t *testing.T) {
	Info("hello")
	Warn("warning")
}
