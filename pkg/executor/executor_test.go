package executor

import (
	"bytes"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestNewShellExecutor(t *testing.T) {
	exec := NewShellExecutor()

	if exec.Stdout != os.Stdout {
		t.Error("Expected Stdout to be os.Stdout")
	}

	if exec.Stderr != os.Stderr {
		t.Error("Expected Stderr to be os.Stderr")
	}
}

func TestExecute(t *testing.T) {
	// 保存原始环境变量
	originalShell := os.Getenv("SHELL")
	originalPSModulePath := os.Getenv("PSModulePath")
	defer func() {
		os.Setenv("SHELL", originalShell)
		os.Setenv("PSModulePath", originalPSModulePath)
	}()

	// 创建一个自定义输出的执行器
	var stdout, stderr bytes.Buffer
	exec := &ShellExecutor{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	// 测试简单的echo命令
	testCases := []struct {
		name    string
		command string
		setup   func()
		check   func(t *testing.T, stdout, stderr string, err error)
	}{
		{
			name:    "echo command",
			command: "echo test",
			setup:   func() {},
			check: func(t *testing.T, stdout, stderr string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if !strings.Contains(stdout, "test") {
					t.Errorf("Expected output to contain 'test', got %s", stdout)
				}
			},
		},
		{
			name:    "invalid command",
			command: "invalidcommand",
			setup:   func() {},
			check: func(t *testing.T, stdout, stderr string, err error) {
				if err == nil {
					t.Error("Expected error for invalid command, got none")
				}
			},
		},
	}

	// Windows特定的测试用例
	if runtime.GOOS == "windows" {
		// PowerShell测试
		testCases = append(testCases, struct {
			name    string
			command string
			setup   func()
			check   func(t *testing.T, stdout, stderr string, err error)
		}{
			name:    "powershell command",
			command: "Write-Output 'test'",
			setup: func() {
				os.Setenv("PSModulePath", "some-path")
			},
			check: func(t *testing.T, stdout, stderr string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if !strings.Contains(stdout, "test") {
					t.Errorf("Expected output to contain 'test', got %s", stdout)
				}
			},
		})

		// CMD测试
		testCases = append(testCases, struct {
			name    string
			command string
			setup   func()
			check   func(t *testing.T, stdout, stderr string, err error)
		}{
			name:    "cmd command",
			command: "echo test",
			setup: func() {
				os.Setenv("PSModulePath", "")
			},
			check: func(t *testing.T, stdout, stderr string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if !strings.Contains(stdout, "test") {
					t.Errorf("Expected output to contain 'test', got %s", stdout)
				}
			},
		})
	}

	// Unix系统特定的测试用例
	if runtime.GOOS != "windows" {
		testCases = append(testCases, struct {
			name    string
			command string
			setup   func()
			check   func(t *testing.T, stdout, stderr string, err error)
		}{
			name:    "shell environment test",
			command: "echo $SHELL",
			setup: func() {
				os.Setenv("SHELL", "/bin/bash")
			},
			check: func(t *testing.T, stdout, stderr string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if !strings.Contains(stdout, "/bin/bash") {
					t.Errorf("Expected output to contain '/bin/bash', got %s", stdout)
				}
			},
		})
	}

	// 运行所有测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 重置输出缓冲区
			stdout.Reset()
			stderr.Reset()

			// 设置测试环境
			tc.setup()

			// 执行命令
			err := exec.Execute(tc.command)

			// 检查结果
			tc.check(t, stdout.String(), stderr.String(), err)
		})
	}
}