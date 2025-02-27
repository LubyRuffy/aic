package sysinfo

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestGetSystemInfo(t *testing.T) {
	// 保存原始环境变量
	originalShell := os.Getenv("SHELL")
	originalPSModulePath := os.Getenv("PSModulePath")
	defer func() {
		os.Setenv("SHELL", originalShell)
		os.Setenv("PSModulePath", originalPSModulePath)
	}()

	// 测试Unix系统shell检测
	testCases := []struct {
		name     string
		shell    string
		expected string
	}{
		{"bash shell", "/bin/bash", "bash"},
		{"zsh shell", "/usr/bin/zsh", "zsh"},
		{"fish shell", "/usr/local/bin/fish", "fish"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if runtime.GOOS != "windows" {
				os.Setenv("SHELL", tc.shell)
				info, err := GetSystemInfo()
				if err != nil {
					t.Fatalf("GetSystemInfo() error = %v", err)
				}
				if info.Shell != tc.expected {
					t.Errorf("Shell = %v, want %v", info.Shell, tc.expected)
				}
			}
		})
	}

	// 测试Windows系统shell检测
	if runtime.GOOS == "windows" {
		// 测试PowerShell
		t.Run("windows powershell", func(t *testing.T) {
			os.Setenv("PSModulePath", "some-path")
			info, err := GetSystemInfo()
			if err != nil {
				t.Fatalf("GetSystemInfo() error = %v", err)
			}
			if info.Shell != "powershell" {
				t.Errorf("Shell = %v, want powershell", info.Shell)
			}
		})

		// 测试CMD
		t.Run("windows cmd", func(t *testing.T) {
			os.Setenv("PSModulePath", "")
			info, err := GetSystemInfo()
			if err != nil {
				t.Fatalf("GetSystemInfo() error = %v", err)
			}
			if info.Shell != "cmd" {
				t.Errorf("Shell = %v, want cmd", info.Shell)
			}
		})
	}

	// 测试基本系统信息
	info, err := GetSystemInfo()
	if err != nil {
		t.Fatalf("GetSystemInfo() error = %v", err)
	}

	// 验证操作系统信息
	if info.OS != runtime.GOOS {
		t.Errorf("OS = %v, want %v", info.OS, runtime.GOOS)
	}

	// 验证当前目录
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	if info.CurrentDir != cwd {
		t.Errorf("CurrentDir = %v, want %v", info.CurrentDir, cwd)
	}
}

func TestGetOSVersion(t *testing.T) {
	// 测试不同操作系统的版本获取
	version := getOSVersion()

	switch runtime.GOOS {
	case "darwin":
		if !strings.HasPrefix(version, "macOS") {
			t.Errorf("Version = %v, want prefix 'macOS'", version)
		}

	case "windows":
		// 在Windows上验证是否返回了OS环境变量的值
		expected := os.Getenv("OS")
		if version != expected {
			t.Errorf("Version = %v, want %v", version, expected)
		}

	case "linux":
		// 在Linux上验证是否返回了有效的版本字符串
		if version == "" {
			t.Error("Version is empty for Linux")
		}

	default:
		// 对于其他操作系统，验证是否返回GOOS
		if version != runtime.GOOS {
			t.Errorf("Version = %v, want %v", version, runtime.GOOS)
		}
	}
}
