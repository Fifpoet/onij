package model

/*********** 表字段常量定义 **********/

import (
	"gorm.io/gorm"
	"onij/util/cdb"
	"time"
)

const (
	TableName_Album = "album"

	Album_Id          = "id"
	Album_Name        = "name"
	Album_Profile     = "profile"
	Album_AlbumType   = "album_type"
	Album_ArtistIds   = "artist_ids"
	Album_IssueTime   = "issue_time"
	Album_CoverFileId = "cover_file_id"
	Album_LiveUrl     = "live_url"
	Album_ThirdId     = "third_id"
	Album_CreatedAt   = "created_at"
	Album_UpdatedAt   = "updated_at"
	Album_DeletedAt   = "deleted_at"
)

/*********** 表结构定义 **********/

// Album 表 album 结构定义
type Album struct {
	Id          int64          `gorm:"column:id" json:"id"`                       // 主键id
	Name        string         `gorm:"column:name" json:"name"`                   // 名称
	Profile     string         `gorm:"column:profile" json:"profile"`             // 简介
	AlbumType   int32          `gorm:"column:album_type" json:"album_type"`       // 专辑类型
	ArtistIds   string         `gorm:"column:artist_ids" json:"artist_ids"`       // 艺术家id列表
	IssueTime   time.Time      `gorm:"column:issue_time" json:"issue_time"`       // 发行时间
	CoverFileId int64          `gorm:"column:cover_file_id" json:"cover_file_id"` // 封面文件id
	LiveUrl     string         `gorm:"column:live_url" json:"live_url"`           // 视频地址
	ThirdId     int64          `gorm:"column:third_id" json:"third_id"`           // 三方id
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`       // 创建时间
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`       // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`       // 删除时间
}

func NewAlbum() *Album {
	return &Album{}
}

func (Album) TableName() string {
	return TableName_Album
}

func (Album) UpdatableColumns() []string {
	return []string{
		Album_Name,
		Album_Profile,
		Album_AlbumType,
		Album_ArtistIds,
		Album_IssueTime,
		Album_CoverFileId,
		Album_LiveUrl,
		Album_ThirdId,
		Album_CreatedAt,
		Album_UpdatedAt,
		Album_DeletedAt,
	}
}

/*********** 查询器与更新器 **********/

// AlbumQuerier 表 album 查询器
type AlbumQuerier cdb.BasicQuerier

func (q AlbumQuerier) super() cdb.BasicQuerier {
	return (cdb.BasicQuerier)(q)
}

func (q AlbumQuerier) ToOptions() []cdb.Option {
	return q.super().ToOptions()
}

func (q AlbumQuerier) WithAsc(cols ...string) AlbumQuerier {
	q.super().WithAsc(cols...)
	return q
}

func (q AlbumQuerier) WithDesc(cols ...string) AlbumQuerier {
	q.super().WithDesc(cols...)
	return q
}

func (q AlbumQuerier) WithLimit(limit int) AlbumQuerier {
	q.super().WithLimit(limit)
	return q
}

func (q AlbumQuerier) WithUnscoped() AlbumQuerier {
	q.super().WithUnscoped()
	return q
}

func (q AlbumQuerier) Id(v any) AlbumQuerier {
	q.super().Add(Album_Id, v)
	return q
}

func (q AlbumQuerier) Name(v any) AlbumQuerier {
	q.super().Add(Album_Name, v)
	return q
}

func (q AlbumQuerier) Profile(v any) AlbumQuerier {
	q.super().Add(Album_Profile, v)
	return q
}

func (q AlbumQuerier) AlbumType(v any) AlbumQuerier {
	q.super().Add(Album_AlbumType, v)
	return q
}

func (q AlbumQuerier) ArtistIds(v any) AlbumQuerier {
	q.super().Add(Album_ArtistIds, v)
	return q
}

func (q AlbumQuerier) IssueTime(v any) AlbumQuerier {
	q.super().Add(Album_IssueTime, v)
	return q
}

func (q AlbumQuerier) CoverFileId(v any) AlbumQuerier {
	q.super().Add(Album_CoverFileId, v)
	return q
}

func (q AlbumQuerier) LiveUrl(v any) AlbumQuerier {
	q.super().Add(Album_LiveUrl, v)
	return q
}

func (q AlbumQuerier) ThirdId(v any) AlbumQuerier {
	q.super().Add(Album_ThirdId, v)
	return q
}

func (q AlbumQuerier) CreatedAt(v any) AlbumQuerier {
	q.super().Add(Album_CreatedAt, v)
	return q
}

func (q AlbumQuerier) UpdatedAt(v any) AlbumQuerier {
	q.super().Add(Album_UpdatedAt, v)
	return q
}

func (q AlbumQuerier) DeletedAt(v any) AlbumQuerier {
	q.super().Add(Album_DeletedAt, v)
	return q
}

// AlbumUpdater 表 album 更新器
type AlbumUpdater cdb.BasicUpdater

func (u AlbumUpdater) super() cdb.BasicUpdater {
	return (cdb.BasicUpdater)(u)
}

func (u AlbumUpdater) Id(v int64) AlbumUpdater {
	u.super().Add(Album_Id, v)
	return u
}

func (u AlbumUpdater) Name(v string) AlbumUpdater {
	u.super().Add(Album_Name, v)
	return u
}

func (u AlbumUpdater) Profile(v string) AlbumUpdater {
	u.super().Add(Album_Profile, v)
	return u
}

func (u AlbumUpdater) AlbumType(v int32) AlbumUpdater {
	u.super().Add(Album_AlbumType, v)
	return u
}

func (u AlbumUpdater) ArtistIds(v string) AlbumUpdater {
	u.super().Add(Album_ArtistIds, v)
	return u
}

func (u AlbumUpdater) IssueTime(v time.Time) AlbumUpdater {
	u.super().Add(Album_IssueTime, v)
	return u
}

func (u AlbumUpdater) CoverFileId(v int64) AlbumUpdater {
	u.super().Add(Album_CoverFileId, v)
	return u
}

func (u AlbumUpdater) LiveUrl(v string) AlbumUpdater {
	u.super().Add(Album_LiveUrl, v)
	return u
}

func (u AlbumUpdater) ThirdId(v int64) AlbumUpdater {
	u.super().Add(Album_ThirdId, v)
	return u
}

func (u AlbumUpdater) CreatedAt(v time.Time) AlbumUpdater {
	u.super().Add(Album_CreatedAt, v)
	return u
}

func (u AlbumUpdater) UpdatedAt(v time.Time) AlbumUpdater {
	u.super().Add(Album_UpdatedAt, v)
	return u
}

func (u AlbumUpdater) DeletedAt(v time.Time) AlbumUpdater {
	u.super().Add(Album_DeletedAt, v)
	return u
}

func (u AlbumUpdater) ToMap() map[string]any {
	return u.super().ToMap()
}

func (u AlbumUpdater) IsEmpty() bool {
	return u.super().IsEmpty()
}
