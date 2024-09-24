package mysql

import (
	"gorm.io/gorm"
	"time"
)

type MusicDal interface {
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
	Title       string `json:"title"`
	ArtistIds   string `json:"artist_ids"`
	Composer    int    `json:"composer"`
	Writer      int    `json:"writer"`
	IssueYear   int    `json:"issue_year"`
	Language    string `json:"language"`
	PerformType string `json:"perform_type"`
	Instrument  string `json:"instrument"`
	Concert     string `json:"concert"`
	ConcertYear int    `json:"concert_year"`
	Sequence    int    `json:"sequence"`
	MvUrl       string `json:"mv_url"`
	CoverOss    string `json:"cover_oss"`
	MpOss       string `json:"mp_oss"`
	LyricOss    string `json:"lyric_oss"`
	SheetOss    string `json:"sheet_oss"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
