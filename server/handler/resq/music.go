package resq

import (
	"mime/multipart"
	"onij/infra/mysql"
)

type UpsertMusicReq struct {
	RootId      int                   `json:"music_id"`
	Title       string                `json:"title"`
	ArtistIds   []int                 `json:"artist_ids"`
	Composer    int                   `json:"composer"`
	Writer      int                   `json:"writer"`
	IssueYear   int                   `json:"issue_year"`
	Language    string                `json:"language"`
	PerformType string                `json:"perform_type"`
	Instrument  string                `json:"instrument"`
	Concert     string                `json:"concert"`
	ConcertYear int                   `json:"concert_year"`
	Sequence    int                   `json:"sequence"`
	MvUrl       string                `json:"mv_url"`
	CoverOss    *multipart.FileHeader `json:"cover_oss"`
	MpOss       *multipart.FileHeader `json:"mp_oss"`
	LyricOss    *multipart.FileHeader `json:"lyric_oss"`
	SheetOss    *multipart.FileHeader `json:"sheet_oss"`
}

func (u *UpsertMusicReq) ToModel() (m *mysql.Music, cover, mp, lyric, sheet *multipart.FileHeader) {

	return nil, nil, nil, nil, nil
}
