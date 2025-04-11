package biz

import (
	"onij/infra/mysql"
	"onij/util"
	"onij/util/boost/exp"
)

type AlbumPrime struct {
	Id           *int64
	Name         string
	ArtistIds    []int64
	Mp3FileId    int64
	LyricsFileId int64
	ComposerId   *int64
	WriterId     *int64
	AlbumId      *int64
	MvUrl        *string
	RootAlbumId  *int64
	IssueTime    *int32
}

func (m *AlbumPrime) Upsert() (album *mysql.Album, toUpdate map[string]any) {
	if m.Id == nil {
		return &mysql.Album{
			Id:          util.IdGen.Generate(),
			RootId:      exp.ValueOrZero(m.RootAlbumId),
			Name:        m.Name,
			ArtistIds:   util.Int64List2Str(m.ArtistIds),
			ComposerId:  exp.ValueOrZero(m.ComposerId),
			WriterId:    exp.ValueOrZero(m.WriterId),
			IssueTime:   exp.ValueOrZero(m.IssueTime),
			PerformType: 0,
			AlbumId:     exp.ValueOrZero(m.AlbumId),
			MvUrl:       exp.ValueOrZero(m.MvUrl),
			Mp3FileId:   m.Mp3FileId,
			LyricFileId: m.LyricsFileId,
		}, nil
	} else {
		toUpdate = map[string]any{
			"name":          m.Name,
			"artist_ids":    util.Int64List2Str(m.ArtistIds),
			"mp3_file_id":   m.Mp3FileId,
			"lyric_file_id": m.LyricsFileId,
		}
		if m.ComposerId != nil {
			toUpdate["composer_id"] = *m.ComposerId
		}
		if m.WriterId != nil {
			toUpdate["writer_id"] = *m.WriterId
		}
		if m.AlbumId != nil {
			toUpdate["album_id"] = *m.AlbumId
		}
		if m.MvUrl != nil {
			toUpdate["mv_url"] = *m.MvUrl
		}
		if m.RootAlbumId != nil {
			toUpdate["root_id"] = *m.RootAlbumId
		}
		if m.IssueTime != nil {
			toUpdate["issue_time"] = *m.IssueTime
		}
		return &mysql.Album{}, toUpdate
	}
}
