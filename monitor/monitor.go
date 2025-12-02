package monitor

import (
	"context"
	"time"

	"github.com/lance4117/gofuse/fileio"
	"github.com/lance4117/gofuse/logger"
	"github.com/shirou/gopsutil/v4/process"
)

// Monitor 监控器
type Monitor struct {
	pid        int
	interval   time.Duration
	duration   time.Duration
	collectors []Collector
	writer     fileio.Writer
}

// NewCustomMonitor 通用构造函数
func NewCustomMonitor(pid int, interval time.Duration, collectors []Collector, writer fileio.Writer) *Monitor {
	return &Monitor{
		pid:        pid,
		interval:   interval,
		collectors: collectors,
		writer:     writer,
	}
}

// NewDefaultMonitor 默认构造函数
func NewDefaultMonitor(pid int, path string) *Monitor {
	writer, err := fileio.NewCSVWriter(path)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return &Monitor{
		pid:        pid,
		interval:   time.Second,
		duration:   time.Minute,
		collectors: []Collector{NewCPUCollector(), NewMemoryCollector(), NewIOCollector(), NewDiskCollector(path), NewNetCollector()},
		writer:     writer,
	}
}

// Run 执行监控
// 用 ctx 管控整个生命周期
func (m *Monitor) Run(ctx context.Context, showLog bool) error {
	logger.Infof("monitor started, pid=%d, interval=%s, duration=%s", m.pid, m.interval, m.duration)

	p, err := process.NewProcess(int32(m.pid))
	if err != nil {
		return err
	}

	// 准备表头
	headers := []string{"Time"}
	for _, c := range m.collectors {
		headers = append(headers, c.Names()...)
	}
	if err := m.writer.Write(headers); err != nil {
		return err
	}
	defer m.writer.Close()

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Infof("monitor finished, pid=%d", m.pid)
			return nil
		case now := <-ticker.C:
			row := []string{now.Format(time.DateTime)}
			for _, c := range m.collectors {
				val, err := c.Collect(p, now)
				if err != nil {
					row = append(row, "ERR")
					logger.Error("collector %s failed: %v", c.Names(), err)
				} else {
					row = append(row, val...)
				}
			}
			if err := m.writer.Write(row); err != nil {
				return err
			}
			if showLog {
				logger.Infof("recorded sample for pid=%d at %s", m.pid, now.Format(time.DateTime))
			}
		}
	}
}
