package mysql

import (
	"errors"
	"fmt"
	"log"
	"onij/util"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MusicDal interface {
	GetById(id int) (*Music, error)
	GetByIds(ids ...int64) ([]*Music, error)

	Upsert(music *Music, toUpdate map[string]any) error
	DelById(id int) (*Music, error)
	GetByFullName(name string) ([]*Music, error)
	SearchByArtistAndNameAndTag(artistIds, writerIds, composeIds []int64, tagTypes []int32, name string, pageInfo util.Page) ([]*Music, error)
	GetByTitleArtistPerType(title string, artistId, performType, page, size int) ([]*Music, error)
}

type musicDal struct {
	db *gorm.DB
}

func NewMusicDal(db *gorm.DB) MusicDal {
	return &musicDal{db: db}
}

type Music struct {
	Id          int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	RootId      int64  `json:"root_id"`
	Name        string `json:"name" gorm:"not null;uniqueIndex:uni_idx_music"`
	FullName    string `json:"full_name"`
	ArtistIds   string `json:"artist_ids" gorm:"not null;uniqueIndex:uni_idx_music"`
	ComposerId  int64  `json:"composer_id"`
	WriterId    int64  `json:"writer_id"`
	IssueTime   int32  `json:"issue_time"`
	PerformType int32  `json:"perform_type"`
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
func (m *musicDal) GetByIds(ids ...int64) ([]*Music, error) {
	var musics []*Music
	err := m.db.First(&musics, "id in ?", ids).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		log.Printf("GetByIds, get music error: %v \n", err)
		return nil, err
	}
	return musics, nil
}

func (m *musicDal) Upsert(music *Music, toUpdate map[string]any) error {
	if toUpdate == nil {
		if err := m.db.Clauses(clause.OnConflict{DoNothing: true}).Create(music).Error; err != nil {
			log.Printf("Upsert, save music err: %v", err)
			return err
		}
	} else {
		if err := m.db.Model(music).Updates(toUpdate).Error; err != nil {
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

func (m *musicDal) GetByFullName(name string) ([]*Music, error) {
	var musics []*Music
	err := m.db.Where("full_name = ?", name).Find(&musics).Error
	if err != nil {
		log.Printf("GetByFullName, get music error: %v \n", err)
		return nil, err
	}
	return musics, nil
}

func (m *musicDal) SearchByArtistAndNameAndTag(artistIds, writerIds, composeIds []int64, tagTypes []int32, name string, pageInfo util.Page) ([]*Music, error) {
	var musics []*Music
	db := m.db

	// 处理artist_ids LIKE条件（满足一个即可）
	if len(artistIds) > 0 {
		var orConditions []string
		for _, id := range artistIds {
			orConditions = append(orConditions, fmt.Sprintf("artist_ids LIKE '%%%d%%'", id))
		}
		db = db.Where(strings.Join(orConditions, " OR "))
	}

	// 处理writer_id IN条件
	if len(writerIds) > 0 {
		db = db.Where("writer_id IN ?", writerIds)
	}

	// 处理composer_id IN条件
	if len(composeIds) > 0 {
		db = db.Where("composer_id IN ?", composeIds)
	}

	// 处理tag_type条件（需要连表查询）
	if len(tagTypes) > 0 {
		db = db.Joins("JOIN tag ON tag.resource_id = music.id AND tag.resource_type = ?", "music").
			Where("tag.tag_type IN ?", tagTypes)
	}

	// 处理name条件
	if name != "" {
		db = db.Where("full_name LIKE ?", "%"+name+"%")
	}

	err := db.
		Offset(pageInfo.OffsetNum()).
		Limit(pageInfo.LimitNum()).
		Find(&musics).Error

	if err != nil {
		log.Printf("SearchByArtistAndNameAndTag, get music error: %v \n", err)
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
