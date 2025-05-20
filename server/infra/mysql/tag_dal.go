package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

type TagDal interface {
	Save(tags ...*Tag) error
	DelById(id int) error
}

type tagDal struct {
	db *gorm.DB
}

func NewTagDal(db *gorm.DB) TagDal {
	return &tagDal{db: db}
}

type Tag struct {
	Id         int64       `json:"id" gorm:"primaryKey;autoIncrement"`
	ResourceId     int64       `json:"origin" gorm:"index:uk_resource_biz_group_type_target,unique"`
	ResourceType int32       `json:"origin_type"`
	TagBiz     int32       `json:"tag_biz" gorm:"index:uk_resource_biz_group_type_target,unique"`
	TagGroup   int32       `json:"tag_group" gorm:"index:uk_resource_biz_group_type_target,unique"`
	TagType    int32       `json:"tag_type" gorm:"index:uk_resource_biz_group_type_target,unique"`
	TargetId     int64       `json:"target" gorm:"index:uk_resource_biz_group_type_target,unique"`
	TargetType int32       `json:"target_type"`
	ListShow bool 			`json:"list_show"`
	Extra      string    `json:"extra"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}


func (t *tagDal) Save(tags ...*Tag) error {
	err := t.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(tags).Error
	if err != nil {
		log.Printf("Save, save tag failed: err = %v \n", err)
		return  err
	}
	return nil
}

func (t *tagDal) DelById(id int) error {
	err := t.db.Where("id = ?", id).Delete(&Tag{}).Error
	if err != nil {
		log.Printf("DelById, delete tag failed: err = %v \n", err)
		return err
	}
	return nil
}
