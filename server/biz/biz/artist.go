package biz

import (
	"onij/model"
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

func (a *ArtistPrime) Model() *model.Artist {
	return &model.Artist{
		Id:           a.Id,
		Name:         exp.ValueOrZero(a.Name),
		AvatarFileId: exp.ValueOrZero(a.AvatarFileId),
	}
}
