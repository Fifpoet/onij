package resq

import (
	"mime/multipart"
	"onij/infra/mysql"
	"onij/util"
)

type UpsertMusicReq struct {
	Id          int                   `json:"id" form:"id"`
	RootId      int                   `json:"root_id" form:"root_id"`
	Title       string                `json:"title" form:"title"`
	ArtistIds   []int                 `json:"artist_ids" form:"artist_ids"`
	Composer    int                   `json:"composer" form:"composer"`
	Writer      int                   `json:"writer" form:"writer"`
	IssueYear   int                   `json:"issue_year" form:"issue_year"`
	Language    int                   `json:"language" form:"language"`
	PerformType int                   `json:"perform_type" form:"perform_type"`
	Instrument  int                   `json:"instrument" form:"instrument"`
	Concert     string                `json:"concert" form:"concert"`
	ConcertYear int                   `json:"concert_year" form:"concert_year"`
	Sequence    int                   `json:"sequence" form:"sequence"`
	MvUrl       string                `json:"mv_url" form:"mv_url"`
	CoverOss    int                   `json:"cover_oss" form:"cover_oss"`
	MpOss       int                   `json:"mp_oss" form:"mp_oss" gorm:"not null"`
	LyricOss    int                   `json:"lyric_oss" form:"lyric_oss"`
	SheetOss    int                   `json:"sheet_oss" form:"sheet_oss"`
	Cover       *multipart.FileHeader `form:"cover"`
	Mp          *multipart.FileHeader `form:"mp"`
	Lyric       *multipart.FileHeader `form:"lyric"`
	Sheet       *multipart.FileHeader `form:"sheet"`
}

func (u *UpsertMusicReq) ToModel() (m *mysql.Music, cover, mp, lyric, sheet *multipart.FileHeader) {
	return &mysql.Music{
		Id:          u.Id,
		RootId:      u.RootId,
		Title:       u.Title,
		ArtistIds:   util.ListToDb(u.ArtistIds),
		Composer:    u.Composer,
		Writer:      u.Writer,
		IssueYear:   u.IssueYear,
		Language:    u.Language,
		PerformType: u.PerformType,
		Instrument:  u.Instrument,
		Concert:     u.Concert,
		ConcertYear: u.ConcertYear,
		Sequence:    u.Sequence,
		MvUrl:       u.MvUrl,
		CoverOss:    u.CoverOss,
		MpOss:       u.MpOss,
		LyricOss:    u.LyricOss,
		SheetOss:    u.SheetOss,
	}, u.Cover, u.Mp, u.Lyric, u.Sheet
}
