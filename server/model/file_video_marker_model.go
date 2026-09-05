package model

import (
	"time"

	"gorm.io/gorm"
	"onij/util/cdb"
)

const (
	TableName_FileVideoMarker = "file_video_marker"

	FileVideoMarker_Id        = "id"
	FileVideoMarker_FileId    = "file_id"
	FileVideoMarker_TimeMs    = "time_ms"
	FileVideoMarker_Label     = "label"
	FileVideoMarker_CreatedAt = "created_at"
	FileVideoMarker_UpdatedAt = "updated_at"
	FileVideoMarker_DeletedAt = "deleted_at"
)

type FileVideoMarker struct {
	Id        int64          `gorm:"column:id" json:"id"`
	FileId    int64          `gorm:"column:file_id" json:"file_id"`
	TimeMs    int64          `gorm:"column:time_ms" json:"time_ms"`
	Label     string         `gorm:"column:label" json:"label"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (FileVideoMarker) TableName() string { return TableName_FileVideoMarker }

func (FileVideoMarker) UpdatableColumns() []string {
	return []string{
		FileVideoMarker_FileId,
		FileVideoMarker_TimeMs,
		FileVideoMarker_Label,
		FileVideoMarker_CreatedAt,
		FileVideoMarker_UpdatedAt,
		FileVideoMarker_DeletedAt,
	}
}

type FileVideoMarkerQuerier cdb.BasicQuerier

func (q FileVideoMarkerQuerier) super() cdb.BasicQuerier { return cdb.BasicQuerier(q) }
func (q FileVideoMarkerQuerier) ToOptions() []cdb.Option { return q.super().ToOptions() }
func (q FileVideoMarkerQuerier) WithAsc(cols ...string) FileVideoMarkerQuerier {
	q.super().WithAsc(cols...)
	return q
}
func (q FileVideoMarkerQuerier) WithDesc(cols ...string) FileVideoMarkerQuerier {
	q.super().WithDesc(cols...)
	return q
}
func (q FileVideoMarkerQuerier) Id(v any) FileVideoMarkerQuerier {
	q.super().Add(FileVideoMarker_Id, v)
	return q
}
func (q FileVideoMarkerQuerier) FileId(v any) FileVideoMarkerQuerier {
	q.super().Add(FileVideoMarker_FileId, v)
	return q
}
func (q FileVideoMarkerQuerier) TimeMs(v any) FileVideoMarkerQuerier {
	q.super().Add(FileVideoMarker_TimeMs, v)
	return q
}

type FileVideoMarkerUpdater cdb.BasicUpdater

func (u FileVideoMarkerUpdater) super() cdb.BasicUpdater { return cdb.BasicUpdater(u) }
func (u FileVideoMarkerUpdater) ToMap() map[string]any   { return u.super().ToMap() }
func (u FileVideoMarkerUpdater) IsEmpty() bool           { return u.super().IsEmpty() }
func (u FileVideoMarkerUpdater) TimeMs(v int64) FileVideoMarkerUpdater {
	u.super().Add(FileVideoMarker_TimeMs, v)
	return u
}
func (u FileVideoMarkerUpdater) Label(v string) FileVideoMarkerUpdater {
	u.super().Add(FileVideoMarker_Label, v)
	return u
}
