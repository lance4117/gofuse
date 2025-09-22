package monitor

import (
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

// GetPidByName 需要提供完整程序名字 rg:"foo.exe"
func GetPidByName(name string) ([]int, error) {
	pids, err := process.Pids()
	if err != nil {
		return nil, err
	}
	var result []int
	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			continue // 可能进程刚好退出，忽略
		}

		exe, err := proc.Name()
		if err != nil {
			continue
		}

		if strings.EqualFold(exe, name) { // 不区分大小写匹配
			result = append(result, int(pid))
		}
	}
	return result, nil
}
