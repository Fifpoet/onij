package mysql

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type AlbumDal interface {
	Save(album *Album) error
}
type albumDal struct {
	db *gorm.DB
}

func NewAlbumDal(db *gorm.DB) AlbumDal {
	return &albumDal{db: db}
}

type Album struct {
	Id          int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	RootId      int64  `json:"root_id"`
	Name        string `json:"name"`
	ArtistId    int64  `json:"artist_id"`
	MusicId     int64  `json:"music_id"`
	CoverFileId int64  `json:"cover_file_id"`
	IssueTime   int32  `json:"issue_time"`

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
