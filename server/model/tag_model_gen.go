package model

/*********** 表字段常量定义 **********/

import (
	"gorm.io/gorm"
	"onij/util/cdb"
	"time"
)

const (
	TableName_Tag = "tag"

	Tag_Id           = "id"
	Tag_ResourceId   = "resource_id"
	Tag_ResourceType = "resource_type"
	Tag_TagBiz       = "tag_biz"
	Tag_TagGroup     = "tag_group"
	Tag_TagType      = "tag_type"
	Tag_TargetId     = "target_id"
	Tag_TargetType   = "target_type"
	Tag_Extra        = "extra"
	Tag_CreatedAt    = "created_at"
	Tag_UpdatedAt    = "updated_at"
	Tag_DeletedAt    = "deleted_at"
)

/*********** 表结构定义 **********/

// Tag 表 tag 结构定义
type Tag struct {
	Id           int64          `gorm:"column:id" json:"id"`                       // 主键id
	ResourceId   int64          `gorm:"column:resource_id" json:"resource_id"`     // 资源id
	ResourceType int32          `gorm:"column:resource_type" json:"resource_type"` // 资源类型
	TagBiz       int32          `gorm:"column:tag_biz" json:"tag_biz"`             // 标签业务
	TagGroup     int32          `gorm:"column:tag_group" json:"tag_group"`         // 标签组
	TagType      int32          `gorm:"column:tag_type" json:"tag_type"`           // 标签类型
	TargetId     int64          `gorm:"column:target_id" json:"target_id"`         // 目标id
	TargetType   int32          `gorm:"column:target_type" json:"target_type"`     // 目标类型
	Extra        string         `gorm:"column:extra" json:"extra"`                 // 额外信息
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`       // 创建时间
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`       // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`       // 删除时间
}

func NewTag() *Tag {
	return &Tag{}
}

func (Tag) TableName() string {
	return TableName_Tag
}

func (Tag) UpdatableColumns() []string {
	return []string{
		Tag_ResourceType,
		Tag_TargetType,
		Tag_Extra,
		Tag_CreatedAt,
		Tag_UpdatedAt,
		Tag_DeletedAt,
	}
}

/*********** 查询器与更新器 **********/

// TagQuerier 表 tag 查询器
type TagQuerier cdb.BasicQuerier

func (q TagQuerier) super() cdb.BasicQuerier {
	return (cdb.BasicQuerier)(q)
}

func (q TagQuerier) ToOptions() []cdb.Option {
	return q.super().ToOptions()
}

func (q TagQuerier) WithAsc(cols ...string) TagQuerier {
	q.super().WithAsc(cols...)
	return q
}

func (q TagQuerier) WithDesc(cols ...string) TagQuerier {
	q.super().WithDesc(cols...)
	return q
}

func (q TagQuerier) WithLimit(limit int) TagQuerier {
	q.super().WithLimit(limit)
	return q
}

func (q TagQuerier) WithUnscoped() TagQuerier {
	q.super().WithUnscoped()
	return q
}

func (q TagQuerier) UkResourceBizGroupTypeTarget(resourceId int64, tagBiz int32, tagGroup int32, tagType int32, targetId int64) TagQuerier {
	return q.ResourceId(resourceId).TagBiz(tagBiz).TagGroup(tagGroup).TagType(tagType).TargetId(targetId)
}

func (q TagQuerier) Id(v any) TagQuerier {
	q.super().Add(Tag_Id, v)
	return q
}

func (q TagQuerier) ResourceId(v any) TagQuerier {
	q.super().Add(Tag_ResourceId, v)
	return q
}

func (q TagQuerier) ResourceType(v any) TagQuerier {
	q.super().Add(Tag_ResourceType, v)
	return q
}

func (q TagQuerier) TagBiz(v any) TagQuerier {
	q.super().Add(Tag_TagBiz, v)
	return q
}

func (q TagQuerier) TagGroup(v any) TagQuerier {
	q.super().Add(Tag_TagGroup, v)
	return q
}

func (q TagQuerier) TagType(v any) TagQuerier {
	q.super().Add(Tag_TagType, v)
	return q
}

func (q TagQuerier) TargetId(v any) TagQuerier {
	q.super().Add(Tag_TargetId, v)
	return q
}

func (q TagQuerier) TargetType(v any) TagQuerier {
	q.super().Add(Tag_TargetType, v)
	return q
}

func (q TagQuerier) Extra(v any) TagQuerier {
	q.super().Add(Tag_Extra, v)
	return q
}

func (q TagQuerier) CreatedAt(v any) TagQuerier {
	q.super().Add(Tag_CreatedAt, v)
	return q
}

func (q TagQuerier) UpdatedAt(v any) TagQuerier {
	q.super().Add(Tag_UpdatedAt, v)
	return q
}

func (q TagQuerier) DeletedAt(v any) TagQuerier {
	q.super().Add(Tag_DeletedAt, v)
	return q
}

// TagUpdater 表 tag 更新器
type TagUpdater cdb.BasicUpdater

func (u TagUpdater) super() cdb.BasicUpdater {
	return (cdb.BasicUpdater)(u)
}

func (u TagUpdater) Id(v int64) TagUpdater {
	u.super().Add(Tag_Id, v)
	return u
}

func (u TagUpdater) ResourceId(v int64) TagUpdater {
	u.super().Add(Tag_ResourceId, v)
	return u
}

func (u TagUpdater) ResourceType(v int32) TagUpdater {
	u.super().Add(Tag_ResourceType, v)
	return u
}

func (u TagUpdater) TagBiz(v int32) TagUpdater {
	u.super().Add(Tag_TagBiz, v)
	return u
}

func (u TagUpdater) TagGroup(v int32) TagUpdater {
	u.super().Add(Tag_TagGroup, v)
	return u
}

func (u TagUpdater) TagType(v int32) TagUpdater {
	u.super().Add(Tag_TagType, v)
	return u
}

func (u TagUpdater) TargetId(v int64) TagUpdater {
	u.super().Add(Tag_TargetId, v)
	return u
}

func (u TagUpdater) TargetType(v int32) TagUpdater {
	u.super().Add(Tag_TargetType, v)
	return u
}

func (u TagUpdater) Extra(v string) TagUpdater {
	u.super().Add(Tag_Extra, v)
	return u
}

func (u TagUpdater) CreatedAt(v time.Time) TagUpdater {
	u.super().Add(Tag_CreatedAt, v)
	return u
}

func (u TagUpdater) UpdatedAt(v time.Time) TagUpdater {
	u.super().Add(Tag_UpdatedAt, v)
	return u
}

func (u TagUpdater) DeletedAt(v time.Time) TagUpdater {
	u.super().Add(Tag_DeletedAt, v)
	return u
}

func (u TagUpdater) ToMap() map[string]any {
	return u.super().ToMap()
}

func (u TagUpdater) IsEmpty() bool {
	return u.super().IsEmpty()
}
