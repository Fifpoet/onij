package model

import (
	"time"

	"gorm.io/gorm"
	"onij/util/cdb"
)

const (
	TableName_MusicCollection = "music_collection"

	MusicCollection_Id         = "id"
	MusicCollection_Name       = "name"
	MusicCollection_CoverUrl   = "cover_url"
	MusicCollection_SearchText = "search_text"
	MusicCollection_CreatedAt  = "created_at"
	MusicCollection_UpdatedAt  = "updated_at"
	MusicCollection_DeletedAt  = "deleted_at"
)

type MusicCollection struct {
	Id         int64          `gorm:"column:id" json:"id"`
	Name       string         `gorm:"column:name" json:"name"`
	CoverUrl   string         `gorm:"column:cover_url" json:"cover_url"`
	SearchText string         `gorm:"column:search_text" json:"search_text"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MusicCollection) TableName() string { return TableName_MusicCollection }

func (MusicCollection) UpdatableColumns() []string {
	return []string{
		MusicCollection_Name,
		MusicCollection_CoverUrl,
		MusicCollection_SearchText,
		MusicCollection_CreatedAt,
		MusicCollection_UpdatedAt,
		MusicCollection_DeletedAt,
	}
}

type MusicCollectionQuerier cdb.BasicQuerier

func (q MusicCollectionQuerier) super() cdb.BasicQuerier { return cdb.BasicQuerier(q) }
func (q MusicCollectionQuerier) ToOptions() []cdb.Option { return q.super().ToOptions() }
func (q MusicCollectionQuerier) WithAsc(cols ...string) MusicCollectionQuerier {
	q.super().WithAsc(cols...)
	return q
}
func (q MusicCollectionQuerier) WithDesc(cols ...string) MusicCollectionQuerier {
	q.super().WithDesc(cols...)
	return q
}
func (q MusicCollectionQuerier) Id(v any) MusicCollectionQuerier {
	q.super().Add(MusicCollection_Id, v)
	return q
}

type MusicCollectionUpdater cdb.BasicUpdater

func (u MusicCollectionUpdater) super() cdb.BasicUpdater { return cdb.BasicUpdater(u) }
func (u MusicCollectionUpdater) ToMap() map[string]any   { return u.super().ToMap() }
func (u MusicCollectionUpdater) IsEmpty() bool           { return u.super().IsEmpty() }
func (u MusicCollectionUpdater) Name(v string) MusicCollectionUpdater {
	u.super().Add(MusicCollection_Name, v)
	return u
}
func (u MusicCollectionUpdater) CoverUrl(v string) MusicCollectionUpdater {
	u.super().Add(MusicCollection_CoverUrl, v)
	return u
}
func (u MusicCollectionUpdater) SearchText(v string) MusicCollectionUpdater {
	u.super().Add(MusicCollection_SearchText, v)
	return u
}
