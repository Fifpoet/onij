package biz

import (
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/exp"
)

type AlbumMusicPrime struct {
	Id          int64
	Name        *string
	TimeLength  *int32
	MusicId     *int64
	AlbumId     *int64
	IsAvailable *bool
}

func (p *AlbumMusicPrime) Model() *model.AlbumMusic {
	if p == nil {
		return nil
	}
	return &model.AlbumMusic{
		Id:          p.Id,
		Name:        exp.ValueOrZero(p.Name),
		TimeLength:  exp.ValueOrZero(p.TimeLength),
		MusicId:     exp.ValueOrZero(p.MusicId),
		AlbumId:     exp.ValueOrZero(p.AlbumId),
		IsAvailable: util.BoolToInt16(exp.ValueOrZero(p.IsAvailable)),
	}
}
func AlbumMusic(p *AlbumMusicPrime) *api.AlbumMusic {
	return &api.AlbumMusic{
		Name:        exp.ValueOrZero(p.Name),
		TimeLength:  exp.ValueOrZero(p.TimeLength),
		MusicId:     p.Id,
		IsAvailable: exp.ValueOrZero(p.IsAvailable),
	}
}
func AlbumMusicToBiz(a *model.AlbumMusic) *AlbumMusicPrime {
	if a == nil {
		return nil
	}
	return &AlbumMusicPrime{
		Id:          a.Id,
		Name:        exp.Ptr(a.Name),
		TimeLength:  exp.Ptr(a.TimeLength),
		MusicId:     exp.Ptr(a.MusicId),
		AlbumId:     exp.Ptr(a.AlbumId),
		IsAvailable: exp.Ptr(a.IsAvailable == 1),
	}
}
