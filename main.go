package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/LubyRuffy/aic/pkg/color"
	"github.com/LubyRuffy/aic/pkg/executor"
	"github.com/LubyRuffy/aic/pkg/ollama"
)

// Version information, will be injected during build via ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Parse command line arguments
	model := flag.String("model", "qwen2.5-coder", "Ollama model name to use")
	verbose := flag.Bool("verbose", false, "Enable verbose mode to print actual commands")
	ollamaURL := flag.String("ollama-url", "http://localhost:11434", "Ollama service address")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Display version information
	if *showVersion {
		fmt.Printf("aic version %s\nbuilt on %s\ncommit hash %s\n", version, date, commit)
		os.Exit(0)
	}

	// Get prompt
	args := flag.Args()
	if len(args) == 0 {
		color.Warning("Usage: aic [--model model_name] [--verbose] [--ollama-url ollama_address] [--version] <prompt>\n")
		os.Exit(1)
	}
	prompt := strings.Join(args, " ")

	// Print debug information in verbose mode
	if *verbose {
		color.Info("Version: %s (built on %s, commit %s)\n", version, date, commit)
		color.Info("Ollama URL: %s\n", *ollamaURL)
		color.Info("Model: %s\n", *model)
		color.Info("Prompt: %s\n", prompt)
	}

	// Create Ollama client
	client := ollama.NewClient(*ollamaURL, *verbose)

	// Generate command
	response, err := client.Generate(*model, prompt)
	if err != nil {
		color.Error("Error generating command: %v\n", err)
		os.Exit(1)
	}

	// Print actual command in verbose mode
	if *verbose {
		color.Info("Generated command: %s\n", response)
	}

	// Create command executor
	exec := executor.NewShellExecutor()

	// Execute command
	if err := exec.Execute(response); err != nil {
		color.Error("Error executing command: %v\n", err)
		os.Exit(1)
	}
}
