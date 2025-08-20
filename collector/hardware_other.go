//go:build !windows

package collector

import (
	"os/exec"
)

// setWindowsProcessAttributes 设置Windows进程属性（非Windows系统为空实现）
func setWindowsProcessAttributes(cmd *exec.Cmd) {
	// 非Windows系统不需要特殊设置
}
