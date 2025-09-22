package monitor

import (
	"testing"

	"gitee.com/lance4117/GoFuse/logger"
)

func TestPsutils(t *testing.T) {
	pid := 28828
	monitor := NewDefaultMonitor(pid, "C:\\")
	if err := monitor.Run(); err != nil {
		logger.Fatal(err)
	}
}

func TestNamePid(t *testing.T) {
	pid, err := GetPidByName("blogd.exe")
	t.Log(pid, err)
}
