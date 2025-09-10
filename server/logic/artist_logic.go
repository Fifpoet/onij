package logic

import (
	"context"
	"onij/biz/prm"
	"onij/infra"
	"onij/model"
	"onij/util"
	"time"
)

type ArtistLogic interface {
	Upload(ctx context.Context, param *prm.UploadArtistParam) (*prm.UploadArtistResult, error)
	Search(ctx context.Context, param *prm.SearchArtistParam) (*prm.SearchArtistResult, error)
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
		Id:           util.IdGen.Generate(),
		Name:         param.Name,
		AvatarFileId: 0, // TODO: 根据 param.AvatarFileId 处理
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := l.ArtistDal.Save(ctx, artist)
	if err != nil {
		return nil, err
	}

	return &prm.UploadArtistResult{
		ArtistId: artist.Id,
	}, nil
}

func (l *artistLogic) Search(ctx context.Context, param *prm.SearchArtistParam) (*prm.SearchArtistResult, error) {
	var arts []*model.Artist
	var err error

	if param.TagType != nil || param.TagGroup != nil {
		// 通过标签搜索
		tags, err := l.TagDal.GetByGroupType(ctx, param.TagGroup, param.TagType)
		if err != nil {
			return nil, err
		}

		// 过滤出艺术家相关的标签
		var artistIds []int64
		for _, tag := range tags {
			// TODO: 需要根据实际的资源类型常量来判断
			// if tag.ResourceType == int32(api.ResourceType_RT_Artist) {
			artistIds = append(artistIds, tag.ResourceId)
			// }
		}

		if len(artistIds) > 0 {
			arts, err = l.ArtistDal.GetByIds(ctx, artistIds...)
			if err != nil {
				return nil, err
			}
		}
	} else {
		// 通过关键词搜索
		arts, _, err = l.ArtistDal.GetByKeyword(ctx, param.Keyword, 0, param.Limit)
		if err != nil {
			return nil, err
		}
	}

	return &prm.SearchArtistResult{
		Artists: arts,
	}, nil
}
