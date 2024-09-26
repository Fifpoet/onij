package mysql

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"onij/util"
	"os"
	"path/filepath"
	"time"
)

type FileDal interface {
	CreateFileFormLocal(localFilePath string, biz int) (int, error)
	CreateFileFromForm(fileHeader *multipart.FileHeader, biz int) (int, error)
	DelByKey(key string) (*File, error)
	DelByIds(id []int) ([]*File, error)

	GetByKey(key string) (*File, error)
	GetByHash(hash string) (*File, error)
	GetByIds(ids []int) ([]*File, error)
}

type fileDal struct {
	db *gorm.DB
}

func NewFileDal(db *gorm.DB) FileDal {
	return &fileDal{db: db}
}

// File .
// Key 为uuid, Hash为文件哈希值. upload接口返回
// 插入时有Hash相同的文件则 直接返回对应信息
// Path 不为空则代表本地文件路径
type File struct {
	Id   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
	Key  string `json:"key" gorm:"unique"`
	Biz  int    `json:"biz"`
	Size int    `json:"size"`
	Path string `json:"path"`
	Hash string `json:"hash"  gorm:"unique"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (f *fileDal) CreateFileFormLocal(localFilePath string, biz int) (int, error) {
	// check hash, upload oss
	hash, err := util.GetLocalFileHash(localFilePath)
	if err != nil {
		return 0, err
	}
	fil, err := f.GetByHash(hash)
	if err != nil {
		return 0, err
	}
	if fil != nil {
		// 文件已存在
		return fil.Id, nil
	}

	key, err := util.UploadFile(localFilePath)
	if err != nil {
		return 0, err
	}
	stat, err := os.Stat(localFilePath)
	if err != nil {
		log.Printf("CreateFileFormLocal, get stat file failed: err = %v \n", err)
		return 0, err
	}

	// save db
	newFil := &File{
		Name: filepath.Base(localFilePath),
		Key:  key,
		Biz:  biz,
		Size: int(stat.Size()),
		Path: localFilePath,
		Hash: hash,
	}
	err = f.db.Save(newFil).Error
	if err != nil {
		log.Printf("CreateFileFormLocal, save file failed: err = %v \n", err)
		return 0, err
	}
	return fil.Id, nil
}

// CreateFileFromForm 从表单文件中上传并写入 OSS
func (f *fileDal) CreateFileFromForm(fileHeader *multipart.FileHeader, biz int) (int, error) {
	// check hash, upload oss
	if fileHeader == nil {
		return 0, nil
	}
	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		log.Printf("CreateFileFromForm, open file failed: err = %v \n", err)
		return 0, err
	}
	hash, err := util.GetFileHash(file)
	if err != nil {
		return 0, err
	}
	fil, err := f.GetByHash(hash)
	if err != nil {
		return 0, err
	}
	if fil != nil {
		// 文件已存在
		return fil.Id, nil
	}

	// upload oss
	key, err := util.UploadFromReader(file, fileHeader.Size)
	if err != nil {
		return 0, err
	}

	// save db
	newFil := &File{
		Name: fileHeader.Filename,
		Key:  key,
		Biz:  biz,
		Size: int(fileHeader.Size),
		Path: "",
		Hash: hash,
	}
	err = f.db.Save(newFil).Error
	if err != nil {
		log.Printf("CreateFileFromForm, save file failed: err = %v \n", err)
		return 0, err
	}

	return fil.Id, nil
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
	res, err := f.GetByIds(ids)
	if err != nil {
		return nil, err
	}
	err = f.db.Delete(&File{}, "id = ?", ids).Error
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

func (f *fileDal) GetByIds(ids []int) ([]*File, error) {
	var res []*File
	err := f.db.Where("id = ?", ids).First(res).Error
	if err != nil {
		log.Printf("GetByIds, get file failed: err = %v \n", err)
		return nil, err
	}
	return res, nil
}
