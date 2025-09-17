package monitor

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

// Collector 定义采集器接口
type Collector interface {
	Names() []string
	Collect(p *process.Process, now time.Time) ([]string, error)
}

// CPUCollector CPU 使用率采集器
type CPUCollector struct{}

// NewCPUCollector 初始化CPU采集器
func NewCPUCollector() *CPUCollector {
	return &CPUCollector{}
}

func (c *CPUCollector) Names() []string { return []string{"CPU(%)"} }
func (c *CPUCollector) Collect(p *process.Process, _ time.Time) ([]string, error) {
	percent, err := p.CPUPercent()
	if err != nil {
		return nil, err
	}
	return []string{fmt.Sprintf("%.2f", percent)}, nil
}

// MemoryCollector 内存采集器
type MemoryCollector struct{}

// NewMemoryCollector 初始化内存采集器
func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{}
}

func (m *MemoryCollector) Names() []string { return []string{"Memory(MB)"} }
func (m *MemoryCollector) Collect(p *process.Process, _ time.Time) ([]string, error) {
	info, err := p.MemoryInfo()
	if err != nil {
		return nil, err
	}
	return []string{fmt.Sprintf("%d", info.RSS/1024/1024)}, nil
}

// IOCollector IO 采集器（含速率计算）
type IOCollector struct {
	prevRead  uint64
	prevWrite uint64
	prevTS    time.Time
}

// NewIOCollector 初始化IO采集器(含速率计算）
func NewIOCollector() *IOCollector {
	return &IOCollector{}
}

func (c *IOCollector) Names() []string {
	return []string{"Read(KB/s)", "Write(KB/s)", "Read(KB)", "Write(KB)"}
}

func (c *IOCollector) Collect(p *process.Process, now time.Time) ([]string, error) {
	ioStat, err := p.IOCounters()
	if err != nil {
		return nil, err
	}

	readKB := ioStat.ReadBytes / 1024
	writeKB := ioStat.WriteBytes / 1024

	var readRate, writeRate uint64
	if !c.prevTS.IsZero() {
		deltaT := now.Sub(c.prevTS).Seconds()
		if deltaT > 1e-9 {
			readRate = uint64(float64(readKB-c.prevRead) / deltaT)
			writeRate = uint64(float64(writeKB-c.prevWrite) / deltaT)
		}
	}
	c.prevRead, c.prevWrite, c.prevTS = readKB, writeKB, now

	return []string{
		fmt.Sprintf("%d", readRate),
		fmt.Sprintf("%d", writeRate),
		fmt.Sprintf("%d", readKB),
		fmt.Sprintf("%d", writeKB),
	}, nil
}

// DiskCollector 磁盘采集器（必须初始化path）
type DiskCollector struct {
	path string
}

// NewDiskCollector 初始化磁盘采集器
func NewDiskCollector(path string) *DiskCollector {
	return &DiskCollector{path: path}
}

func (d *DiskCollector) Names() []string {
	return []string{"DiskUsed(MB)", "DiskFree(MB)", "DiskUsed(%)"}
}

func (d *DiskCollector) Collect(_ *process.Process, _ time.Time) ([]string, error) {
	usage, err := disk.Usage(d.path)
	if err != nil {
		return nil, err
	}
	return []string{
		fmt.Sprintf("%d", usage.Used/1024/1024),
		fmt.Sprintf("%d", usage.Free/1024/1024),
		fmt.Sprintf("%.2f", usage.UsedPercent),
	}, nil
}

// NetCollector 网络信息采集器
type NetCollector struct {
	prevRecv, prevSent uint64
	prevTS             time.Time
}

// NewNetCollector 初始化网络采集器
func NewNetCollector() *NetCollector {
	return &NetCollector{}
}

func (c *NetCollector) Names() []string {
	return []string{"NetRecv(KB/s)", "NetSent(KB/s)"}
}

func (c *NetCollector) Collect(_ *process.Process, now time.Time) ([]string, error) {
	io, err := net.IOCounters(false)
	if err != nil || len(io) == 0 {
		return []string{"0", "0"}, nil
	}
	recv := io[0].BytesRecv / 1024
	sent := io[0].BytesSent / 1024

	var recvRate, sentRate uint64
	if !c.prevTS.IsZero() {
		delta := now.Sub(c.prevTS).Seconds()
		if delta > 1e-9 {
			// 使用浮点数计算，然后转换为整数
			recvRate = uint64(float64(recv-c.prevRecv) / delta)
			sentRate = uint64(float64(sent-c.prevSent) / delta)
		}
	}
	c.prevRecv, c.prevSent = recv, sent
	c.prevTS = now

	return []string{
		fmt.Sprintf("%d", recvRate),
		fmt.Sprintf("%d", sentRate),
	}, nil
}
