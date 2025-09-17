package system

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// ScanFloat 获取输入的一个数字
func ScanFloat() (float64, error) {
	input, err := ScanStr()
	if err != nil {
		return 0, err
	}
	// 去除前后空白字符和换行
	input = strings.TrimSpace(input)

	// 尝试转换为float64
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// ScanStr 获取输入的字符
func ScanStr() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}
