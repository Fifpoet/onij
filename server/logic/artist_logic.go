package logic

import (
	"context"
	"gorm.io/gorm"
	"onij/biz/prm"
	"onij/infra"
	"onij/model"
	"onij/util"
	"time"
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

//func (l *artistLogic) Search(ctx context.Context, param *prm.SearchArtistParam) (*prm.SearchArtistResult, error) {
//	var arts []*model.Artist
//	var err error
//	if param.TagType != nil || param.TagGroup != nil {
//		var tags []*model.Tag
//		tags, err = l.TagDal.GetByGroupType(param.TagGroup, param.TagType)
//		if err != nil {
//			return nil, err
//		}
//		artistIds := collext.Select(tags, func(tag *model.Tag) (int64, bool) {
//			return tag.ResourceId, tag.ResourceType == int32(api.ResourceType_RT_Artist)
//		})
//		arts, err = l.ArtistDal.GetByIds(artistIds...)
//		if err != nil {
//			return nil, err
//		}
//	} else {
//		arts, err = l.ArtistDal.GetLikeNameAndType(param.Keyword, nil)
//		if err != nil {
//			return nil, err
//		}
//	}
//	return &prm.SearchArtistResult{
//		Artists: arts,
//	}, nil
//}

func (l *artistLogic) Upload(ctx context.Context, param *prm.UploadArtistParam) (*prm.UploadArtistResult, error) {
	artist := &model.Artist{
		Id:           util.IdGen.Generate(),
		Name:         param.Name,
		AvatarFileId: param.AvatarFileId,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
		DeletedAt:    gorm.DeletedAt{},
	}
	if _, err := l.ArtistDal.Save(ctx, artist); err != nil {
		return nil, err
	}
	return &prm.UploadArtistResult{
		ArtistId: artist.Id,
	}, nil
}
