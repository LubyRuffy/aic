package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/LubyRuffy/aic/pkg/sysinfo"
)

// Client 是Ollama API的客户端
// 它提供了与Ollama服务进行交互的方法
type Client struct {
	BaseURL string
}

// Request 是发送给Ollama的请求结构
type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	System string `json:"system"`
	Stream bool   `json:"stream"`
}

// Response 是Ollama的响应结构
type Response struct {
	Response string `json:"response"`
}

// ErrorResponse 是Ollama的错误响应结构
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewClient 创建一个新的Ollama客户端
func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

// Generate 发送生成请求到Ollama服务
func (c *Client) Generate(model, prompt string) (string, error) {
	// 获取系统信息
	sysInfo, err := sysinfo.GetSystemInfo()
	if err != nil {
		return "", fmt.Errorf("failed to get system info: %w", err)
	}

	// 构建系统提示词
	systemPrompt := fmt.Sprintf(`You are a command line assistant, please generate commands that match the current system environment based on user's description.

## Response Format
- Only provide the command in response, no explanation.
- The command MUST be complete and executable.
- If no corresponding command exists, return "<err_cannot_generate_command>".
- NEVER return natural language responses or greetings.
- NEVER return incomplete or invalid shell commands.

## Command Examples
Here are some examples of valid and invalid responses:

### Invalid Responses (DO NOT USE):
Input: "hi"
Output: "Hello! How can I assist you today?"
(This is wrong because it's a natural language response, not a command)

Input: "show me files"
Output: "files"
(This is wrong because it's an incomplete command)

### Valid Commands for Different Environments:

1. For macOS/Linux (bash/zsh):
   Input: "Show disk usage"
   Output: df -h

   Input: "List files in current directory"
   Output: ls -la

2. For Windows (cmd):
   Input: "Show disk usage"
   Output: wmic logicaldisk get size,freespace,caption

   Input: "List files in current directory"
   Output: dir

3. For Windows (PowerShell):
   Input: "Show disk usage"
   Output: Get-PSDrive -PSProvider FileSystem

   Input: "List files in current directory"
   Output: Get-ChildItem

Please ensure the generated command:
1. Is complete and executable
2. Uses the correct syntax for the current shell
3. Includes all necessary flags and parameters
4. Fully complies with the above environment

## Current System Environment:
- OS: %s %s
- Shell Type: %s
- Username: %s
- Home Directory: %s
- Current Directory: %s

`,
		sysInfo.OS, sysInfo.OSVersion,
		sysInfo.Shell,
		sysInfo.Username,
		sysInfo.HomeDir,
		sysInfo.CurrentDir)

	reqData := Request{
		Model:  model,
		Prompt: prompt,
		System: systemPrompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("failed to serialize request data: %w", err)
	}

	resp, err := http.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama service: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
			return "", fmt.Errorf("Ollama service error: %s", errResp.Error)
		}
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response data: %w", err)
	}

	var ollamaResp Response
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to parse response data: %w \n %s", err, body)
	}

	// 检查是否为无法生成命令的错误标记
	if ollamaResp.Response == "<err_cannot_generate_command>" {
		return "", fmt.Errorf("unable to generate command based on your description, please try to be more specific")
	}

	return ollamaResp.Response, nil
}
