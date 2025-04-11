package mysql

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type AlbumDal interface {
	Save(album *Album) error
	GetById(id int64) (*Album, error)
	GetByRelatedId(id int64) ([]*Album, error)
	GetByMusicId(id int64) ([]*Album, error)
	GetByArtistId(id int64, offset, limit int32) ([]*Album, int32, error)
}
type albumDal struct {
	db *gorm.DB
}

func NewAlbumDal(db *gorm.DB) AlbumDal {
	return &albumDal{db: db}
}

type Album struct {
	Id             int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string `json:"name"`
	ArtistId       int64  `json:"artist_id"`
	MusicId        int64  `json:"music_id"`
	RelatedAlbumId int64  `json:"related_album_id"`
	CoverFileId    int64  `json:"cover_file_id"`
	IssueTime      int32  `json:"issue_time"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (f *albumDal) Save(album *Album) error {
	if err := f.db.Create(&album).Error; err != nil {
		log.Printf("save album err: %v", err)
		return err
	}
	return nil
}

func (f *albumDal) GetById(id int64) (*Album, error) {
	var model Album
	result := f.db.Where("id = ?", id).First(&model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &model, nil
}
func (f *albumDal) GetByRelatedId(id int64) ([]*Album, error) {
	var models []*Album
	if err := f.db.Where("related_album_id = ?", id).Find(&models).Error; err != nil {
		log.Printf("get album by music id err: %v", err)
		return nil, err
	}
	return models, nil

}
func (f *albumDal) GetByMusicId(id int64) ([]*Album, error) {
	var models []*Album
	if err := f.db.Where("music_id = ?", id).Find(&models).Error; err != nil {
		log.Printf("get album by music id err: %v", err)
		return nil, err
	}
	return models, nil
}
func (f *albumDal) GetByArtistId(id int64, offset, limit int32) ([]*Album, int32, error) {
	var models []*Album
	var total int64
	if err := f.db.Where("artist_id = ?", id).Model(&Album{}).Count(&total).Error; err != nil {
		log.Printf("count album by artist id err: %v", err)
		return nil, 0, err
	}
	if err := f.db.Where("artist_id = ?", id).Offset(int(offset)).Limit(int(limit)).Find(&models).Error; err != nil {
		log.Printf("get album by artist id err: %v", err)
		return nil, 0, err
	}
	return models, int32(total), nil
}
