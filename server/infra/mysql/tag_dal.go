package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type TagDal interface {
	GetByBizGroupTypeOriginTarget(biz, group, tagType, origin, target int) ([]*Tag, error)
	Save(tag *Tag) (int, error)
	DelById(id int) error
}

type tagDal struct {
	db *gorm.DB
}

func NewTagDal(db *gorm.DB) TagDal {
	return &tagDal{db: db}
}

type Tag struct {
	Id         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Origin     int       `json:"origin" gorm:"index:uk_origin_biz_group_type_target,unique"`
	OriginType int       `json:"origin_type"`
	TagBiz     int       `json:"tag_biz" gorm:"index:uk_origin_biz_group_type_target,unique"`
	TagGroup   int       `json:"tag_group" gorm:"index:uk_origin_biz_group_type_target,unique"`
	TagType    int       `json:"tag_type" gorm:"index:uk_origin_biz_group_type_target,unique"`
	Target     int       `json:"target" gorm:"index:uk_origin_biz_group_type_target,unique"`
	TargetType int       `json:"target_type"`
	Extra      string    `json:"extra"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

func (t *tagDal) GetByBizGroupTypeOriginTarget(biz, group, tagType, origin, target int) ([]*Tag, error) {
	var tags []*Tag
	query := t.db

	if biz != 0 {
		query = query.Where("tag_biz = ?", biz)
	}
	if group != 0 {
		query = query.Where("tag_group = ?", group)
	}
	if tagType != 0 {
		query = query.Where("tag_type = ?", tagType)
	}
	if origin != 0 {
		query = query.Where("origin = ?", origin)
	}
	if target != 0 {
		query = query.Where("target = ?", target)
	}

	err := query.Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (t *tagDal) Save(tag *Tag) (int, error) {
	err := t.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(tag).Error
	if err != nil {
		return 0, err
	}
	return tag.Id, nil
}

func (t *tagDal) DelById(id int) error {
	err := t.db.Where("id = ?", id).Delete(&Tag{}).Error
	return err
}
