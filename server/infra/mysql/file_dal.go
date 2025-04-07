package mysql

import (
	"code.chenji.com/pkg/boost/collection/collext"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"onij/util"
	"time"
)

type FileDal interface {
	Save(file *File) error

	DelByKey(key string) (*File, error)
	DelByIds(id []int) ([]*File, error)

	GetByKey(key string) (*File, error)
	GetByIds(ids []int) ([]*File, error)
	GetUrlByIds(ids ...int) ([]string, error)
}
type fileDal struct {
	db *gorm.DB
}

func NewFileDal(db *gorm.DB) FileDal {
	return &fileDal{db: db}
}

type File struct {
	Id       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Format   int    `json:"format"`
	StoreKey string `json:"store_key" gorm:"unique"`
	Hash     string `json:"hash"  gorm:"unique"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
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

func (f *fileDal) DelByKey(key string) (*File, error) {
	res, err := f.GetByKey(key)
	if err != nil {
		return nil, err
	}
	err = f.db.Delete(&File{}, "key = ?", key).Error
	if err != nil {
		log.Printf("DelByKey, delete file failed: err = %v \n", err)
		return res, err
	}
	return res, nil
}
func (f *fileDal) DelByIds(ids []int) ([]*File, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	res, err := f.GetByIds(ids)
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

func (f *fileDal) GetByKey(key string) (*File, error) {
	res := &File{}
	err := f.db.Where("key = ?", key).First(res).Error
	if err != nil {
		log.Printf("GetByKey, get file failed: err = %v \n", err)
		return nil, err
	}
	return res, nil
}

func (f *fileDal) GetByIds(ids []int) ([]*File, error) {
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

func (f *fileDal) GetUrlByIds(ids ...int) ([]string, error) {
	fs, err := f.GetByIds(ids)
	fsIdx := collext.Map(fs, func(fil *File) int { return fil.Id })
	if err != nil {
		return nil, err
	}
	var res []string
	for _, id := range ids {
		if fsIdx[id] != nil && fsIdx[id].StoreKey != "" {
			res = append(res, util.DownloadFile(fsIdx[id].StoreKey))
		} else {
			res = append(res, "")
		}
	}
	return res, nil
}
