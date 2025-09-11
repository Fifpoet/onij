package biz

import (
	"onij/model"
	"onij/model/api"
	"onij/util/boost/exp"
)

type ArtistPrime struct {
	Id           int64
	Name         *string
	AvatarFileId *int64
	CreatedAt    *int64
	UpdatedAt    *int64

	AvatarFile *FilePrime
}

func (p *ArtistPrime) Model() *model.Artist {
	if p == nil {
		return nil
	}
	return &model.Artist{
		Id:           p.Id,
		Name:         exp.ValueOrZero(p.Name),
		AvatarFileId: exp.ValueOrZero(p.AvatarFileId),
	}
}

func Artist(p *ArtistPrime) *api.Artist {
	if p == nil {
		return nil
	}
	return &api.Artist{
		Id:            p.Id,
		Name:          exp.ValueOrZero(p.Name),
		AvatarFileUrl: p.AvatarFile.Url(),
	}
}

func ArtistToBiz(a *model.Artist) *ArtistPrime {
	if a == nil {
		return nil
	}
	return &ArtistPrime{
		Id:           a.Id,
		Name:         exp.Ptr(a.Name),
		AvatarFileId: exp.Ptr(a.AvatarFileId),
		CreatedAt:    exp.Ptr(a.CreatedAt.Unix()),
		UpdatedAt:    exp.Ptr(a.UpdatedAt.Unix()),
	}
}
