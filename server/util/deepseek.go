package util

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func DeepseekAPIKey() string {
	if v := strings.TrimSpace(os.Getenv("dskey")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("DSKEY"))
}

type DeepseekToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type DeepseekToolCall struct {
	ID       string                   `json:"id"`
	Type     string                   `json:"type"`
	Function DeepseekToolCallFunction `json:"function"`
}

type DeepseekChatMessage struct {
	Role       string             `json:"role"`
	Content    string             `json:"content,omitempty"`
	ToolCalls  []DeepseekToolCall `json:"tool_calls,omitempty"`
	ToolCallID string             `json:"tool_call_id,omitempty"`
	Name       string             `json:"name,omitempty"`
}

type DeepseekToolDef struct {
	Type     string         `json:"type"`
	Function map[string]any `json:"function"`
}

type DeepseekUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type DeepseekChatResult struct {
	Content      string
	ToolCalls    []DeepseekToolCall
	FinishReason string
	Usage        DeepseekUsage
}

func deepseekPost(ctx context.Context, apiKey string, payload map[string]any, timeout time.Duration) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.deepseek.com/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	cli := &http.Client{Timeout: timeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("deepseek http %d: %s", resp.StatusCode, string(raw))
	}
	return raw, nil
}

// DeepseekChat 简易非流式对话（无 tools）
func DeepseekChat(ctx context.Context, apiKey string, messages []DeepseekChatMessage, temperature float64) (*DeepseekChatResult, error) {
	return DeepseekChatWithTools(ctx, apiKey, messages, nil, temperature)
}

// DeepseekChatWithTools 非流式对话，可携带 tools；返回 content 与/或 tool_calls。
func DeepseekChatWithTools(
	ctx context.Context,
	apiKey string,
	messages []DeepseekChatMessage,
	tools []DeepseekToolDef,
	temperature float64,
) (*DeepseekChatResult, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("dskey not configured")
	}
	if temperature <= 0 {
		temperature = 0.3
	}
	payload := map[string]any{
		"model":       "deepseek-chat",
		"messages":    messages,
		"temperature": temperature,
		"stream":      false,
	}
	if len(tools) > 0 {
		payload["tools"] = tools
	}
	raw, err := deepseekPost(ctx, apiKey, payload, 90*time.Second)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Choices []struct {
			FinishReason string `json:"finish_reason"`
			Message      struct {
				Content   string             `json:"content"`
				ToolCalls []DeepseekToolCall `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
		Usage DeepseekUsage `json:"usage"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	out := &DeepseekChatResult{Usage: parsed.Usage}
	if len(parsed.Choices) > 0 {
		ch := parsed.Choices[0]
		out.Content = strings.TrimSpace(ch.Message.Content)
		out.ToolCalls = ch.Message.ToolCalls
		out.FinishReason = ch.FinishReason
	}
	return out, nil
}

// DeepseekChatStream 流式对话；onDelta 每次收到增量文本时回调。不携带 tools。
func DeepseekChatStream(
	ctx context.Context,
	apiKey string,
	messages []DeepseekChatMessage,
	temperature float64,
	onDelta func(delta string) error,
) (*DeepseekChatResult, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("dskey not configured")
	}
	if temperature <= 0 {
		temperature = 0.3
	}
	payload := map[string]any{
		"model":       "deepseek-chat",
		"messages":    messages,
		"temperature": temperature,
		"stream":      true,
		"stream_options": map[string]any{
			"include_usage": true,
		},
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.deepseek.com/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "text/event-stream")

	cli := &http.Client{Timeout: 5 * time.Minute}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("deepseek http %d: %s", resp.StatusCode, string(raw))
	}

	var (
		content strings.Builder
		usage   DeepseekUsage
	)
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			break
		}
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
			Usage *DeepseekUsage `json:"usage"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if chunk.Usage != nil {
			usage = *chunk.Usage
		}
		delta := ""
		if len(chunk.Choices) > 0 {
			delta = chunk.Choices[0].Delta.Content
		}
		if delta == "" {
			continue
		}
		content.WriteString(delta)
		if onDelta != nil {
			if err := onDelta(delta); err != nil {
				return nil, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &DeepseekChatResult{Content: content.String(), Usage: usage}, nil
}
