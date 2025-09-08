package monitor

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"gitee.com/lance4117/GoFuse/logger"
	"gitee.com/lance4117/GoFuse/times"
	"github.com/shirou/gopsutil/v4/process"
)

// Monitor 定义监控参数
type Monitor struct {
	pid      int           // 要监控的进程 PID
	interval time.Duration // 采样间隔
	duration time.Duration // 采样时长
	output   string        // 导出文件路径 (CSV)
}

// Record 单条采集数据
type Record struct {
	Time     string
	CPU      float64
	MemoryMB uint64
	ReadKB   uint64
	WriteKB  uint64
}

// IO快照，用于计算速率
type ioSnapshot struct {
	read  uint64
	write uint64
	ts    time.Time
}

var prev *ioSnapshot

// NewCustom 创建一个自定义设置
func NewCustom(cfg Monitor) *Monitor {
	return &cfg
}

// NewDefault 创建一个只接受pid的默认设置
func NewDefault(pid int) *Monitor {
	cfg := Monitor{
		pid:      pid,
		interval: time.Second,                                             // 每 1 秒采样一次
		duration: 1 * time.Minute,                                         // 采样 1 分钟
		output:   fmt.Sprintf("monitor-%d-%d.csv", pid, times.NowMilli()), // 输出文件名字
	}
	return &cfg
}

// Run 执行监控任务
func (m *Monitor) Run() error {
	logger.Infof("monitor started, pid=%d, output=%s, interval=%s, duration=%s",
		m.pid, m.output, m.interval, m.duration)

	p, err := process.NewProcess(int32(m.pid))
	if err != nil {
		return err
	}

	file, err := os.Create(m.output)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写 CSV 表头
	err = writer.Write([]string{"Time", "CPU(%)", "Memory(MB)", "Read(KB/s)", "Write(KB/s)"})
	if err != nil {
		return err
	}

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	timeout := time.After(m.duration)

	// 清空之前的IO快照数据
	prev = nil

	for {
		select {
		case <-timeout:
			logger.Infof("monitor finished, pid=%d, output=%s", m.pid, m.output)
			return nil
		case <-ticker.C:
			rec, err := collect(p)
			if err != nil {
				logger.Error(err)
				continue
			}
			err = writer.Write([]string{
				rec.Time,
				fmt.Sprintf("%.2f", rec.CPU),
				fmt.Sprintf("%d", rec.MemoryMB),
				fmt.Sprintf("%d", rec.ReadKB),
				fmt.Sprintf("%d", rec.WriteKB),
			})
			if err != nil {
				return err
			}
			writer.Flush()
			logger.Infof("recorded sample for pid=%d at %s", m.pid, rec.Time)
		}
	}
}

// collect 从进程采集一次数据
func collect(p *process.Process) (*Record, error) {
	cpuPercent, err := p.CPUPercent()
	if err != nil {
		return nil, err
	}

	memInfo, err := p.MemoryInfo()
	if err != nil {
		return nil, err
	}

	ioStat, err := p.IOCounters()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	readKB := ioStat.ReadBytes / 1024
	writeKB := ioStat.WriteBytes / 1024

	var readRate, writeRate uint64
	if prev != nil {
		deltaT := now.Sub(prev.ts).Seconds()
		if deltaT >= 1e-9 { // 确保时间差足够大，避免除零错误
			readRate = uint64(float64(readKB-prev.read) / deltaT)
			writeRate = uint64(float64(writeKB-prev.write) / deltaT)
		}
	}
	prev = &ioSnapshot{read: readKB, write: writeKB, ts: now}

	return &Record{
		Time:     now.Format(time.DateTime),
		CPU:      cpuPercent,
		MemoryMB: memInfo.RSS / 1024 / 1024,
		ReadKB:   readRate,
		WriteKB:  writeRate,
	}, nil
}