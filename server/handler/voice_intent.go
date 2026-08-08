package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type voiceIntentReq struct {
	Text string `json:"text"`
}

type voiceAction struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args,omitempty"`
}

type voiceIntentResp struct {
	Reply   string        `json:"reply"`
	Actions []voiceAction `json:"actions"`
}

func deepseekAPIKey() string {
	if v := strings.TrimSpace(os.Getenv("dskey")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("DSKEY"))
}

// VoiceIntent 仅用于「搜索」口令：从 ASR 文本抽出可搜索的歌名/歌手关键词
// @router /voice/intent [POST]
func VoiceIntent(ctx context.Context, c *app.RequestContext) {
	var req voiceIntentReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	text := strings.TrimSpace(req.Text)
	if text == "" {
		c.String(consts.StatusBadRequest, "text is empty")
		return
	}
	key := deepseekAPIKey()
	if key == "" {
		c.String(consts.StatusInternalServerError, "dskey not configured")
		return
	}

	resp, err := callDeepSeekSearchParse(ctx, key, text)
	if err != nil {
		c.String(consts.StatusBadGateway, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp)
}

func callDeepSeekSearchParse(ctx context.Context, apiKey, text string) (*voiceIntentResp, error) {
	tools := []map[string]any{
		{
			"type": "function",
			"function": map[string]any{
				"name":        "search_music",
				"description": "根据用户语音搜索意图，返回用于网易云搜索的结构化字段",
				"parameters": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "最终搜索关键词（可含歌手+歌名，空格分隔）",
						},
						"artist": map[string]any{
							"type":        "string",
							"description": "歌手名；没有则空字符串",
						},
						"title": map[string]any{
							"type":        "string",
							"description": "歌名；没有则空字符串",
						},
					},
					"required": []string{"query"},
				},
			},
		},
	}

	system := "解析音乐搜索，必须调用 search_music。" +
		"输入常为语音识别结果，可轻度纠正明显同音错字（如凌俊杰→林俊杰）；不确定勿改，勿臆造。" +
		"去掉「搜索/的歌」等指令词；query=检索关键词（优先「歌手 歌名」）；artist/title 能拆则拆否则空。"

	payload := map[string]any{
		"model": "deepseek-chat",
		"messages": []map[string]any{
			{"role": "system", "content": system},
			{"role": "user", "content": text},
		},
		"tools":       tools,
		"tool_choice": map[string]any{"type": "function", "function": map[string]any{"name": "search_music"}},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.deepseek.com/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 45 * time.Second}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(httpResp.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, fmt.Errorf("deepseek http %d: %s", httpResp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content   string `json:"content"`
				ToolCalls []struct {
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	if len(parsed.Choices) == 0 {
		return nil, fmt.Errorf("deepseek empty choices")
	}

	msg := parsed.Choices[0].Message
	out := &voiceIntentResp{
		Reply:   strings.TrimSpace(msg.Content),
		Actions: make([]voiceAction, 0, 1),
	}
	for _, tc := range msg.ToolCalls {
		if tc.Function.Name != "search_music" {
			continue
		}
		args := map[string]any{}
		if strings.TrimSpace(tc.Function.Arguments) != "" {
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
		}
		query, _ := args["query"].(string)
		query = strings.TrimSpace(query)
		if query == "" {
			artist, _ := args["artist"].(string)
			title, _ := args["title"].(string)
			query = strings.TrimSpace(strings.TrimSpace(artist) + " " + strings.TrimSpace(title))
			args["query"] = query
		}
		if query == "" {
			continue
		}
		out.Actions = append(out.Actions, voiceAction{
			Name: "search_music",
			Args: args,
		})
	}
	if len(out.Actions) == 0 {
		return nil, fmt.Errorf("deepseek did not return search_music")
	}
	return out, nil
}
