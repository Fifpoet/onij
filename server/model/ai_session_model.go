package model

import (
	"time"

	"gorm.io/gorm"
	"onij/util/cdb"
)

const (
	TableName_AiSession = "ai_session"

	AiSession_Id         = "id"
	AiSession_Title      = "title"
	AiSession_Kind       = "kind"
	AiSession_TokenTotal = "token_total"
	AiSession_CreatedAt  = "created_at"
	AiSession_UpdatedAt  = "updated_at"
	AiSession_DeletedAt  = "deleted_at"

	AiSessionKindChat  int32 = 0
	AiSessionKindVoice int32 = 1

	AiSessionTitleVoice = "语音助手"
)

type AiSession struct {
	Id         int64          `gorm:"column:id" json:"id"`
	Title      string         `gorm:"column:title" json:"title"`
	Kind       int32          `gorm:"column:kind" json:"kind"`
	TokenTotal int64          `gorm:"column:token_total" json:"token_total"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (AiSession) TableName() string { return TableName_AiSession }

func (AiSession) UpdatableColumns() []string {
	return []string{
		AiSession_Title,
		AiSession_Kind,
		AiSession_TokenTotal,
		AiSession_CreatedAt,
		AiSession_UpdatedAt,
		AiSession_DeletedAt,
	}
}

type AiSessionQuerier cdb.BasicQuerier

func (q AiSessionQuerier) super() cdb.BasicQuerier { return cdb.BasicQuerier(q) }
func (q AiSessionQuerier) ToOptions() []cdb.Option { return q.super().ToOptions() }
func (q AiSessionQuerier) WithAsc(cols ...string) AiSessionQuerier {
	q.super().WithAsc(cols...)
	return q
}
func (q AiSessionQuerier) WithDesc(cols ...string) AiSessionQuerier {
	q.super().WithDesc(cols...)
	return q
}
func (q AiSessionQuerier) Id(v any) AiSessionQuerier {
	q.super().Add(AiSession_Id, v)
	return q
}
func (q AiSessionQuerier) Kind(v any) AiSessionQuerier {
	q.super().Add(AiSession_Kind, v)
	return q
}

type AiSessionUpdater cdb.BasicUpdater

func (u AiSessionUpdater) super() cdb.BasicUpdater { return cdb.BasicUpdater(u) }
func (u AiSessionUpdater) ToMap() map[string]any   { return u.super().ToMap() }
func (u AiSessionUpdater) IsEmpty() bool           { return u.super().IsEmpty() }
