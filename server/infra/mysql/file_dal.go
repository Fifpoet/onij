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
	Save(files ...*File) error

	DelByIds(id []int64) ([]*File, error)

	GetByParentAndHash(parentId int64, keys ...string) ([]*File, error)
	GetByIds(ids ...int64) ([]*File, error)
	GetListByParentIdAndKeyword(parentId int64, keyword string, page util.Page) ([]*File, int32, error)
	GetFolderPathByParentId(parentId int64) ([]string, error)
	Delete(id int64) error
}
type fileDal struct {
	db *gorm.DB
}

func NewFileDal(db *gorm.DB) FileDal {
	return &fileDal{db: db}
}

// File 同一个文件夹下的hash去重
type File struct {
	Id       int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Format   int32  `json:"format"`
	StoreKey string `json:"store_key"`
	Hash     string `json:"hash" gorm:"not null;uniqueIndex:uk_hash"`
	ParentId int64  `json:"parent_id" gorm:"not null;uniqueIndex:uk_hash"`
	Size     int64  `json:"size"`

	OriginAt  int64          `json:"origin_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (f *fileDal) Delete(id int64) error {
	return f.db.Delete(&File{}, "id = ?", id).Error
}

func (f *fileDal) GetFolderPathByParentId(parentId int64) ([]string, error) {
	// 递归查询file, 直到parentId为0
	folders := []string{}
	for parentId != 0 {
		var file File
		err := f.db.Where("id = ?", parentId).First(&file).Error
		if err != nil {
			return nil, err
		}
		parentId = file.ParentId
		folders = append([]string{file.Name}, folders...)
	}
	return folders, nil
}

func (f *fileDal) GetListByParentIdAndKeyword(parentId int64, keyword string, page util.Page) ([]*File, int32, error) {
	var res []*File
	tx := f.db.Where("parent_id = ? AND name LIKE ?", parentId, "%"+keyword+"%").Model(&File{})
	count := int64(0)
	err := tx.Count(&count).Error
	if err != nil {
		log.Printf("GetByParentId, count file failed: err = %v \n", err)
		return nil, 0, err
	}
	err = tx.Offset(page.OffsetNum()).Limit(page.LimitNum()).Find(&res).Error
	if err != nil {
		log.Printf("GetByParentId, find file failed: err = %v \n", err)
		return nil, 0, err
	}
	return res, int32(count), nil
}

func (f *fileDal) Save(files ...*File) error {
	err := f.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(files).Error
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

func (f *fileDal) GetByParentAndHash(parentId int64, keys ...string) ([]*File, error) {
	res := []*File{}
	err := f.db.Where("parent_id = ? AND hash IN (?)", parentId, keys).Find(&res).Error
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
