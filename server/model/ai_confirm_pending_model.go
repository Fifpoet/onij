package model

import (
	"time"

	"gorm.io/gorm"
	"onij/util/cdb"
)

const (
	TableName_AiConfirmPending = "ai_confirm_pending"

	AiConfirmPending_Id         = "id"
	AiConfirmPending_SessionId  = "session_id"
	AiConfirmPending_ToolCallId = "tool_call_id"
	AiConfirmPending_ToolName   = "tool_name"
	AiConfirmPending_DraftArgs  = "draft_args"
	AiConfirmPending_UiJson     = "ui_json"
	AiConfirmPending_Status     = "status"
	AiConfirmPending_ExpiresAt  = "expires_at"
	AiConfirmPending_CreatedAt  = "created_at"
	AiConfirmPending_UpdatedAt  = "updated_at"
	AiConfirmPending_DeletedAt  = "deleted_at"

	AiConfirmStatusPending   = "pending"
	AiConfirmStatusConfirmed = "confirmed"
	AiConfirmStatusCancelled = "cancelled"
	AiConfirmStatusExpired   = "expired"
)

type AiConfirmPending struct {
	Id         int64          `gorm:"column:id" json:"id"`
	SessionId  int64          `gorm:"column:session_id" json:"session_id"`
	ToolCallId string         `gorm:"column:tool_call_id" json:"tool_call_id"`
	ToolName   string         `gorm:"column:tool_name" json:"tool_name"`
	DraftArgs  string         `gorm:"column:draft_args" json:"draft_args"`
	UiJson     string         `gorm:"column:ui_json" json:"ui_json"`
	Status     string         `gorm:"column:status" json:"status"`
	ExpiresAt  time.Time      `gorm:"column:expires_at" json:"expires_at"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (AiConfirmPending) TableName() string { return TableName_AiConfirmPending }

func (AiConfirmPending) UpdatableColumns() []string {
	return []string{
		AiConfirmPending_SessionId,
		AiConfirmPending_ToolCallId,
		AiConfirmPending_ToolName,
		AiConfirmPending_DraftArgs,
		AiConfirmPending_UiJson,
		AiConfirmPending_Status,
		AiConfirmPending_ExpiresAt,
		AiConfirmPending_CreatedAt,
		AiConfirmPending_UpdatedAt,
		AiConfirmPending_DeletedAt,
	}
}

type AiConfirmPendingQuerier cdb.BasicQuerier

func (q AiConfirmPendingQuerier) super() cdb.BasicQuerier { return cdb.BasicQuerier(q) }
func (q AiConfirmPendingQuerier) ToOptions() []cdb.Option { return q.super().ToOptions() }
func (q AiConfirmPendingQuerier) WithAsc(cols ...string) AiConfirmPendingQuerier {
	q.super().WithAsc(cols...)
	return q
}
func (q AiConfirmPendingQuerier) WithDesc(cols ...string) AiConfirmPendingQuerier {
	q.super().WithDesc(cols...)
	return q
}
func (q AiConfirmPendingQuerier) Id(v any) AiConfirmPendingQuerier {
	q.super().Add(AiConfirmPending_Id, v)
	return q
}
func (q AiConfirmPendingQuerier) SessionId(v any) AiConfirmPendingQuerier {
	q.super().Add(AiConfirmPending_SessionId, v)
	return q
}
func (q AiConfirmPendingQuerier) Status(v any) AiConfirmPendingQuerier {
	q.super().Add(AiConfirmPending_Status, v)
	return q
}

type AiConfirmPendingUpdater cdb.BasicUpdater

func (u AiConfirmPendingUpdater) super() cdb.BasicUpdater { return cdb.BasicUpdater(u) }
func (u AiConfirmPendingUpdater) ToMap() map[string]any   { return u.super().ToMap() }
func (u AiConfirmPendingUpdater) IsEmpty() bool           { return u.super().IsEmpty() }
func (u AiConfirmPendingUpdater) Status(v any) AiConfirmPendingUpdater {
	u.super().Add(AiConfirmPending_Status, v)
	return u
}
func (u AiConfirmPendingUpdater) UpdatedAt(v any) AiConfirmPendingUpdater {
	u.super().Add(AiConfirmPending_UpdatedAt, v)
	return u
}
