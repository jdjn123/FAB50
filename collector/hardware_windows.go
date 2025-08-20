//go:build windows

package collector

import (
	"os/exec"
	"syscall"
)

const (
	CREATE_NO_WINDOW = 0x08000000
)

// setWindowsProcessAttributes 设置Windows进程属性
func setWindowsProcessAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: CREATE_NO_WINDOW,
	}
}
