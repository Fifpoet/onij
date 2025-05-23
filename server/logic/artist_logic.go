package logic

import (
	"context"
	"onij/biz/prm"
	"onij/infra"
	"onij/infra/mysql"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
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
	var arts []*mysql.Artist
	var err error
	if param.TagType != nil || param.TagGroup != nil {
		var tags []*mysql.Tag
		tags, err = l.TagDal.GetByGroupType(param.TagGroup, param.TagType)
		if err != nil {
			return nil, err
		}
		artistIds := collext.Select(tags, func(tag *mysql.Tag) (int64, bool) {
			return tag.ResourceId, tag.ResourceType == int32(api.ResourceType_RT_Artist)
		})
		arts, err = l.ArtistDal.GetByIds(artistIds...)
		if err != nil {
			return nil, err
		}
	} else {
		arts, err = l.ArtistDal.GetLikeNameAndType(param.Keyword, nil)
		if err != nil {
			return nil, err
		}
	}
	return &prm.SearchArtistResult{
		Artists: arts,
	}, nil
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
