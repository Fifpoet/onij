package mysql

import "time"

// Tag设计:
// 1. 需要联动 tod 图片 card poem, 一级业务分类

type TagDal interface {
}

type tagDal struct {
}

func NewTagDal() TagDal {
	return &tagDal{}
}

type Tag struct {
	Id         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string    `json:"name"`
	TagBiz     int       `json:"tag_biz"`
	TagGroup   int       `json:"tag_group"`
	TagType    int       `json:"tag_type"`
	Target     int       `json:"target"`
	TargetType int       `json:"target_type"`
	Extra      string    `json:"extra"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}
