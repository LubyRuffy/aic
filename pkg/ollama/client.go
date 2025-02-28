package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/LubyRuffy/aic/pkg/sysinfo"
)

// Client 是Ollama API的客户端
// 它提供了与Ollama服务进行交互的方法
type Client struct {
	BaseURL string
	Verbose bool
}

type Options struct {
	Temperature float32 `json:"temperature"`
}

// Request 是发送给Ollama的请求结构
type Request struct {
	Model   string  `json:"model"`
	Prompt  string  `json:"prompt"`
	System  string  `json:"system"`
	Options Options `json:"options"`
	Stream  bool    `json:"stream"`
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
func NewClient(baseURL string, verbose bool) *Client {
	return &Client{BaseURL: baseURL, Verbose: verbose}
}

func genSystemPrompt() (string, error) {
	// 获取系统信息
	sysInfo, err := sysinfo.GetSystemInfo()
	if err != nil {
		return "", fmt.Errorf("failed to get system info: %w", err)
	}

	// 构建环境变量列表
	envKeys := make([]string, 0, len(sysInfo.EnvVars))
	for k := range sysInfo.EnvVars {
		envKeys = append(envKeys, k)
	}

	// 构建系统提示词
	systemPrompt := fmt.Sprintf(`You are a command line assistant, please generate commands that match the current system environment based on user's description.

## Response Format
- Only provide the command in response, no explanation.
- The command MUST be complete and executable.
- If no corresponding command exists, return "<err_cannot_generate_command>".
- NEVER return natural language responses or greetings.
- NEVER return incomplete or invalid shell commands.
- NEVER return "Im sorry"

## Command Examples
Here are some examples of valid and invalid responses:

### Invalid Responses (DO NOT USE):
Input: "hi"
Output: "Hello! How can I assist you today?"
(This is wrong because it's a natural language response, not a command)

Input: "show me files"
Output: "files"
(This is wrong because it's an incomplete command)

Input: "request qq.com with q param equal a+b"
Output: "curl -s 'http://qq.com/?q=a%%2Bb'"
(This is wrong because the URL contains unescaped special characters)

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
4. Properly handles special characters and URLs:
- Always URL-encode special characters in URLs
- Escape spaces with %%20 or quotes
- Use proper quotes for arguments containing spaces
- Escape special shell characters when needed
5. Fully complies with the above environment
6. Avoids using APIs that require API keys whenever possible, if an API key is required, verifies that the necessary environment variables are set first
7. For internet requests, follow these specific rules:
   - For weather queries, use "https://wttr.in/"

## Special Character Handling Examples:
1. URLs with special characters:
Input: "request qq.com with q param equal a+b"
Output: curl -s "http://qq.com/?q=a%%2Bb"

2. Commands with spaces in arguments:
Input: "create folder named 'my documents'"
Output: mkdir "my documents"

3. Commands with special shell characters:
Input: "find files with name containing '&'"
Output: find . -name "*\&*"

## Current System Environment:
- OS: %s %s
- Shell Type: %s
- Username: %s
- Home Directory: %s
- Current Directory: %s
- Environment Variables: %s
`,
		sysInfo.OS, sysInfo.OSVersion,
		sysInfo.Shell,
		sysInfo.Username,
		sysInfo.HomeDir,
		sysInfo.CurrentDir,
		strings.Join(envKeys, ", "))
	return systemPrompt, nil
}

// Generate 发送生成请求到Ollama服务
func (c *Client) Generate(model, prompt string) (string, error) {
	systemPrompt, err := genSystemPrompt()
	if err != nil {
		return "", fmt.Errorf("failed to genSystemPrompt: %w", err)
	}

	if c.Verbose {
		fmt.Println("System Prompt:")
		fmt.Println(systemPrompt)
	}

	reqData := Request{
		Model:  model,
		Prompt: prompt,
		System: systemPrompt,
		Stream: false,
		Options: Options{
			Temperature: 0.95,
		},
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
			return "", fmt.Errorf("ollama service error: %s", errResp.Error)
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
