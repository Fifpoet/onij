package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

type MetaDal interface {
	GetByMetaEnumCode(metaEnumCode []int) ([]*Meta, error)
	Save(metas []*Meta) error
}

type metaDal struct {
	db *gorm.DB
}

func NewMetaDal(db *gorm.DB) MetaDal {
	return &metaDal{db: db}
}

type Meta struct {
	Id           int            `json:"id" gorm:"primaryKey;autoIncrement"`
	MetaEnumCode int            `json:"meta_enum_code" gorm:"uniqueIndex:uni_meta"`
	Value        int            `json:"value" gorm:"uniqueIndex:uni_meta"`
	Name         string         `json:"name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

// GetByMetaEnumCode retrieves all Meta records that match the given MetaEnumCode
func (m *metaDal) GetByMetaEnumCode(metaEnumCode []int) ([]*Meta, error) {
	var metas []*Meta
	err := m.db.Where("meta_enum_code IN ?", metaEnumCode).Order("created_at desc").Find(&metas).Error
	if err != nil {
		return nil, err
	}
	return metas, nil
}

func (m *metaDal) Save(metas []*Meta) error {
	err := m.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(metas).Error
	if err != nil {
		log.Printf("Save, save meta err: %v", err)
		return err
	}
	return err
}
