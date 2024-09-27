package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type MusicDal interface {
	GetById(id int) (*Music, error)

	Save(music *Music) (int, error)
	DelById(id int) (*Music, error)
}

type musicDal struct {
	db *gorm.DB
}

func NewMusicDal(db *gorm.DB) MusicDal {
	return &musicDal{db: db}
}

type Music struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	RootId      int    `json:"root_id"`
	Title       string `json:"title" gorm:"not null;uniqueIndex:uni_idx_music"`
	ArtistIds   string `json:"artist_ids" gorm:"not null;uniqueIndex:uni_idx_music"`
	Composer    int    `json:"composer"`
	Writer      int    `json:"writer"`
	Length      int    `json:"length"`
	IssueYear   int    `json:"issue_year"`
	Language    int    `json:"language"`
	PerformType int    `json:"perform_type"`
	Instrument  int    `json:"instrument"`
	Concert     string `json:"concert" gorm:"uniqueIndex:uni_idx_music"`
	ConcertYear int    `json:"concert_year"`
	Sequence    int    `json:"sequence"`
	MvUrl       string `json:"mv_url"`
	CoverOss    int    `json:"cover_oss"`
	MpOss       int    `json:"mp_oss" gorm:"not null"`
	LyricOss    int    `json:"lyric_oss"`
	SheetOss    int    `json:"sheet_oss"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (m *musicDal) GetById(id int) (*Music, error) {
	var music Music
	err := m.db.First(&music, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &music, nil
}

func (m *musicDal) Save(music *Music) (int, error) {
	err := m.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(music).Error
	if err != nil {
		return 0, err
	}
	return music.Id, nil
}

func (m *musicDal) DelById(id int) (*Music, error) {
	mus, err := m.GetById(id)
	if err != nil {
		return nil, err
	}

	err = m.db.Delete(&Music{}, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return mus, nil
}
