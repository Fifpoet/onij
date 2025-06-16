package mysql

import (
	"errors"
	"log"
	"onij/util"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FileDal interface {
	Save(file *File) error

	DelByIds(id []int64) ([]*File, error)

	GetByHash(key string) (*File, error)
	GetByIds(ids ...int64) ([]*File, error)
	GetByParentId(parentId int64, page util.Page) ([]*File, error)
	CountByParentId(parentId int64) (int32, error)
}
type fileDal struct {
	db *gorm.DB
}

func NewFileDal(db *gorm.DB) FileDal {
	return &fileDal{db: db}
}

type File struct {
	Id       int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Format   int32  `json:"format"`
	StoreKey string `json:"store_key"`
	Hash     string `json:"hash" gorm:"not null;uniqueIndex:uk_hash"`
	ParentId int64  `json:"parent_id"`

	OriginAt  *time.Time     `json:"origin_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (f *fileDal) GetByParentId(parentId int64, page util.Page) ([]*File, error) {
	var res []*File
	err := f.db.Where("parent_id = ?", parentId).Offset(page.OffsetNum()).Limit(page.LimitNum()).Find(&res).Error
	if err != nil {
		log.Printf("GetByParentId, get file failed: err = %v \n", err)
		return nil, err
	}
	return res, nil
}

func (f *fileDal) Save(file *File) error {
	err := f.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(file).Error
	if err != nil {
		log.Printf("save file failed: err = %v \n", err)
		return err
	}
	return nil
}

func (f *fileDal) DelByIds(ids []int64) ([]*File, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	res, err := f.GetByIds(ids...)
	if err != nil {
		return nil, err
	}
	err = f.db.Delete(&File{}, "id IN ?", ids).Error
	if err != nil {
		log.Printf("DelByIds, delete file failed: err = %v \n", err)
		return res, err
	}
	return res, nil
}

func (f *fileDal) GetByHash(hash string) (*File, error) {
	res := &File{}
	err := f.db.Where("hash = ?", hash).First(res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		log.Printf("GetByHash, get file failed: err = %v \n", err)
		return nil, err
	}
	return res, nil
}

func (f *fileDal) GetByIds(ids ...int64) ([]*File, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var res []*File
	err := f.db.Where("id IN ?", ids).Find(&res).Error
	if err != nil {
		log.Printf("GetByIds, get file failed: err = %v \n", err)
		return nil, err
	}
	return res, nil
}

func (f *fileDal) CountByParentId(parentId int64) (int32, error) {
	var count int64
	err := f.db.Model(&File{}).Where("parent_id = ?", parentId).Count(&count).Error
	if err != nil {
		log.Printf("CountByParentId, get file count failed: err = %v \n", err)
		return 0, err
	}
	return int32(count), nil
}
