package model

/*********** 表字段常量定义 **********/

import (
	"gorm.io/gorm"
	"onij/util/cdb"
	"time"
)

const (
	TableName_Practice = "practice"

	Practice_Id           = "id"
	Practice_PracticeType = "practice_type"
	Practice_Content      = "content"
	Practice_PracticeAt   = "practice_at"
	Practice_Duration     = "duration"
	Practice_CreatedAt    = "created_at"
	Practice_UpdatedAt    = "updated_at"
	Practice_DeletedAt    = "deleted_at"
)

/*********** 表结构定义 **********/

// Practice 表 practice 结构定义
type Practice struct {
	Id           int64          `gorm:"column:id" json:"id"`                       // 主键id
	PracticeType int32          `gorm:"column:practice_type" json:"practice_type"` // 练习类型
	Content      string         `gorm:"column:content" json:"content"`             // 练习内容
	PracticeAt   int64          `gorm:"column:practice_at" json:"practice_at"`     // 练习时间戳
	Duration     int32          `gorm:"column:duration" json:"duration"`           // 时长(分钟)
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`       // 创建时间
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`       // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`       // 删除时间
}

func NewPractice() *Practice {
	return &Practice{}
}

func (Practice) TableName() string {
	return TableName_Practice
}

func (Practice) UpdatableColumns() []string {
	return []string{
		Practice_PracticeType,
		Practice_Content,
		Practice_PracticeAt,
		Practice_Duration,
		Practice_CreatedAt,
		Practice_UpdatedAt,
		Practice_DeletedAt,
	}
}

/*********** 查询器与更新器 **********/

// PracticeQuerier 表 practice 查询器
type PracticeQuerier cdb.BasicQuerier

func (q PracticeQuerier) super() cdb.BasicQuerier {
	return (cdb.BasicQuerier)(q)
}

func (q PracticeQuerier) ToOptions() []cdb.Option {
	return q.super().ToOptions()
}

func (q PracticeQuerier) WithAsc(cols ...string) PracticeQuerier {
	q.super().WithAsc(cols...)
	return q
}

func (q PracticeQuerier) WithDesc(cols ...string) PracticeQuerier {
	q.super().WithDesc(cols...)
	return q
}

func (q PracticeQuerier) WithLimit(limit int) PracticeQuerier {
	q.super().WithLimit(limit)
	return q
}

func (q PracticeQuerier) WithUnscoped() PracticeQuerier {
	q.super().WithUnscoped()
	return q
}

func (q PracticeQuerier) IdxPracticeAt(practiceAt int64) PracticeQuerier {
	return q.PracticeAt(practiceAt)
}

func (q PracticeQuerier) Id(v any) PracticeQuerier {
	q.super().Add(Practice_Id, v)
	return q
}

func (q PracticeQuerier) PracticeType(v any) PracticeQuerier {
	q.super().Add(Practice_PracticeType, v)
	return q
}

func (q PracticeQuerier) Content(v any) PracticeQuerier {
	q.super().Add(Practice_Content, v)
	return q
}

func (q PracticeQuerier) PracticeAt(v any) PracticeQuerier {
	q.super().Add(Practice_PracticeAt, v)
	return q
}

func (q PracticeQuerier) Duration(v any) PracticeQuerier {
	q.super().Add(Practice_Duration, v)
	return q
}

func (q PracticeQuerier) CreatedAt(v any) PracticeQuerier {
	q.super().Add(Practice_CreatedAt, v)
	return q
}

func (q PracticeQuerier) UpdatedAt(v any) PracticeQuerier {
	q.super().Add(Practice_UpdatedAt, v)
	return q
}

func (q PracticeQuerier) DeletedAt(v any) PracticeQuerier {
	q.super().Add(Practice_DeletedAt, v)
	return q
}

// PracticeUpdater 表 practice 更新器
type PracticeUpdater cdb.BasicUpdater

func (u PracticeUpdater) super() cdb.BasicUpdater {
	return (cdb.BasicUpdater)(u)
}

func (u PracticeUpdater) Id(v int64) PracticeUpdater {
	u.super().Add(Practice_Id, v)
	return u
}

func (u PracticeUpdater) PracticeType(v int32) PracticeUpdater {
	u.super().Add(Practice_PracticeType, v)
	return u
}

func (u PracticeUpdater) Content(v string) PracticeUpdater {
	u.super().Add(Practice_Content, v)
	return u
}

func (u PracticeUpdater) PracticeAt(v int64) PracticeUpdater {
	u.super().Add(Practice_PracticeAt, v)
	return u
}

func (u PracticeUpdater) Duration(v int32) PracticeUpdater {
	u.super().Add(Practice_Duration, v)
	return u
}

func (u PracticeUpdater) CreatedAt(v time.Time) PracticeUpdater {
	u.super().Add(Practice_CreatedAt, v)
	return u
}

func (u PracticeUpdater) UpdatedAt(v time.Time) PracticeUpdater {
	u.super().Add(Practice_UpdatedAt, v)
	return u
}

func (u PracticeUpdater) DeletedAt(v time.Time) PracticeUpdater {
	u.super().Add(Practice_DeletedAt, v)
	return u
}

func (u PracticeUpdater) ToMap() map[string]any {
	return u.super().ToMap()
}

func (u PracticeUpdater) IsEmpty() bool {
	return u.super().IsEmpty()
}
