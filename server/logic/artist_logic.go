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
	Search(ctx context.Context, prm *prm.SearchArtistParam) (*prm.SearchArtistResult, error)
}

type artistLogic struct {
	*infra.AllInfra
}

func NewArtistLogic(i *infra.AllInfra) ArtistLogic {
	return &artistLogic{
		AllInfra: i,
	}
}

func (l *artistLogic) Search(ctx context.Context, param *prm.SearchArtistParam) (*prm.SearchArtistResult, error) {
	
}

func (l *artistLogic) Upload(ctx context.Context, param *prm.UploadArtistParam) (*prm.UploadArtistResult, error) {
	artist := &mysql.Artist{
		Id:         util.IdGen.Generate(),
		Name:       param.Name,
		ArtistType: int32(param.ArtistType),
	}
	if err := l.ArtistDal.Save(artist); err != nil {
		return nil, err
	}
	return &prm.UploadArtistResult{
		ArtistId: artist.Id,
	}, nil
}
