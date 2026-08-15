package prm

import "onij/util"

type AiSessionDTO struct {
	Id         int64  `json:"id,string"`
	Title      string `json:"title"`
	Kind       int32  `json:"kind"`
	TokenTotal int64  `json:"token_total"`
	UpdatedAt  int64  `json:"updated_at"`
	CreatedAt  int64  `json:"created_at"`
}

type AiMessageDTO struct {
	Id              int64  `json:"id,string"`
	SessionId       int64  `json:"session_id,string"`
	Role            string `json:"role"`
	Content         string `json:"content"`
	ToolName        string `json:"tool_name,omitempty"`
	ToolCallId      string `json:"tool_call_id,omitempty"`
	TokenPrompt     int32  `json:"token_prompt"`
	TokenCompletion int32  `json:"token_completion"`
	TokenTotal      int32  `json:"token_total"`
	CreatedAt       int64  `json:"created_at"`
}

type ListAiSessionResult struct {
	Items []*AiSessionDTO
}

func (r *ListAiSessionResult) Resp() map[string]any {
	return map[string]any{"code": util.BaseCodeOK, "message": util.BaseMsgOK, "items": r.Items}
}

type AiSessionItemResult struct {
	Item *AiSessionDTO
}

func (r *AiSessionItemResult) Resp() map[string]any {
	return map[string]any{"code": util.BaseCodeOK, "message": util.BaseMsgOK, "item": r.Item}
}

type ListAiMessageResult struct {
	Items      []*AiMessageDTO
	TokenTotal int64
}

func (r *ListAiMessageResult) Resp() map[string]any {
	return map[string]any{
		"code":        util.BaseCodeOK,
		"message":     util.BaseMsgOK,
		"items":       r.Items,
		"token_total": r.TokenTotal,
	}
}

type AiChatParam struct {
	SessionId int64
	Content   string
}

type AiChatResult struct {
	Session  *AiSessionDTO
	Messages []*AiMessageDTO
	Reply    string
}

func (r *AiChatResult) Resp() map[string]any {
	return map[string]any{
		"code":     util.BaseCodeOK,
		"message":  util.BaseMsgOK,
		"session":  r.Session,
		"messages": r.Messages,
		"reply":    r.Reply,
	}
}

type OkAiResult struct{}

func (r *OkAiResult) Resp() map[string]any {
	return map[string]any{"code": util.BaseCodeOK, "message": util.BaseMsgOK}
}

type AiConfirmOption struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Meta  any    `json:"meta,omitempty"`
}

type AiConfirmUI struct {
	Type     string             `json:"type"` // yes_no | single_select | multi_select
	Title    string             `json:"title"`
	YesLabel string             `json:"yes_label,omitempty"`
	NoLabel  string             `json:"no_label,omitempty"`
	Options  []AiConfirmOption  `json:"options,omitempty"`
	Min      int                `json:"min,omitempty"`
	Max      int                `json:"max,omitempty"`
}

type AiConfirmDTO struct {
	ConfirmId  int64        `json:"confirm_id,string"`
	SessionId  int64        `json:"session_id,string"`
	ToolCallId string       `json:"tool_call_id"`
	ToolName   string       `json:"tool_name"`
	DraftArgs  any          `json:"draft_args"`
	UI         AiConfirmUI  `json:"ui"`
	ExpiresAt  int64        `json:"expires_at"`
}

type AiConfirmDecision struct {
	Type      string   `json:"type"` // yes_no | single_select | multi_select | cancel
	Value     *bool    `json:"value,omitempty"`
	OptionId  string   `json:"option_id,omitempty"`
	OptionIds []string `json:"option_ids,omitempty"`
}

type AiConfirmParam struct {
	ConfirmId int64
	Decision  AiConfirmDecision
}
