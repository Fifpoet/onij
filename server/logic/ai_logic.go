package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"onij/biz/prm"
	"onij/infra"
	"onij/model"
	"onij/util"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	aiHistoryMaxMessages    = 40
	aiVoiceHistoryUserTurns = 3
	aiVoiceSessionTitle     = model.AiSessionTitleVoice
)

type AiStreamEmitter func(event string, payload any) error

type AiLogic interface {
	CreateSession(ctx context.Context) (*prm.AiSessionItemResult, error)
	ListSessions(ctx context.Context) (*prm.ListAiSessionResult, error)
	ListMessages(ctx context.Context, sessionId int64) (*prm.ListAiMessageResult, error)
	DeleteSession(ctx context.Context, sessionId int64) (*prm.OkAiResult, error)
	Chat(ctx context.Context, param *prm.AiChatParam) (*prm.AiChatResult, error)
	ChatStream(ctx context.Context, param *prm.AiChatParam, emit AiStreamEmitter) error
	ConfirmStream(ctx context.Context, param *prm.AiConfirmParam, emit AiStreamEmitter) error
}

type aiLogic struct {
	*infra.AllInfra
}

func NewAiLogic(i *infra.AllInfra) AiLogic {
	return &aiLogic{AllInfra: i}
}

func (l *aiLogic) CreateSession(ctx context.Context) (*prm.AiSessionItemResult, error) {
	now := time.Now()
	row := &model.AiSession{
		Id:        util.IdGen.Generate(),
		Title:     "新对话",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := l.AiSessionDal.Save(ctx, row); err != nil {
		return nil, err
	}
	return &prm.AiSessionItemResult{Item: toAiSessionDTO(row)}, nil
}

func (l *aiLogic) ListSessions(ctx context.Context) (*prm.ListAiSessionResult, error) {
	voice, err := l.ensureVoiceSession(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := l.AiSessionDal.ListRecent(ctx, 50)
	if err != nil {
		return nil, err
	}
	items := make([]*prm.AiSessionDTO, 0, len(rows)+1)
	items = append(items, toAiSessionDTO(voice))
	for _, r := range rows {
		if r == nil || r.Id == voice.Id {
			continue
		}
		items = append(items, toAiSessionDTO(r))
	}
	return &prm.ListAiSessionResult{Items: items}, nil
}

func (l *aiLogic) ListMessages(ctx context.Context, sessionId int64) (*prm.ListAiMessageResult, error) {
	sess, err := l.AiSessionDal.GetById(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, fmt.Errorf("session not found")
	}
	rows, err := l.AiMessageDal.ListBySessionId(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	items := make([]*prm.AiMessageDTO, 0, len(rows))
	for _, r := range rows {
		items = append(items, toAiMessageDTO(r))
	}
	return &prm.ListAiMessageResult{Items: items, TokenTotal: sess.TokenTotal}, nil
}

func (l *aiLogic) DeleteSession(ctx context.Context, sessionId int64) (*prm.OkAiResult, error) {
	if sessionId <= 0 {
		return nil, fmt.Errorf("session_id is empty")
	}
	sess, err := l.AiSessionDal.GetById(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, fmt.Errorf("session not found")
	}
	if sess.Kind == model.AiSessionKindVoice {
		return nil, fmt.Errorf("语音助手会话不可删除")
	}
	_ = l.AiMessageDal.DeleteBySessionId(ctx, sessionId)
	if err = l.AiSessionDal.DeleteById(ctx, sessionId); err != nil {
		return nil, err
	}
	return &prm.OkAiResult{}, nil
}

func (l *aiLogic) ensureVoiceSession(ctx context.Context) (*model.AiSession, error) {
	row, err := l.AiSessionDal.GetByKind(ctx, model.AiSessionKindVoice)
	if err != nil {
		return nil, err
	}
	if row != nil {
		if row.Title != aiVoiceSessionTitle {
			row.Title = aiVoiceSessionTitle
			row.UpdatedAt = time.Now()
			_, _ = l.AiSessionDal.Save(ctx, row)
		}
		return row, nil
	}
	now := time.Now()
	row = &model.AiSession{
		Id:        util.IdGen.Generate(),
		Title:     aiVoiceSessionTitle,
		Kind:      model.AiSessionKindVoice,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err = l.AiSessionDal.Save(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (l *aiLogic) ensureSession(ctx context.Context, sessionId int64) (*model.AiSession, error) {
	if sessionId > 0 {
		sess, err := l.AiSessionDal.GetById(ctx, sessionId)
		if err != nil {
			return nil, err
		}
		if sess == nil {
			return nil, fmt.Errorf("session not found")
		}
		return sess, nil
	}
	now := time.Now()
	row := &model.AiSession{
		Id:        util.IdGen.Generate(),
		Title:     "新对话",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := l.AiSessionDal.Save(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (l *aiLogic) Chat(ctx context.Context, param *prm.AiChatParam) (*prm.AiChatResult, error) {
	content := strings.TrimSpace(param.Content)
	if content == "" {
		return nil, fmt.Errorf("content is empty")
	}
	key := util.DeepseekAPIKey()
	if key == "" {
		return nil, fmt.Errorf("dskey not configured")
	}

	sess, err := l.ensureSession(ctx, param.SessionId)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	userMsg := &model.AiMessage{
		Id:        util.IdGen.Generate(),
		SessionId: sess.Id,
		Role:      model.AiRoleUser,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err = l.AiMessageDal.Save(ctx, userMsg); err != nil {
		return nil, err
	}
	if sess.Kind != model.AiSessionKindVoice && (sess.Title == "新对话" || sess.Title == "") {
		sess.Title = truncateRunes(content, 24)
	}

	final, err := l.runAgentLoop(ctx, key, sess, nil)
	if err != nil {
		return nil, err
	}
	if final.Paused != nil {
		return &prm.AiChatResult{
			Session: toAiSessionDTO(sess),
			Reply:   "",
		}, fmt.Errorf("needs_confirm:%d", final.Paused.ConfirmId)
	}
	reply := ""
	if final.Message != nil {
		reply = final.Message.Content
	}
	return &prm.AiChatResult{
		Session:  toAiSessionDTO(sess),
		Messages: []*prm.AiMessageDTO{toAiMessageDTO(userMsg), toAiMessageDTO(final.Message)},
		Reply:    reply,
	}, nil
}

func (l *aiLogic) ChatStream(ctx context.Context, param *prm.AiChatParam, emit AiStreamEmitter) error {
	if emit == nil {
		return fmt.Errorf("emitter is nil")
	}
	content := strings.TrimSpace(param.Content)
	if content == "" {
		return fmt.Errorf("content is empty")
	}
	key := util.DeepseekAPIKey()
	if key == "" {
		return fmt.Errorf("dskey not configured")
	}

	sess, err := l.ensureSession(ctx, param.SessionId)
	if err != nil {
		return err
	}
	if err = emit("session", toAiSessionDTO(sess)); err != nil {
		return err
	}

	now := time.Now()
	userMsg := &model.AiMessage{
		Id:        util.IdGen.Generate(),
		SessionId: sess.Id,
		Role:      model.AiRoleUser,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err = l.AiMessageDal.Save(ctx, userMsg); err != nil {
		return err
	}
	if sess.Kind != model.AiSessionKindVoice && (sess.Title == "新对话" || sess.Title == "") {
		sess.Title = truncateRunes(content, 24)
	}
	if err = emit("user", toAiMessageDTO(userMsg)); err != nil {
		return err
	}

	final, err := l.runAgentLoop(ctx, key, sess, emit)
	if err != nil {
		return err
	}
	if final.Paused != nil {
		_ = emit("needs_confirm", final.Paused)
		_ = emit("paused", map[string]any{
			"session":    toAiSessionDTO(sess),
			"confirm_id": final.Paused.ConfirmId,
		})
		return nil
	}
	return emit("done", map[string]any{
		"session": toAiSessionDTO(sess),
		"message": toAiMessageDTO(final.Message),
		"reply":   final.Message.Content,
	})
}

type agentLoopOutcome struct {
	Message *model.AiMessage
	Paused  *prm.AiConfirmDTO
}

func (l *aiLogic) runAgentLoop(ctx context.Context, key string, sess *model.AiSession, emit AiStreamEmitter) (*agentLoopOutcome, error) {
	tools := aiToolDefs()
	maxSteps := aiAgentMaxSteps()
	var lastAssistant *model.AiMessage

	for step := 0; step < maxSteps; step++ {
		history, err := l.AiMessageDal.ListBySessionId(ctx, sess.Id)
		if err != nil {
			return nil, err
		}
		llmMessages := buildDeepseekHistory(history, sess)

		dsResp, err := util.DeepseekChatWithTools(ctx, key, llmMessages, tools, 0.3)
		if err != nil {
			return nil, err
		}

		now := time.Now()
		if len(dsResp.ToolCalls) > 0 {
			assistantMsg := &model.AiMessage{
				Id:              util.IdGen.Generate(),
				SessionId:       sess.Id,
				Role:            model.AiRoleAssistant,
				Content:         encodeAssistantToolCallsContent(dsResp.Content, dsResp.ToolCalls),
				ToolName:        aiToolCallsMarker,
				TokenPrompt:     int32(dsResp.Usage.PromptTokens),
				TokenCompletion: int32(dsResp.Usage.CompletionTokens),
				TokenTotal:      int32(dsResp.Usage.TotalTokens),
				CreatedAt:       now,
				UpdatedAt:       now,
			}
			if _, err = l.AiMessageDal.Save(ctx, assistantMsg); err != nil {
				return nil, err
			}
			sess.TokenTotal += int64(dsResp.Usage.TotalTokens)
			lastAssistant = assistantMsg
			sess.UpdatedAt = time.Now()
			_, _ = l.AiSessionDal.Save(ctx, sess)

			if emit != nil {
				_ = emit("assistant_tools", map[string]any{
					"content":    dsResp.Content,
					"tool_calls": dsResp.ToolCalls,
					"step":       step + 1,
				})
			}

			paused, err := l.executeToolBatch(ctx, sess, dsResp.ToolCalls, emit)
			if err != nil {
				return nil, err
			}
			if paused != nil {
				return &agentLoopOutcome{Paused: paused}, nil
			}
			continue
		}

		reply := dsResp.Content
		if reply == "" {
			reply = "（模型未返回内容）"
		}
		assistantMsg := &model.AiMessage{
			Id:              util.IdGen.Generate(),
			SessionId:       sess.Id,
			Role:            model.AiRoleAssistant,
			Content:         reply,
			TokenPrompt:     int32(dsResp.Usage.PromptTokens),
			TokenCompletion: int32(dsResp.Usage.CompletionTokens),
			TokenTotal:      int32(dsResp.Usage.TotalTokens),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		if _, err = l.AiMessageDal.Save(ctx, assistantMsg); err != nil {
			return nil, err
		}
		sess.TokenTotal += int64(dsResp.Usage.TotalTokens)
		sess.UpdatedAt = time.Now()
		if _, err = l.AiSessionDal.Save(ctx, sess); err != nil {
			return nil, err
		}
		if emit != nil {
			_ = emit("delta", map[string]string{"content": reply})
		}
		return &agentLoopOutcome{Message: assistantMsg}, nil
	}

	stop := &model.AiMessage{
		Id:        util.IdGen.Generate(),
		SessionId: sess.Id,
		Role:      model.AiRoleAssistant,
		Content:   "步骤过多已停止，请缩小问题或换一种说法再试。",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if _, err := l.AiMessageDal.Save(ctx, stop); err != nil {
		return nil, err
	}
	sess.UpdatedAt = time.Now()
	_, _ = l.AiSessionDal.Save(ctx, sess)
	if emit != nil {
		_ = emit("delta", map[string]string{"content": stop.Content})
	}
	_ = lastAssistant
	return &agentLoopOutcome{Message: stop}, nil
}

func (l *aiLogic) executeToolBatch(
	ctx context.Context,
	sess *model.AiSession,
	calls []util.DeepseekToolCall,
	emit AiStreamEmitter,
) (*prm.AiConfirmDTO, error) {
	pausedAt := -1
	for i, tc := range calls {
		name := tc.Function.Name
		if emit != nil {
			_ = emit("tool_call", map[string]any{
				"id":        tc.ID,
				"name":      name,
				"arguments": tc.Function.Arguments,
			})
		}

		if toolRequiresConfirm(name) {
			args := map[string]any{}
			if strings.TrimSpace(tc.Function.Arguments) != "" {
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
					exec := failValidation("invalid json arguments")
					if err := l.saveToolResult(ctx, sess, tc, name, exec, emit); err != nil {
						return nil, err
					}
					continue
				}
			}
			confirm, err := l.createConfirmPending(ctx, sess, tc.ID, name, tc.Function.Arguments, args)
			if err != nil {
				exec := failBiz(err.Error())
				if err := l.saveToolResult(ctx, sess, tc, name, exec, emit); err != nil {
					return nil, err
				}
				continue
			}
			pausedAt = i
			// 同批后续 tool 写 skipped，保证协议闭合
			for j := i + 1; j < len(calls); j++ {
				rest := calls[j]
				skipped := aiToolExecResult{
					OK:      false,
					Error:   "skipped",
					Payload: map[string]any{"detail": "paused for human confirm on earlier tool"},
				}
				if emit != nil {
					_ = emit("tool_call", map[string]any{
						"id":        rest.ID,
						"name":      rest.Function.Name,
						"arguments": rest.Function.Arguments,
					})
				}
				if err := l.saveToolResult(ctx, sess, rest, rest.Function.Name, skipped, emit); err != nil {
					return nil, err
				}
			}
			return confirm, nil
		}

		exec := l.executeAiTool(ctx, name, tc.Function.Arguments)
		if exec.ClientAction != nil && emit != nil {
			_ = emit("client_action", exec.ClientAction)
		}
		if err := l.saveToolResult(ctx, sess, tc, name, exec, emit); err != nil {
			return nil, err
		}
	}
	_ = pausedAt
	return nil, nil
}

func (l *aiLogic) saveToolResult(
	ctx context.Context,
	sess *model.AiSession,
	tc util.DeepseekToolCall,
	name string,
	exec aiToolExecResult,
	emit AiStreamEmitter,
) error {
	resultContent := toolResultJSON(exec)
	toolMsg := &model.AiMessage{
		Id:         util.IdGen.Generate(),
		SessionId:  sess.Id,
		Role:       model.AiRoleTool,
		Content:    resultContent,
		ToolName:   name,
		ToolCallId: tc.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if _, err := l.AiMessageDal.Save(ctx, toolMsg); err != nil {
		return err
	}
	if emit != nil {
		_ = emit("tool_result", map[string]any{
			"id":      tc.ID,
			"name":    name,
			"ok":      exec.OK,
			"content": resultContent,
		})
	}
	return nil
}

func toAiSessionDTO(s *model.AiSession) *prm.AiSessionDTO {
	if s == nil {
		return nil
	}
	return &prm.AiSessionDTO{
		Id:         s.Id,
		Title:      s.Title,
		Kind:       s.Kind,
		TokenTotal: s.TokenTotal,
		UpdatedAt:  s.UpdatedAt.UnixMilli(),
		CreatedAt:  s.CreatedAt.UnixMilli(),
	}
}

func toAiMessageDTO(m *model.AiMessage) *prm.AiMessageDTO {
	if m == nil {
		return nil
	}
	content := m.Content
	if m.ToolName == aiToolCallsMarker {
		if text, _, ok := decodeAssistantToolCallsContent(m.Content); ok {
			if text == "" {
				content = "（调用工具中）"
			} else {
				content = text
			}
		}
	}
	return &prm.AiMessageDTO{
		Id:              m.Id,
		SessionId:       m.SessionId,
		Role:            m.Role,
		Content:         content,
		ToolName:        m.ToolName,
		ToolCallId:      m.ToolCallId,
		TokenPrompt:     m.TokenPrompt,
		TokenCompletion: m.TokenCompletion,
		TokenTotal:      m.TokenTotal,
		CreatedAt:       m.CreatedAt.UnixMilli(),
	}
}

func truncateRunes(s string, n int) string {
	if n <= 0 || utf8.RuneCountInString(s) <= n {
		return s
	}
	runes := []rune(s)
	return string(runes[:n]) + "…"
}
