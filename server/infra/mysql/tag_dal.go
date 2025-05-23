package mysql

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TagDal interface {
	GetByResource(ids ...int64) ([]*Tag, error)
	GetByGroupType(tagGroup, tagType *int32) ([]*Tag, error)
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
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

func (t *tagDal) GetByGroupType(tagGroup, tagType *int32) ([]*Tag, error) {
	var tags []*Tag
	if tagGroup == nil && tagType == nil {
		return nil, nil
	}
	tx := t.db
	if tagGroup != nil {
		tx = tx.Where("tag_group =?", tagGroup)
	}
	if tagType != nil {
		tx = tx.Where("tag_type =?", tagType)
	}
	err := tx.Find(&tags).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err!= nil {
		log.Printf("GetByGroupType, get tags error: %v \n", err)
		return nil, err
	}
	return tags, nil
}

func (t *tagDal) GetByResource(ids ...int64) ([]*Tag, error) {
	var tags []*Tag
	err := t.db.First(&tags, "id in ?", ids).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		log.Printf("GetByResource, get tags error: %v \n", err)
		return nil, err
	}
	return tags, nil
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
