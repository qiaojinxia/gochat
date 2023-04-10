package utils

import (
	"fmt"
	"os/exec"
)

//RunPythonScript 调用 Python 代码的函数
func RunPythonScript(scriptPath string, args ...string) (string, error) {
	// 构造命令行参数
	cmdArgs := append([]string{scriptPath}, args...)

	// 运行 Python 命令行解释器并执行脚本
	cmd := exec.Command("python", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run python script: %v", err)
	}
	// 返回输出结果
	return string(output), nil
}
