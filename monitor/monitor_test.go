package monitor

import (
	"testing"

	"gitee.com/lance4117/GoFuse/logger"
)

func TestPsutils(t *testing.T) {
	pid := 28828
	monitor := NewDefaultMonitor(pid, "/Users/admin/.blog")
	if err := monitor.Run(); err != nil {
		logger.Fatal(err)
	}

}
