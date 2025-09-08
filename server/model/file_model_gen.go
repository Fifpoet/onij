package model

/*********** 表字段常量定义 **********/

import (
	"gorm.io/gorm"
	"onij/util/cdb"
	"time"
)

const (
	TableName_File = "file"

	File_Id        = "id"
	File_Name      = "name"
	File_Format    = "format"
	File_Size      = "size"
	File_StoreKey  = "store_key"
	File_Hash      = "hash"
	File_OriginAt  = "origin_at"
	File_CreatedAt = "created_at"
	File_UpdatedAt = "updated_at"
	File_DeletedAt = "deleted_at"
)

/*********** 表结构定义 **********/

// File 表 file 结构定义
type File struct {
	Id        int64          `gorm:"column:id" json:"id"`                 // 主键id
	Name      string         `gorm:"column:name" json:"name"`             // 文件名称
	Format    int32          `gorm:"column:format" json:"format"`         // 文件格式
	Size      int64          `gorm:"column:size" json:"size"`             // 文件大小
	StoreKey  string         `gorm:"column:store_key" json:"store_key"`   // 文件存储key
	Hash      string         `gorm:"column:hash" json:"hash"`             // 文件hash
	OriginAt  int64          `gorm:"column:origin_at" json:"origin_at"`   // 原始时间
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"` // 删除时间
}

func NewFile() *File {
	return &File{}
}

func (File) TableName() string {
	return TableName_File
}

func (File) UpdatableColumns() []string {
	return []string{
		File_Name,
		File_Format,
		File_Size,
		File_StoreKey,
		File_Hash,
		File_OriginAt,
		File_CreatedAt,
		File_UpdatedAt,
		File_DeletedAt,
	}
}

/*********** 查询器与更新器 **********/

// FileQuerier 表 file 查询器
type FileQuerier cdb.BasicQuerier

func (q FileQuerier) super() cdb.BasicQuerier {
	return (cdb.BasicQuerier)(q)
}

func (q FileQuerier) ToOptions() []cdb.Option {
	return q.super().ToOptions()
}

func (q FileQuerier) WithAsc(cols ...string) FileQuerier {
	q.super().WithAsc(cols...)
	return q
}

func (q FileQuerier) WithDesc(cols ...string) FileQuerier {
	q.super().WithDesc(cols...)
	return q
}

func (q FileQuerier) WithLimit(limit int) FileQuerier {
	q.super().WithLimit(limit)
	return q
}

func (q FileQuerier) WithUnscoped() FileQuerier {
	q.super().WithUnscoped()
	return q
}

func (q FileQuerier) Id(v any) FileQuerier {
	q.super().Add(File_Id, v)
	return q
}

func (q FileQuerier) Name(v any) FileQuerier {
	q.super().Add(File_Name, v)
	return q
}

func (q FileQuerier) Format(v any) FileQuerier {
	q.super().Add(File_Format, v)
	return q
}

func (q FileQuerier) Size(v any) FileQuerier {
	q.super().Add(File_Size, v)
	return q
}

func (q FileQuerier) StoreKey(v any) FileQuerier {
	q.super().Add(File_StoreKey, v)
	return q
}

func (q FileQuerier) Hash(v any) FileQuerier {
	q.super().Add(File_Hash, v)
	return q
}

func (q FileQuerier) OriginAt(v any) FileQuerier {
	q.super().Add(File_OriginAt, v)
	return q
}

func (q FileQuerier) CreatedAt(v any) FileQuerier {
	q.super().Add(File_CreatedAt, v)
	return q
}

func (q FileQuerier) UpdatedAt(v any) FileQuerier {
	q.super().Add(File_UpdatedAt, v)
	return q
}

func (q FileQuerier) DeletedAt(v any) FileQuerier {
	q.super().Add(File_DeletedAt, v)
	return q
}

// FileUpdater 表 file 更新器
type FileUpdater cdb.BasicUpdater

func (u FileUpdater) super() cdb.BasicUpdater {
	return (cdb.BasicUpdater)(u)
}

func (u FileUpdater) Id(v int64) FileUpdater {
	u.super().Add(File_Id, v)
	return u
}

func (u FileUpdater) Name(v string) FileUpdater {
	u.super().Add(File_Name, v)
	return u
}

func (u FileUpdater) Format(v int32) FileUpdater {
	u.super().Add(File_Format, v)
	return u
}

func (u FileUpdater) Size(v int64) FileUpdater {
	u.super().Add(File_Size, v)
	return u
}

func (u FileUpdater) StoreKey(v string) FileUpdater {
	u.super().Add(File_StoreKey, v)
	return u
}

func (u FileUpdater) Hash(v string) FileUpdater {
	u.super().Add(File_Hash, v)
	return u
}

func (u FileUpdater) OriginAt(v int64) FileUpdater {
	u.super().Add(File_OriginAt, v)
	return u
}

func (u FileUpdater) CreatedAt(v time.Time) FileUpdater {
	u.super().Add(File_CreatedAt, v)
	return u
}

func (u FileUpdater) UpdatedAt(v time.Time) FileUpdater {
	u.super().Add(File_UpdatedAt, v)
	return u
}

func (u FileUpdater) DeletedAt(v time.Time) FileUpdater {
	u.super().Add(File_DeletedAt, v)
	return u
}

func (u FileUpdater) ToMap() map[string]any {
	return u.super().ToMap()
}

func (u FileUpdater) IsEmpty() bool {
	return u.super().IsEmpty()
}
