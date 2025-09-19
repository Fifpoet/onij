package logic

import (
	"context"
	"onij/biz/prm"
	"onij/infra"
	"onij/model"
	"onij/util"
	"onij/util/boost/exp"
)

type ArtistLogic interface {
	Upload(ctx context.Context, param *prm.UploadArtistParam) (*prm.UploadArtistResult, error)
	GetList(ctx context.Context, param *prm.GetArtistListParam) (*prm.GetArtistListResult, error)
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
	artist := &model.Artist{
		Name:         param.Name,
		AvatarFileId: exp.ValueOrZero(param.AvatarFileId),
	}
	artist.Id = util.IdGen.Generate()
	if param.AvatarFileId != nil {
		artist.AvatarFileId = *param.AvatarFileId
	}

	_, err := l.ArtistDal.Save(ctx, artist)
	if err != nil {
		return nil, err
	}

	return &prm.UploadArtistResult{
		ArtistId: artist.Id,
	}, nil
}

func (l *artistLogic) GetList(ctx context.Context, param *prm.GetArtistListParam) (*prm.GetArtistListResult, error) {
	artists, cnt, err := l.ArtistDal.GetByName(ctx, param.Name, param.Keyword, util.Page{
		Page:  param.Page,
		Limit: param.Limit,
	})
	if err != nil {
		return nil, err
	}

	return &prm.GetArtistListResult{
		Artists: artists,
		Total:   cnt,
	}, nil
}
