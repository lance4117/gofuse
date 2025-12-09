package monitor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/lance4117/gofuse/fileio"
	"github.com/lance4117/gofuse/logger"
	"github.com/lance4117/gofuse/times"
)

func TestMonitor(t *testing.T) {
	pidStr := os.Getenv("TEST_MONITOR_PID")
	appName := os.Getenv("TEST_MONITOR_NAME")
	var pid int

	switch {
	case pidStr != "":
		parsed, err := strconv.Atoi(pidStr)
		if err != nil {
			t.Skipf("invalid TEST_MONITOR_PID: %v", err)
		}
		pid = parsed
	case appName != "":
		pids, err := GetPidByName(appName)
		if err != nil || len(pids) == 0 {
			t.Skipf("no pid found for %s: %v", appName, err)
		}
		pid = pids[0]
	default:
		t.Skip("set TEST_MONITOR_PID or TEST_MONITOR_NAME to run monitor integration test")
	}

	interval := time.Second * 1
	collectors := []Collector{NewCPUCollector(), NewMemoryCollector(),
		NewIOCollector(), NewDiskCollector(t.TempDir())}

	// 用 ctx 管控整个生命周期：任务结束 -> cancel -> 监控退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()

	tempFile := filepath.Join(t.TempDir(), fmt.Sprintf("monitor-%d", times.NowMilli()))
	writer, err := fileio.NewCSVWriter(tempFile)
	if err != nil {
		logger.Error(err)
		return
	}
	// 设置监控程序
	m := NewCustomMonitor(pid, interval, collectors, writer)
	go func() {
		err := m.Run(ctx, false)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// 经过Wait函数后调用cancel关闭ctx从而结束monitor
	Wait()
	cancel()

	// 等待上下文超时
	<-ctx.Done()
	logger.Info("时间结束，总共执行了", time.Since(start).Minutes(), "分钟")
}

func Wait() {
	time.Sleep(2 * time.Second)
}

func TestNamePid(t *testing.T) {
	pid, err := GetPidByName("blogd.exe")
	t.Log(pid, err)
}
