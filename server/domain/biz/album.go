package biz

import (
	"onij/infra/mysql"
	"onij/util"
	"onij/util/boost/exp"
)

type AlbumPrime struct {
	Id             int64
	Name           string
	ArtistId       int64
	CoverFileId    int64
	IssueTime      int32
	RelatedAlbumId *int64
	MusicId        *int64
}

func (a *AlbumPrime) Model() *mysql.Album {
	if a.MusicId == nil {
		return &mysql.Album{
			Id:             util.IdGen.Generate(),
			RelatedAlbumId: 0,
			Name:           a.Name,
			ArtistId:       a.ArtistId,
			CoverFileId:    a.CoverFileId,
			IssueTime:      a.IssueTime,
			MusicId:        exp.ValueOrZero(a.MusicId),
		}
	}
	return &mysql.Album{
		Id:             util.IdGen.Generate(),
		Name:           a.Name,
		ArtistId:       a.ArtistId,
		CoverFileId:    a.CoverFileId,
		IssueTime:      a.IssueTime,
		MusicId:        exp.ValueOrZero(a.MusicId),
		RelatedAlbumId: exp.ValueOrZero(a.RelatedAlbumId),
	}
}
