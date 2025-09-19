package model

/*********** 表字段常量定义 **********/

import (
	"gorm.io/gorm"
	"onij/util/cdb"
	"time"
)

const (
	TableName_AlbumMusic = "album_music"

	AlbumMusic_Id          = "id"
	AlbumMusic_Name        = "name"
	AlbumMusic_TimeLength  = "time_length"
	AlbumMusic_MusicId     = "music_id"
	AlbumMusic_AlbumId     = "album_id"
	AlbumMusic_ArtistNames = "artist_names"
	AlbumMusic_ThirdId     = "third_id"
	AlbumMusic_IsAvailable = "is_available"
	AlbumMusic_CreatedAt   = "created_at"
	AlbumMusic_UpdatedAt   = "updated_at"
	AlbumMusic_DeletedAt   = "deleted_at"
)

/*********** 表结构定义 **********/

// AlbumMusic 表 album_music 结构定义
type AlbumMusic struct {
	Id          int64          `gorm:"column:id" json:"id"`                     // 主键id
	Name        string         `gorm:"column:name" json:"name"`                 // 歌曲名称
	TimeLength  int32          `gorm:"column:time_length" json:"time_length"`   // 时长(s)
	MusicId     int64          `gorm:"column:music_id" json:"music_id"`         // 音乐id
	AlbumId     int64          `gorm:"column:album_id" json:"album_id"`         // 专辑id
	ArtistNames string         `gorm:"column:artist_names" json:"artist_names"` // 艺术家名称
	ThirdId     int64          `gorm:"column:third_id" json:"third_id"`         // 三方id
	IsAvailable int16          `gorm:"column:is_available" json:"is_available"` // 是否可用
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`     // 创建时间
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`     // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`     // 删除时间
}

func NewAlbumMusic() *AlbumMusic {
	return &AlbumMusic{}
}

func (AlbumMusic) TableName() string {
	return TableName_AlbumMusic
}

func (AlbumMusic) UpdatableColumns() []string {
	return []string{
		AlbumMusic_Name,
		AlbumMusic_TimeLength,
		AlbumMusic_MusicId,
		AlbumMusic_AlbumId,
		AlbumMusic_ArtistNames,
		AlbumMusic_ThirdId,
		AlbumMusic_IsAvailable,
		AlbumMusic_CreatedAt,
		AlbumMusic_UpdatedAt,
		AlbumMusic_DeletedAt,
	}
}

/*********** 查询器与更新器 **********/

// AlbumMusicQuerier 表 album_music 查询器
type AlbumMusicQuerier cdb.BasicQuerier

func (q AlbumMusicQuerier) super() cdb.BasicQuerier {
	return (cdb.BasicQuerier)(q)
}

func (q AlbumMusicQuerier) ToOptions() []cdb.Option {
	return q.super().ToOptions()
}

func (q AlbumMusicQuerier) WithAsc(cols ...string) AlbumMusicQuerier {
	q.super().WithAsc(cols...)
	return q
}

func (q AlbumMusicQuerier) WithDesc(cols ...string) AlbumMusicQuerier {
	q.super().WithDesc(cols...)
	return q
}

func (q AlbumMusicQuerier) WithLimit(limit int) AlbumMusicQuerier {
	q.super().WithLimit(limit)
	return q
}

func (q AlbumMusicQuerier) WithUnscoped() AlbumMusicQuerier {
	q.super().WithUnscoped()
	return q
}

func (q AlbumMusicQuerier) Id(v any) AlbumMusicQuerier {
	q.super().Add(AlbumMusic_Id, v)
	return q
}

func (q AlbumMusicQuerier) Name(v any) AlbumMusicQuerier {
	q.super().Add(AlbumMusic_Name, v)
	return q
}

func (q AlbumMusicQuerier) TimeLength(v any) AlbumMusicQuerier {
	q.super().Add(AlbumMusic_TimeLength, v)
	return q
}

func (q AlbumMusicQuerier) MusicId(v any) AlbumMusicQuerier {
	q.super().Add(AlbumMusic_MusicId, v)
	return q
}

func (q AlbumMusicQuerier) AlbumId(v any) AlbumMusicQuerier {
	q.super().Add(AlbumMusic_AlbumId, v)
	return q
}

func (q AlbumMusicQuerier) ArtistNames(v any) AlbumMusicQuerier {
	q.super().Add(AlbumMusic_ArtistNames, v)
	return q
}

func (q AlbumMusicQuerier) ThirdId(v any) AlbumMusicQuerier {
	q.super().Add(AlbumMusic_ThirdId, v)
	return q
}

func (q AlbumMusicQuerier) IsAvailable(v any) AlbumMusicQuerier {
	q.super().Add(AlbumMusic_IsAvailable, v)
	return q
}

func (q AlbumMusicQuerier) CreatedAt(v any) AlbumMusicQuerier {
	q.super().Add(AlbumMusic_CreatedAt, v)
	return q
}

func (q AlbumMusicQuerier) UpdatedAt(v any) AlbumMusicQuerier {
	q.super().Add(AlbumMusic_UpdatedAt, v)
	return q
}

func (q AlbumMusicQuerier) DeletedAt(v any) AlbumMusicQuerier {
	q.super().Add(AlbumMusic_DeletedAt, v)
	return q
}

// AlbumMusicUpdater 表 album_music 更新器
type AlbumMusicUpdater cdb.BasicUpdater

func (u AlbumMusicUpdater) super() cdb.BasicUpdater {
	return (cdb.BasicUpdater)(u)
}

func (u AlbumMusicUpdater) Id(v int64) AlbumMusicUpdater {
	u.super().Add(AlbumMusic_Id, v)
	return u
}

func (u AlbumMusicUpdater) Name(v string) AlbumMusicUpdater {
	u.super().Add(AlbumMusic_Name, v)
	return u
}

func (u AlbumMusicUpdater) TimeLength(v int32) AlbumMusicUpdater {
	u.super().Add(AlbumMusic_TimeLength, v)
	return u
}

func (u AlbumMusicUpdater) MusicId(v int64) AlbumMusicUpdater {
	u.super().Add(AlbumMusic_MusicId, v)
	return u
}

func (u AlbumMusicUpdater) AlbumId(v int64) AlbumMusicUpdater {
	u.super().Add(AlbumMusic_AlbumId, v)
	return u
}

func (u AlbumMusicUpdater) ArtistNames(v string) AlbumMusicUpdater {
	u.super().Add(AlbumMusic_ArtistNames, v)
	return u
}

func (u AlbumMusicUpdater) ThirdId(v int64) AlbumMusicUpdater {
	u.super().Add(AlbumMusic_ThirdId, v)
	return u
}

func (u AlbumMusicUpdater) IsAvailable(v int16) AlbumMusicUpdater {
	u.super().Add(AlbumMusic_IsAvailable, v)
	return u
}

func (u AlbumMusicUpdater) CreatedAt(v time.Time) AlbumMusicUpdater {
	u.super().Add(AlbumMusic_CreatedAt, v)
	return u
}

func (u AlbumMusicUpdater) UpdatedAt(v time.Time) AlbumMusicUpdater {
	u.super().Add(AlbumMusic_UpdatedAt, v)
	return u
}

func (u AlbumMusicUpdater) DeletedAt(v time.Time) AlbumMusicUpdater {
	u.super().Add(AlbumMusic_DeletedAt, v)
	return u
}

func (u AlbumMusicUpdater) ToMap() map[string]any {
	return u.super().ToMap()
}

func (u AlbumMusicUpdater) IsEmpty() bool {
	return u.super().IsEmpty()
}
