package test

import (
	"fmt"
	"os/exec"
	"syscall"
	"testing"
)

func Test_ExcPython(t *testing.T) {
	args1 := "{\"table0\":\"data/2NsfmElfZbjbM1X3Zt5raLdI3f4.xlsx\"}"
	py := "C:\\Users\\a\\GolandProjects\\gochat\\data\\2NmcBlwuqPslOlavQOccJSvdDqE.py"
	// 执行 Python 文件
	cmd := exec.Command("python", py)
	output, err := cmd.Output()
	if err != nil {
		// 调用脚本出错
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			fmt.Println("无法获取退出状态码:", err)
			return
		}
		// 获取退出状态码
		ws := exitErr.ProcessState.Sys().(syscall.WaitStatus)
		fmt.Printf("脚本退出状态码: %d\n", ws.ExitStatus())
		return
	}
	fmt.Println(string(output))
}
