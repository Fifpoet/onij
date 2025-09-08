package model

/*********** 表字段常量定义 **********/

import (
	cdb "code.chenji.com/pkg/common/component/db"
	"gorm.io/gorm"
	"time"
)

const (
	TableName_Artist = "artist"

	Artist_Id        = "id"
	Artist_Name      = "name"
	Artist_Avatar    = "avatar"
	Artist_CreatedAt = "created_at"
	Artist_UpdatedAt = "updated_at"
	Artist_DeletedAt = "deleted_at"
)

/*********** 表结构定义 **********/

// Artist 表 artist 结构定义
type Artist struct {
	Id        int64          `gorm:"column:id" json:"id"`                 // 主键id
	Name      string         `gorm:"column:name" json:"name"`             // 艺术家名字
	Avatar    string         `gorm:"column:avatar" json:"avatar"`         // 艺术家头像
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"` // 删除时间
}

func NewArtist() *Artist {
	return &Artist{}
}

func (Artist) TableName() string {
	return TableName_Artist
}

func (Artist) UpdatableColumns() []string {
	return []string{
		Artist_Name,
		Artist_Avatar,
		Artist_CreatedAt,
		Artist_UpdatedAt,
		Artist_DeletedAt,
	}
}

/*********** 查询器与更新器 **********/

// ArtistQuerier 表 artist 查询器
type ArtistQuerier cdb.BasicQuerier

func (q ArtistQuerier) super() cdb.BasicQuerier {
	return (cdb.BasicQuerier)(q)
}

func (q ArtistQuerier) ToOptions() []cdb.Option {
	return q.super().ToOptions()
}

func (q ArtistQuerier) WithAsc(cols ...string) ArtistQuerier {
	q.super().WithAsc(cols...)
	return q
}

func (q ArtistQuerier) WithDesc(cols ...string) ArtistQuerier {
	q.super().WithDesc(cols...)
	return q
}

func (q ArtistQuerier) WithLimit(limit int) ArtistQuerier {
	q.super().WithLimit(limit)
	return q
}

func (q ArtistQuerier) WithUnscoped() ArtistQuerier {
	q.super().WithUnscoped()
	return q
}

func (q ArtistQuerier) Id(v any) ArtistQuerier {
	q.super().Add(Artist_Id, v)
	return q
}

func (q ArtistQuerier) Name(v any) ArtistQuerier {
	q.super().Add(Artist_Name, v)
	return q
}

func (q ArtistQuerier) Avatar(v any) ArtistQuerier {
	q.super().Add(Artist_Avatar, v)
	return q
}

func (q ArtistQuerier) CreatedAt(v any) ArtistQuerier {
	q.super().Add(Artist_CreatedAt, v)
	return q
}

func (q ArtistQuerier) UpdatedAt(v any) ArtistQuerier {
	q.super().Add(Artist_UpdatedAt, v)
	return q
}

func (q ArtistQuerier) DeletedAt(v any) ArtistQuerier {
	q.super().Add(Artist_DeletedAt, v)
	return q
}

// ArtistUpdater 表 artist 更新器
type ArtistUpdater cdb.BasicUpdater

func (u ArtistUpdater) super() cdb.BasicUpdater {
	return (cdb.BasicUpdater)(u)
}

func (u ArtistUpdater) Id(v int64) ArtistUpdater {
	u.super().Add(Artist_Id, v)
	return u
}

func (u ArtistUpdater) Name(v string) ArtistUpdater {
	u.super().Add(Artist_Name, v)
	return u
}

func (u ArtistUpdater) Avatar(v string) ArtistUpdater {
	u.super().Add(Artist_Avatar, v)
	return u
}

func (u ArtistUpdater) CreatedAt(v time.Time) ArtistUpdater {
	u.super().Add(Artist_CreatedAt, v)
	return u
}

func (u ArtistUpdater) UpdatedAt(v time.Time) ArtistUpdater {
	u.super().Add(Artist_UpdatedAt, v)
	return u
}

func (u ArtistUpdater) DeletedAt(v time.Time) ArtistUpdater {
	u.super().Add(Artist_DeletedAt, v)
	return u
}

func (u ArtistUpdater) ToMap() map[string]any {
	return u.super().ToMap()
}

func (u ArtistUpdater) IsEmpty() bool {
	return u.super().IsEmpty()
}
