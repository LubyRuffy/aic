package ollama

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	baseURL := "http://localhost:11434"
	client := NewClient(baseURL)

	if client.BaseURL != baseURL {
		t.Errorf("Expected BaseURL %s, got %s", baseURL, client.BaseURL)
	}
}

func TestGenerate(t *testing.T) {
	// 创建一个测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法和路径
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/generate" {
			t.Errorf("Expected /api/generate path, got %s", r.URL.Path)
		}

		// 验证请求体
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Error decoding request body: %v", err)
		}

		expectedModel := "test-model"
		expectedPrompt := "test prompt"
		if req.Model != expectedModel {
			t.Errorf("Expected model %s, got %s", expectedModel, req.Model)
		}
		if req.Prompt != expectedPrompt {
			t.Errorf("Expected prompt %s, got %s", expectedPrompt, req.Prompt)
		}

		// 返回测试响应
		resp := Response{Response: "test response"}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 使用测试服务器的URL创建客户端
	client := NewClient(server.URL)

	// 测试Generate方法
	response, err := client.Generate("test-model", "test prompt")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedResponse := "test response"
	if response != expectedResponse {
		t.Errorf("Expected response %s, got %s", expectedResponse, response)
	}
}

func TestGenerateError(t *testing.T) {
	testCases := []struct {
		name           string
		handler        func(w http.ResponseWriter, r *http.Request)
		expectedErrStr string
	}{
		{
			name: "internal server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedErrStr: "unexpected status code: 500",
		},
		{
			name: "ollama service error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "model not found"})
			},
			expectedErrStr: "Ollama service error: model not found",
		},
		{
			name: "cannot generate command",
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := Response{Response: "<err_cannot_generate_command>"}
				json.NewEncoder(w).Encode(resp)
			},
			expectedErrStr: "unable to generate command based on your description",
		},
		{
			name: "invalid json response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("invalid json"))
			},
			expectedErrStr: "failed to parse response data",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tc.handler))
			defer server.Close()

			client := NewClient(server.URL)
			_, err := client.Generate("test-model", "test prompt")

			if err == nil {
				t.Error("Expected error, got nil")
			} else if !strings.Contains(err.Error(), tc.expectedErrStr) {
				t.Errorf("Expected error containing '%s', got '%s'", tc.expectedErrStr, err.Error())
			}
		})
	}
}

func TestGenerateNetworkError(t *testing.T) {
	// 测试网络连接错误
	client := NewClient("http://invalid-url")
	_, err := client.Generate("test-model", "test prompt")
	if err == nil {
		t.Error("Expected network error, got nil")
	} else if !strings.Contains(err.Error(), "failed to connect to Ollama service") {
		t.Errorf("Expected network error, got: %v", err)
	}
}
