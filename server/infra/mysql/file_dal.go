package mysql

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"onij/util"
	"os"
	"path/filepath"
	"time"
)

type FileDal interface {
	CreateFileFormLocal(localFilePath string, biz int) (string, error)
	CreateFileFromForm(fileHeader *multipart.FileHeader, biz int) (string, error)
	DelByKey(key string) error

	GetByKey(key string) (*File, error)
	GetById(id int) (*File, error)
}

type fileDal struct {
	db *gorm.DB
}

func NewFileDal(db *gorm.DB) FileDal {
	return &fileDal{db: db}
}

type File struct {
	Id   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
	Key  string `json:"key" gorm:"unique"`
	Biz  int    `json:"biz"`
	Size int    `json:"size"`
	Path string `json:"path"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (f *fileDal) CreateFileFormLocal(localFilePath string, biz int) (string, error) {
	uid := uuid.New()
	// upload oss
	err := util.UploadFile(localFilePath, uid.String())
	if err != nil {
		return "", err
	}
	stat, err := os.Stat(localFilePath)
	if err != nil {
		return "", err
	}

	// save db
	fil := &File{
		Name: filepath.Base(localFilePath),
		Key:  uid.String(),
		Biz:  biz,
		Size: int(stat.Size()),
		Path: localFilePath,
	}
	err = f.db.Save(fil).Error
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}

// CreateFileFromForm 从表单文件中上传并写入 OSS
func (f *fileDal) CreateFileFromForm(fileHeader *multipart.FileHeader, biz int) (string, error) {
	uid := uuid.New()

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// upload oss
	err = util.UploadFromReader(file, fileHeader.Size, uid.String())
	if err != nil {
		return "", err
	}

	// save db
	fil := &File{
		Name: fileHeader.Filename,
		Key:  uid.String(),
		Biz:  biz,
		Size: int(fileHeader.Size),
		Path: "",
	}
	err = f.db.Save(fil).Error
	if err != nil {
		return "", err
	}

	return uid.String(), nil
}

func (f *fileDal) DelByKey(key string) error {
	err := f.db.Delete(&File{}, "key = ?", key).Error
	if err != nil {
		return err
	}
	return util.DeleteFile(key)
}

func (f *fileDal) GetByKey(key string) (*File, error) {
	res := &File{}
	err := f.db.Where("key = ?", key).First(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (f *fileDal) GetById(id int) (*File, error) {
	res := &File{}
	err := f.db.Where("id = ?", id).First(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
