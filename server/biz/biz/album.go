package biz

import (
	"onij/infra/mysql"
	"onij/util"
	"onij/util/boost/exp"
)

type AlbumPrime struct {
	Name        string
	ArtistId    int64
	CoverFileId int64
	IssueTime   int32
	AlbumId     *int64
	MusicId     *int64
}

func (a *AlbumPrime) Model() *mysql.Album {
	if a.MusicId == nil {
		return &mysql.Album{
			Id:          util.IdGen.Generate(),
			RootId:      0,
			Name:        a.Name,
			ArtistId:    a.ArtistId,
			CoverFileId: a.CoverFileId,
			IssueTime:   a.IssueTime,
			MusicId:     exp.ValueOrZero(a.MusicId),
		}
	}
	return &mysql.Album{
		Id:          util.IdGen.Generate(),
		RootId:      exp.ValueOrZero(a.AlbumId),
		Name:        a.Name,
		ArtistId:    a.ArtistId,
		CoverFileId: a.CoverFileId,
		IssueTime:   a.IssueTime,
		MusicId:     exp.ValueOrZero(a.MusicId),
	}
}
