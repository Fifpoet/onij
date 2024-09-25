package resq

import (
	"mime/multipart"
	"onij/infra/mysql"
	"onij/util"
)

type UpsertMusicReq struct {
	Id          int                   `json:"id"`
	RootId      int                   `json:"music_id"`
	Title       string                `json:"title"`
	ArtistIds   []int                 `json:"artist_ids"`
	Composer    int                   `json:"composer"`
	Writer      int                   `json:"writer"`
	IssueYear   int                   `json:"issue_year"`
	Language    string                `json:"language"`
	PerformType int                   `json:"perform_type"`
	Instrument  string                `json:"instrument"`
	Concert     string                `json:"concert"`
	ConcertYear int                   `json:"concert_year"`
	Sequence    int                   `json:"sequence"`
	MvUrl       string                `json:"mv_url"`
	Cover       *multipart.FileHeader `json:"cover"`
	Mp          *multipart.FileHeader `json:"mp"`
	Lyric       *multipart.FileHeader `json:"lyric"`
	Sheet       *multipart.FileHeader `json:"sheet"`
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
	}, u.Cover, u.Mp, u.Lyric, u.Sheet
}
