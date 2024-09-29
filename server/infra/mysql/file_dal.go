package mysql

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"mime/multipart"
	"onij/util"
	"os"
	"path/filepath"
	"time"
)

type FileDal interface {
	CreateLocalFile(localFilePath string, biz int) (int, error)
	CreateFormFile(fileHeader *multipart.FileHeader, biz int) (int, error)
	DelByKey(key string) (*File, error)
	DelByIds(id []int) ([]*File, error)

	GetByKey(key string) (*File, error)
	GetByHash(hash string) (*File, error)
	GetByIds(ids []int) ([]*File, error)
	GetUrlByIds(ids ...int) ([]string, error)
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
	X    int    `json:"x"`
	Y    int    `json:"y"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (f *fileDal) CreateLocalFile(localFilePath string, biz int) (int, error) {
	// check hash, upload oss
	hash, x, y, err := util.GetLocalFileHash(localFilePath)
	if err != nil {
		return 0, err
	}
	fil, err := f.GetByHash(hash)
	if err != nil {
		return 0, err
	}
	if fil != nil {
		// 文件已存在
		log.Printf("CreateLocalFile, file already exist: file: %v \n", fil)
		return fil.Id, nil
	}

	key, err := util.UploadFile(localFilePath)
	if err != nil {
		return 0, err
	}
	stat, err := os.Stat(localFilePath)
	if err != nil {
		log.Printf("CreateLocalFile, get stat file failed: err = %v \n", err)
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
		X:    x,
		Y:    y,
	}
	err = f.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(newFil).Error
	if err != nil {
		log.Printf("CreateLocalFile, save file failed: err = %v \n", err)
		return 0, err
	}
	return newFil.Id, nil
}

// CreateFormFile 从表单文件中上传并写入 OSS
func (f *fileDal) CreateFormFile(fileHeader *multipart.FileHeader, biz int) (int, error) {
	// check hash, upload oss
	if fileHeader.Size == 0 {
		return 0, nil
	}
	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		log.Printf("CreateFormFile, open file failed: err = %v \n", err)
		return 0, err
	}
	hash, x, y, err := util.GetFileHash(file)
	if err != nil {
		return 0, err
	}
	fil, err := f.GetByHash(hash)
	if err != nil {
		return 0, err
	}
	if fil != nil {
		// 文件已存在
		log.Printf("CreateFormFile, file already exist: file: %v \n", fil)
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
		X:    x,
		Y:    y,
	}
	err = f.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(newFil).Error
	if err != nil {
		log.Printf("CreateFormFile, save file failed: err = %v \n", err)
		return 0, err
	}

	return newFil.Id, nil
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
	if err != nil {
		return nil, err
	}
	var res []string
	for _, fil := range fs {
		res = append(res, util.DownloadFile(fil.Key))
	}
	return res, nil
}
