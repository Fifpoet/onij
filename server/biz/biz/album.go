package biz

import (
	"gorm.io/gorm"
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/exp"
	"onij/util/boost/tool"
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

	CoverFile *FilePrime
	Artists   []*ArtistPrime
}

func (p *AlbumPrime) Model() *model.Album {
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
		IssueTime:   time.Time{},
		CoverFileId: 0,
		LiveUrl:     "",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		DeletedAt:   gorm.DeletedAt{},
	}
}
