# aic

[![Go Report Card](https://goreportcard.com/badge/github.com/LubyRuffy/aic)](https://goreportcard.com/report/github.com/LubyRuffy/aic)
[![GitHub release](https://img.shields.io/github/release/LubyRuffy/aic.svg)](https://github.com/LubyRuffy/aic/releases)
[![Build Status](https://github.com/LubyRuffy/aic/workflows/CI/badge.svg)](https://github.com/LubyRuffy/aic/actions)
[![License](https://img.shields.io/github/license/LubyRuffy/aic.svg)](https://github.com/LubyRuffy/aic/blob/main/LICENSE)

AIC (AI Client) 是一个基于 Ollama 的智能命令行工具，它能够将自然语言描述转换为系统命令并执行。通过简单的描述，AIC 可以帮助你快速找到并执行所需的命令，无需记忆复杂的命令语法。

## 特性

- 🤖 基于 Ollama 的智能命令生成
- 🌈 支持多种操作系统（Windows、macOS、Linux）
- 🔧 兼容多种 Shell（bash、zsh、PowerShell、cmd）
- 🎨 美观的彩色输出界面
- 🔍 详细的调试模式
- ⚡ 快速且轻量级

## 安装

### 使用 Go 安装

```bash
go install github.com/LubyRuffy/aic@latest
```

### 从发布页面下载

访问 [GitHub Releases](https://github.com/LubyRuffy/aic/releases) 页面下载适合你系统的预编译二进制文件。

## 前置条件

1. 安装 [Ollama](https://ollama.ai)
2. 拉取所需模型（默认使用 qwen2.5-coder）：
```bash
ollama pull qwen2.5-coder
```

## 使用方法

### 基本用法

```bash
aic "你想执行的操作描述"
```

### 示例

```bash
# 列出当前目录下的所有文件
aic "显示当前目录下的所有文件"

# 查看系统内存使用情况
aic "查看系统内存使用情况"
```

### 命令行参数

```bash
aic [选项] <提示词>

选项：
  -model string
        指定使用的 Ollama 模型 (默认 "qwen2.5-coder")
  -verbose
        启用详细模式，显示生成的实际命令
  -ollama-url string
        指定 Ollama 服务地址 (默认 "http://localhost:11434")
```

## 开发

### 环境要求

- Go 1.21 或更高版本
- golangci-lint（用于代码检查）

### 构建

```bash
# 克隆仓库
git clone https://github.com/LubyRuffy/aic.git
cd aic

# 安装依赖
go mod tidy

# 构建项目
go build
```

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。
