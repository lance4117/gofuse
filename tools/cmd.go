package tools

import (
	"bytes"
	"errors"
	"os/exec"
)

// ExactCmd 执行一个外部命令
// 注意，这里的arg必须要分开
// 这里并不会像 shell 那样帮你解析一整条字符串
func ExactCmd(name string, arg ...string) ([]byte, error) {
	var out bytes.Buffer

	cmd := exec.Command(name, arg...)

	out.Reset()
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(out.String())
	}
	return out.Bytes(), nil
}
