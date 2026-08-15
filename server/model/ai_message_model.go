package model

import (
	"time"

	"gorm.io/gorm"
	"onij/util/cdb"
)

const (
	TableName_AiMessage = "ai_message"

	AiMessage_Id              = "id"
	AiMessage_SessionId       = "session_id"
	AiMessage_Role            = "role"
	AiMessage_Content         = "content"
	AiMessage_ToolName        = "tool_name"
	AiMessage_ToolCallId      = "tool_call_id"
	AiMessage_TokenPrompt     = "token_prompt"
	AiMessage_TokenCompletion = "token_completion"
	AiMessage_TokenTotal      = "token_total"
	AiMessage_CreatedAt       = "created_at"
	AiMessage_UpdatedAt       = "updated_at"
	AiMessage_DeletedAt       = "deleted_at"

	AiRoleUser      = "user"
	AiRoleAssistant = "assistant"
	AiRoleTool      = "tool"
	AiRoleSystem    = "system"
)

type AiMessage struct {
	Id              int64          `gorm:"column:id" json:"id"`
	SessionId       int64          `gorm:"column:session_id" json:"session_id"`
	Role            string         `gorm:"column:role" json:"role"`
	Content         string         `gorm:"column:content" json:"content"`
	ToolName        string         `gorm:"column:tool_name" json:"tool_name"`
	ToolCallId      string         `gorm:"column:tool_call_id" json:"tool_call_id"`
	TokenPrompt     int32          `gorm:"column:token_prompt" json:"token_prompt"`
	TokenCompletion int32          `gorm:"column:token_completion" json:"token_completion"`
	TokenTotal      int32          `gorm:"column:token_total" json:"token_total"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (AiMessage) TableName() string { return TableName_AiMessage }

func (AiMessage) UpdatableColumns() []string {
	return []string{
		AiMessage_SessionId,
		AiMessage_Role,
		AiMessage_Content,
		AiMessage_ToolName,
		AiMessage_ToolCallId,
		AiMessage_TokenPrompt,
		AiMessage_TokenCompletion,
		AiMessage_TokenTotal,
		AiMessage_CreatedAt,
		AiMessage_UpdatedAt,
		AiMessage_DeletedAt,
	}
}

type AiMessageQuerier cdb.BasicQuerier

func (q AiMessageQuerier) super() cdb.BasicQuerier { return cdb.BasicQuerier(q) }
func (q AiMessageQuerier) ToOptions() []cdb.Option { return q.super().ToOptions() }
func (q AiMessageQuerier) WithAsc(cols ...string) AiMessageQuerier {
	q.super().WithAsc(cols...)
	return q
}
func (q AiMessageQuerier) WithDesc(cols ...string) AiMessageQuerier {
	q.super().WithDesc(cols...)
	return q
}
func (q AiMessageQuerier) Id(v any) AiMessageQuerier {
	q.super().Add(AiMessage_Id, v)
	return q
}
func (q AiMessageQuerier) SessionId(v any) AiMessageQuerier {
	q.super().Add(AiMessage_SessionId, v)
	return q
}

type AiMessageUpdater cdb.BasicUpdater

func (u AiMessageUpdater) super() cdb.BasicUpdater { return cdb.BasicUpdater(u) }
func (u AiMessageUpdater) ToMap() map[string]any   { return u.super().ToMap() }
func (u AiMessageUpdater) IsEmpty() bool           { return u.super().IsEmpty() }
