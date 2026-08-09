package model

import (
	"time"

	"gorm.io/gorm"
	"onij/util/cdb"
)

const (
	TableName_TranItem = "tran_item"

	TranItem_Id        = "id"
	TranItem_Kind      = "kind"
	TranItem_Content   = "content"
	TranItem_FileId    = "file_id"
	TranItem_Name      = "name"
	TranItem_Size      = "size"
	TranItem_Format    = "format"
	TranItem_Pinned    = "pinned"
	TranItem_PinnedAt  = "pinned_at"
	TranItem_CreatedAt = "created_at"
	TranItem_UpdatedAt = "updated_at"
	TranItem_DeletedAt = "deleted_at"

	TranKindText = int32(1)
	TranKindFile = int32(2)
)

type TranItem struct {
	Id        int64          `gorm:"column:id" json:"id"`
	Kind      int32          `gorm:"column:kind" json:"kind"`
	Content   string         `gorm:"column:content" json:"content"`
	FileId    int64          `gorm:"column:file_id" json:"file_id"`
	Name      string         `gorm:"column:name" json:"name"`
	Size      int64          `gorm:"column:size" json:"size"`
	Format    int32          `gorm:"column:format" json:"format"`
	Pinned    bool           `gorm:"column:pinned" json:"pinned"`
	PinnedAt  int64          `gorm:"column:pinned_at" json:"pinned_at"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (TranItem) TableName() string { return TableName_TranItem }

func (TranItem) UpdatableColumns() []string {
	return []string{
		TranItem_Kind,
		TranItem_Content,
		TranItem_FileId,
		TranItem_Name,
		TranItem_Size,
		TranItem_Format,
		TranItem_Pinned,
		TranItem_PinnedAt,
		TranItem_CreatedAt,
		TranItem_UpdatedAt,
		TranItem_DeletedAt,
	}
}

type TranItemQuerier cdb.BasicQuerier

func (q TranItemQuerier) super() cdb.BasicQuerier { return cdb.BasicQuerier(q) }
func (q TranItemQuerier) ToOptions() []cdb.Option { return q.super().ToOptions() }
func (q TranItemQuerier) WithAsc(cols ...string) TranItemQuerier {
	q.super().WithAsc(cols...)
	return q
}
func (q TranItemQuerier) WithDesc(cols ...string) TranItemQuerier {
	q.super().WithDesc(cols...)
	return q
}
func (q TranItemQuerier) Id(v any) TranItemQuerier {
	q.super().Add(TranItem_Id, v)
	return q
}
func (q TranItemQuerier) Kind(v any) TranItemQuerier {
	q.super().Add(TranItem_Kind, v)
	return q
}
func (q TranItemQuerier) FileId(v any) TranItemQuerier {
	q.super().Add(TranItem_FileId, v)
	return q
}

type TranItemUpdater cdb.BasicUpdater

func (u TranItemUpdater) super() cdb.BasicUpdater { return cdb.BasicUpdater(u) }
func (u TranItemUpdater) ToMap() map[string]any   { return u.super().ToMap() }
func (u TranItemUpdater) IsEmpty() bool           { return u.super().IsEmpty() }
func (u TranItemUpdater) Pinned(v bool) TranItemUpdater {
	u.super().Add(TranItem_Pinned, v)
	return u
}
func (u TranItemUpdater) PinnedAt(v int64) TranItemUpdater {
	u.super().Add(TranItem_PinnedAt, v)
	return u
}
