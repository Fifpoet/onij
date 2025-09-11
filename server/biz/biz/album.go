package biz

import (
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
	"onij/util/boost/exp"
	"onij/util/boost/tool"
	"time"
)

type AlbumPrime struct {
	Id          int64
	Name        *string
	Profile     *string
	ArtistIds   []int64
	AlbumType   *api.AlbumType
	IssueTime   *int32
	CoverFileId *int64
	LiveUrl     *string
	CreatedAt   *int64
	UpdatedAt   *int64

	CoverFile   *FilePrime
	Artists     []*ArtistPrime
	AlbumMusics []*AlbumMusicPrime
	TagPrime    []*TagPrime // album + musicIds
}

func (p *AlbumPrime) Model() *model.Album {
	if p == nil {
		return nil
	}
	return &model.Album{
		Id: func() int64 {
			if p.Id <= 0 {
				p.Id = util.IdGen.Generate()
			}
			return p.Id
		}(),
		Name:        exp.ValueOrZero(p.Name),
		Profile:     exp.ValueOrZero(p.Profile),
		AlbumType:   int32(exp.ValueOrZero(p.AlbumType)),
		ArtistIds:   tool.ToJson(p.ArtistIds),
		IssueTime:   time.Unix(int64(exp.ValueOrZero(p.IssueTime)), 0),
		CoverFileId: exp.ValueOrZero(p.CoverFileId),
		LiveUrl:     exp.ValueOrZero(p.LiveUrl),
	}
}

func Album(p *AlbumPrime) *api.Album {
	if p == nil {
		return nil
	}
	return &api.Album{
		Id:           p.Id,
		Name:         exp.ValueOrZero(p.Name),
		Profile:      exp.ValueOrZero(p.Profile),
		AlbumType:    exp.ValueOrZero(p.AlbumType),
		Artists:      collext.Pick(p.Artists, Artist),
		IssueTime:    exp.ValueOrZero(p.IssueTime),
		CoverFileUrl: p.CoverFile.Url(),
		LiveUrl:      exp.ValueOrZero(p.LiveUrl),
		AlbumMusics:  collext.Pick(p.AlbumMusics, AlbumMusic),
	}
}

func AlbumToBiz(a *model.Album) *AlbumPrime {
	if a == nil {
		return nil
	}
	return &AlbumPrime{
		Id:          a.Id,
		Name:        exp.Ptr(a.Name),
		Profile:     exp.Ptr(a.Profile),
		AlbumType:   exp.Ptr(api.AlbumType(a.AlbumType)),
		ArtistIds:   *tool.LoadJson[[]int64](a.ArtistIds, true),
		IssueTime:   exp.Ptr(int32(a.IssueTime.Unix())),
		CoverFileId: exp.Ptr(a.CoverFileId),
		LiveUrl:     exp.Ptr(a.LiveUrl),
		CreatedAt:   exp.Ptr(a.CreatedAt.Unix()),
		UpdatedAt:   exp.Ptr(a.UpdatedAt.Unix()),
	}
}
