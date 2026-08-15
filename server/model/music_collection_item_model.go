package model

import (
	"time"

	"gorm.io/gorm"
	"onij/util/cdb"
)

const (
	TableName_MusicCollectionItem = "music_collection_item"

	MusicCollectionItem_Id           = "id"
	MusicCollectionItem_CollectionId = "collection_id"
	MusicCollectionItem_SongId       = "song_id"
	MusicCollectionItem_SongName     = "song_name"
	MusicCollectionItem_ArtistNames  = "artist_names"
	MusicCollectionItem_SortOrder    = "sort_order"
	MusicCollectionItem_CreatedAt    = "created_at"
	MusicCollectionItem_UpdatedAt    = "updated_at"
	MusicCollectionItem_DeletedAt    = "deleted_at"
)

type MusicCollectionItem struct {
	Id           int64          `gorm:"column:id" json:"id"`
	CollectionId int64          `gorm:"column:collection_id" json:"collection_id"`
	SongId       int64          `gorm:"column:song_id" json:"song_id"`
	SongName     string         `gorm:"column:song_name" json:"song_name"`
	ArtistNames  string         `gorm:"column:artist_names" json:"artist_names"`
	SortOrder    int32          `gorm:"column:sort_order" json:"sort_order"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MusicCollectionItem) TableName() string { return TableName_MusicCollectionItem }

func (MusicCollectionItem) UpdatableColumns() []string {
	return []string{
		MusicCollectionItem_CollectionId,
		MusicCollectionItem_SongId,
		MusicCollectionItem_SongName,
		MusicCollectionItem_ArtistNames,
		MusicCollectionItem_SortOrder,
		MusicCollectionItem_CreatedAt,
		MusicCollectionItem_UpdatedAt,
		MusicCollectionItem_DeletedAt,
	}
}

type MusicCollectionItemQuerier cdb.BasicQuerier

func (q MusicCollectionItemQuerier) super() cdb.BasicQuerier { return cdb.BasicQuerier(q) }
func (q MusicCollectionItemQuerier) ToOptions() []cdb.Option { return q.super().ToOptions() }
func (q MusicCollectionItemQuerier) WithAsc(cols ...string) MusicCollectionItemQuerier {
	q.super().WithAsc(cols...)
	return q
}
func (q MusicCollectionItemQuerier) WithDesc(cols ...string) MusicCollectionItemQuerier {
	q.super().WithDesc(cols...)
	return q
}
func (q MusicCollectionItemQuerier) Id(v any) MusicCollectionItemQuerier {
	q.super().Add(MusicCollectionItem_Id, v)
	return q
}
func (q MusicCollectionItemQuerier) CollectionId(v any) MusicCollectionItemQuerier {
	q.super().Add(MusicCollectionItem_CollectionId, v)
	return q
}
func (q MusicCollectionItemQuerier) SongId(v any) MusicCollectionItemQuerier {
	q.super().Add(MusicCollectionItem_SongId, v)
	return q
}

type MusicCollectionItemUpdater cdb.BasicUpdater

func (u MusicCollectionItemUpdater) super() cdb.BasicUpdater { return cdb.BasicUpdater(u) }
func (u MusicCollectionItemUpdater) ToMap() map[string]any   { return u.super().ToMap() }
func (u MusicCollectionItemUpdater) IsEmpty() bool           { return u.super().IsEmpty() }
