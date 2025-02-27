package main

import (
	"flag"
	"os"
	"strings"

	"github.com/LubyRuffy/aic/pkg/color"
	"github.com/LubyRuffy/aic/pkg/executor"
	"github.com/LubyRuffy/aic/pkg/ollama"
)

func main() {
	// 解析命令行参数
	model := flag.String("model", "qwen2.5-coder", "Ollama model name to use")
	verbose := flag.Bool("verbose", false, "Enable verbose mode to print actual commands")
	ollamaURL := flag.String("ollama-url", "http://localhost:11434", "Ollama service address")
	flag.Parse()

	// 获取prompt
	args := flag.Args()
	if len(args) == 0 {
		color.Warning("Usage: aic [--model model_name] [--verbose] [--ollama-url ollama_address] <prompt>\n")
		os.Exit(1)
	}
	prompt := strings.Join(args, " ")

	// 创建Ollama客户端
	client := ollama.NewClient(*ollamaURL)

	// 生成命令
	response, err := client.Generate(*model, prompt)
	if err != nil {
		color.Error("Error generating command: %v\n", err)
		os.Exit(1)
	}

	// 在verbose模式下打印实际命令
	if *verbose {
		color.Info("Generated command: %s\n", response)
	}

	// 创建命令执行器
	exec := executor.NewShellExecutor()

	// 执行命令
	if err := exec.Execute(response); err != nil {
		color.Error("Error executing command: %v\n", err)
		os.Exit(1)
	}
}