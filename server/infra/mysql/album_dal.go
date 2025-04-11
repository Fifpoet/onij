package mysql

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

type AlbumDal interface {
	Upsert(album *Album, toUpdate map[string]any) error
	Save(Album *Album) error

	DelByIds(id []int) ([]*Album, error)

	GetByHash(key string) (*Album, error)
	GetByIds(ids []int) ([]*Album, error)
}
type albumDal struct {
	db *gorm.DB
}

func NewAlbumDal(db *gorm.DB) AlbumDal {
	return &AlbumDal{db: db}
}

type Album struct {
	Id          int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name"`
	ArtistId    int64  `json:"artist_id"`
	MusicId     int64  `json:"music_id"`
	CoverFileId int64  `json:"cover_file_id"`
	IssueTime   int32  `json:"issue_time"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (f *albumDal) Upsert(album *Album, toUpdate map[string]any) error {
	if album == nil && toUpdate == nil {
		return nil
	}
	if album != nil {
		if err := f.db.Create(&album).Error; err != nil {
			log.Printf("Upsert, save album err: %v", err)
			return err
		}
	}
	if toUpdate != nil {
		if err := f.db.Model(&Music{}).Updates(toUpdate).Error; err != nil {
			log.Printf("Upsert, update album err: %v", err)
			return err
		}
	}
	return nil
}

func (f *albumDal) Save(Album *Album) error {
	err := f.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(Album).Error
	if err != nil {
		log.Printf("save Album failed: err = %v \n", err)
		return err
	}
	return nil
}

func (f *albumDal) DelByIds(ids []int) ([]*Album, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	res, err := f.GetByIds(ids)
	if err != nil {
		return nil, err
	}
	err = f.db.Delete(&Album{}, "id IN ?", ids).Error
	if err != nil {
		log.Printf("DelByIds, delete Album failed: err = %v \n", err)
		return res, err
	}
	return res, nil
}

func (f *albumDal) GetByHash(hash string) (*Album, error) {
	res := &Album{}
	err := f.db.Where("hash = ?", hash).First(res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		log.Printf("GetByHash, get Album failed: err = %v \n", err)
		return nil, err
	}
	return res, nil
}

func (f *albumDal) GetByIds(ids []int) ([]*Album, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var res []*Album
	err := f.db.Where("id IN ?", ids).Find(&res).Error
	if err != nil {
		log.Printf("GetByIds, get Album failed: err = %v \n", err)
		return nil, err
	}
	return res, nil
}
