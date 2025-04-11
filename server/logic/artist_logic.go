package logic

import (
	"context"
	"onij/biz/prm"
	"onij/infra"
	"onij/infra/mysql"
	"onij/util"
)

type ArtistLogic interface {
	Upload(ctx context.Context, prm *prm.UploadArtistParam) (*prm.UploadArtistResult, error)
}

type artistLogic struct {
	*infra.AllInfra
}

func NewArtistLogic(i *infra.AllInfra) ArtistLogic {
	return &artistLogic{
		AllInfra: i,
	}
}

func (l *artistLogic) Upload(ctx context.Context, param *prm.UploadArtistParam) (*prm.UploadArtistResult, error) {
	if err := l.ArtistDal.Save(&mysql.Artist{
		Id:         util.IdGen.Generate(),
		Name:       param.Name,
		ArtistType: int32(param.ArtistType),
	}); err != nil {
		return nil, err
	}
	return &prm.UploadArtistResult{}, nil
}
