package monitor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lance4117/gofuse/errs"
	"github.com/lance4117/gofuse/fileio"
	"github.com/lance4117/gofuse/logger"
	"github.com/lance4117/gofuse/times"
)

func TestMonitor(t *testing.T) {
	appName := "app"
	pids, err := GetPidByName(appName)
	if err != nil {
		logger.Error(err)
		return
	}
	if pids == nil {
		logger.Error(errs.ErrNil)
		return
	}
	pid := pids[0]
	interval := time.Second * 1
	collectors := []Collector{NewCPUCollector(), NewMemoryCollector(),
		NewIOCollector(), NewDiskCollector("D:\\code\\zerod")}

	// 用 ctx 管控整个生命周期：任务结束 -> cancel -> 监控退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()

	// 设置监控程序
	m := NewCustomMonitor(pid, interval, collectors, fileio.NewCSVFileIO(fmt.Sprintf("monitor-%d", times.NowMilli())))
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
	time.Sleep(5 * time.Second)
	return
}

func TestNamePid(t *testing.T) {
	pid, err := GetPidByName("blogd.exe")
	t.Log(pid, err)
}
