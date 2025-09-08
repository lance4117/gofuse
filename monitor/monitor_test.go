package monitor

import (
	"os"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

func TestPsutils(t *testing.T) {

	pid := 11044

	monitor := NewCustom(Monitor{
		pid:      pid,
		interval: 500 * time.Millisecond, // 缩短采样间隔以更好地测试修复效果
		duration: 10 * time.Second,       // 缩短测试时长
		output:   "test-monitor.csv",
	})

	err := monitor.Run()
	if err != nil {
		t.Fatal(err)
	}

}

// 添加一个辅助测试函数，验证在极短时间内连续调用collect不会出错
func TestCollectRapidly(t *testing.T) {
	pid := os.Getpid()
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		t.Fatal(err)
	}

	// 快速连续调用collect两次，测试修复的除零错误
	_, err = collect(p)
	if err != nil {
		t.Fatal(err)
	}

	// 立即再次调用，这在修复前可能会导致除零错误
	_, err = collect(p)
	if err != nil {
		t.Fatal(err)
	}
}
