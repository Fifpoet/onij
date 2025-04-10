package mysql

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type MusicDal interface {
	GetById(id int) (*Music, error)

	Upsert(music *Music, toUpdate map[string]any) error
	DelById(id int) (*Music, error)
	GetByTitle(title string) ([]*Music, error)
	GetByArtist(artistId int) ([]*Music, error)
	GetByTitleArtistPerType(title string, artistId, performType, page, size int) ([]*Music, error)
}

type musicDal struct {
	db *gorm.DB
}

func NewMusicDal(db *gorm.DB) MusicDal {
	return &musicDal{db: db}
}

type Music struct {
	Id          uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	RootId      int64  `json:"root_id"`
	Name        string `json:"name" gorm:"not null;uniqueIndex:uni_idx_music"`
	ArtistIds   string `json:"artist_ids" gorm:"not null;uniqueIndex:uni_idx_music"`
	ComposerId  int64  `json:"composer_id"`
	WriterId    int64  `json:"writer_id"`
	IssueTime   int32  `json:"issue_time"`
	PerformType int32  `json:"perform_type"`
	AlbumId     int64  `json:"album_id" gorm:"uniqueIndex:uni_idx_music"`
	MvUrl       string `json:"mv_url"`
	Mp3FileId   int64  `json:"mp3_file_id" gorm:"not null"`
	LyricFileId int64  `json:"lyric_file_id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (m *musicDal) GetById(id int) (*Music, error) {
	var music Music
	err := m.db.First(&music, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		log.Printf("GetById, get music error: %v \n", err)
		return nil, err
	}
	return &music, nil
}

func (m *musicDal) Upsert(music *Music, toUpdate map[string]any) error {
	if music == nil && toUpdate == nil {
		return nil
	}
	if music != nil {
		if err := m.db.Create(&music).Error; err != nil {
			log.Printf("Upsert, save music err: %v", err)
			return err
		}
	}
	if toUpdate != nil {
		if err := m.db.Model(&Music{}).Updates(toUpdate).Error; err != nil {
			log.Printf("Upsert, update music err: %v", err)
			return err
		}
	}
	return nil
}

func (m *musicDal) DelById(id int) (*Music, error) {
	mus, err := m.GetById(id)
	if err != nil {
		return nil, err
	}

	err = m.db.Delete(&Music{}, "id = ?", id).Error
	if err != nil {
		log.Printf("DelById, del music error: %v \n", err)
		return nil, err
	}
	return mus, nil
}

func (m *musicDal) GetByTitle(title string) ([]*Music, error) {
	var musics []*Music
	err := m.db.Where("title LIKE ?", "%"+title+"%").Find(&musics).Error
	if err != nil {
		log.Printf("GetByTitle, get music error: %v \n", err)
		return nil, err
	}
	return musics, nil
}

// GetByArtist 通过 artist_id 模糊匹配查询
func (m *musicDal) GetByArtist(artistId int) ([]*Music, error) {
	var musics []*Music
	artistIdStr := fmt.Sprintf(",%d,", artistId)
	err := m.db.Where("artist_ids LIKE ?", "%"+artistIdStr+"%").Find(&musics).Error
	if err != nil {
		log.Printf("GetByArtist, get music error: %v \n", err)
		return nil, err
	}
	return musics, nil
}

func (m *musicDal) GetByTitleArtistPerType(title string, artistId, performType, page, size int) ([]*Music, error) {
	var musics []*Music
	query := m.db

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if artistId != 0 {
		artistIdStr := fmt.Sprintf(",%d,", artistId)
		query = query.Where("artist_ids LIKE ?", "%"+artistIdStr+"%")
	}
	if performType != 0 {
		query = query.Where("perform_type = ?", performType)
	}
	err := query.Limit(size).Offset((page - 1) * size).Limit(size).Find(&musics).Error
	if err != nil {
		log.Printf("GetByTitleArtistPerType, get music error: %v \n", err)
		return nil, err
	}
	return musics, nil
}
