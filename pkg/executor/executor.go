package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

// CommandExecutor 是命令执行器的接口
type CommandExecutor interface {
	Execute(command string) error
}

// ShellExecutor 是基于shell的命令执行器
type ShellExecutor struct {
	Stdout io.Writer
	Stderr io.Writer
}

// NewShellExecutor 创建一个新的shell命令执行器
func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Execute 执行shell命令
func (e *ShellExecutor) Execute(command string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// 在Windows上检查是否是PowerShell
		if os.Getenv("PSModulePath") != "" {
			cmd = exec.Command("powershell", "-Command", command)
		} else {
			cmd = exec.Command("cmd", "/C", command)
		}
	default:
		// 对于Unix系统（Linux和macOS），使用默认的shell
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
		cmd = exec.Command(shell, "-c", command)
	}

	cmd.Stdout = e.Stdout
	cmd.Stderr = e.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing command: %v", err)
	}

	return nil
}
