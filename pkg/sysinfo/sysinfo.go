package sysinfo

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// SystemInfo 包含系统环境的相关信息
type SystemInfo struct {
	OS         string
	OSVersion  string
	Shell      string
	Username   string
	HomeDir    string
	CurrentDir string
	EnvVars    map[string]string // 环境变量列表
}

// GetSystemInfo 获取当前系统的环境信息
func GetSystemInfo() (*SystemInfo, error) {
	// 获取当前用户信息
	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("error getting current user: %v", err)
	}

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current directory: %v", err)
	}

	// 获取shell类型
	shell := os.Getenv("SHELL")
	if runtime.GOOS == "windows" {
		// 在Windows上检查是否是PowerShell
		if os.Getenv("PSModulePath") != "" {
			shell = "powershell"
		} else {
			shell = "cmd"
		}
	} else {
		// 对于Unix系统，从SHELL环境变量中提取shell名称
		shell = filepath.Base(shell)
	}

	// 获取环境变量
	envVars := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envVars[parts[0]] = parts[1]
		}
	}

	return &SystemInfo{
		OS:         runtime.GOOS,
		OSVersion:  getOSVersion(),
		Shell:      shell,
		Username:   currentUser.Username,
		HomeDir:    currentUser.HomeDir,
		CurrentDir: cwd,
		EnvVars:    envVars,
	}, nil
}

// getOSVersion 获取操作系统版本
func getOSVersion() string {
	switch runtime.GOOS {
	case "darwin":
		// 在macOS上使用sw_vers命令获取版本
		if out, err := os.ReadFile("/System/Library/CoreServices/SystemVersion.plist"); err == nil {
			if idx := strings.Index(string(out), "<string>10."); idx != -1 {
				ver := string(out[idx+8:])
				if idx = strings.Index(ver, "</string>"); idx != -1 {
					return "macOS " + ver[:idx]
				}
			}
		}
		return "macOS"

	case "windows":
		// 在Windows上使用ver命令获取版本
		return os.Getenv("OS")

	case "linux":
		// 在Linux上读取/etc/os-release文件获取版本
		if out, err := os.ReadFile("/etc/os-release"); err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					ver := strings.Trim(line[12:], "\"") // 移除PRETTY_NAME=和引号
					return ver
				}
			}
		}
		return "Linux"

	default:
		return runtime.GOOS
	}
}
