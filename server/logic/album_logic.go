package logic

import (
	"context"
	"errors"
	"onij/biz/getter"
	"onij/biz/prm"
	"onij/infra"
	"onij/model"
	"onij/util"
	"onij/util/boost/collection/collext"
	"onij/util/boost/tool"
	"time"
)

type AlbumLogic interface {
	Upload(ctx context.Context, param *prm.UploadAlbumParam) (*prm.UploadAlbumResult, error)
	GetDetail(ctx context.Context, param *prm.GetAlbumDetailParam) (*prm.GetAlbumDetailResult, error)
}

type albumLogic struct {
	*infra.AllInfra
}

func NewAlbumLogic(i *infra.AllInfra) AlbumLogic {
	return &albumLogic{
		AllInfra: i,
	}
}

func (l *albumLogic) Upload(ctx context.Context, param *prm.UploadAlbumParam) (*prm.UploadAlbumResult, error) {
	album := &model.Album{
		Id:          util.IdGen.Generate(),
		Name:        param.Name,
		Profile:     param.Profile,
		AlbumType:   int32(param.AlbumType),
		ArtistIds:   tool.ToJson(param.ArtistIds),
		IssueTime:   time.Unix(int64(param.IssueTime), 0),
		CoverFileId: param.CoverFileId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := l.AlbumDal.Save(ctx, album)
	if err != nil {
		return nil, err
	}

	return &prm.UploadAlbumResult{
		AlbumId: album.Id,
	}, nil
}

func (l *albumLogic) GetDetail(ctx context.Context, param *prm.GetAlbumDetailParam) (*prm.GetAlbumDetailResult, error) {
	mainAlbum, err := l.AlbumDal.GetById(ctx, param.AlbumId)
	if err != nil {
		return nil, err
	}

	// 获取专辑相关的音乐
	// TODO: 需要实现 GetByAlbumId 方法或者通过其他方式获取专辑音乐
	// relatedAlbums, err := l.AlbumDal.GetByRelatedId(mainAlbum.RelatedAlbumId)
	// if err != nil {
	// 	return nil, err
	// }
	// musics, err := l.MusicDal.GetByIds(ctx, collext.Pick(relatedAlbums, getter.AlbumMusicId)...)
	// if err != nil {
	// 	return nil, err
	// }

	// 获取艺术家信息
	artistIds := *tool.LoadJson[[]int64](mainAlbum.ArtistIds, true)
	artists, err := l.ArtistDal.GetByIds(ctx, artistIds...)
	if err != nil {
		return nil, err
	}
	if len(artists) == 0 {
		return nil, errors.New("artists not found")
	}

	// 获取封面文件
	coverFiles, err := l.FileDal.GetByIds(ctx, mainAlbum.CoverFileId)
	if err != nil {
		return nil, err
	}
	if len(coverFiles) == 0 {
		return nil, errors.New("cover file not found")
	}

	// TODO: 需要实现音乐和标签的获取
	// artGroup, tagGroup, err := l.getMusicArtistAndTags(musics...)
	// if err != nil {
	// 	return nil, err
	// }

	return &prm.GetAlbumDetailResult{
		Album:           mainAlbum,
		Musics:          []*model.Music{},      // TODO: 获取实际音乐列表
		AlbumMusics:     []*model.AlbumMusic{}, // TODO: 获取专辑音乐关联
		CoverUrl:        util.DownloadFile(coverFiles[0].StoreKey),
		Artists:         artists,
		MusicArtistsMap: make(map[int64][]*model.Artist), // TODO: 获取实际数据
		MusicTagsMap:    make(map[int64][]*model.Tag),    // TODO: 获取实际数据
	}, nil
}

func (l *albumLogic) getMusicArtistAndTags(ctx context.Context, musics ...*model.Music) (map[int64][]*model.Artist, map[int64][]*model.Tag, error) {
	if len(musics) == 0 {
		return make(map[int64][]*model.Artist), make(map[int64][]*model.Tag), nil
	}

	musicArtistsMap := collext.MapKV(musics, getter.MusicId, getter.MusicArtistIds)
	artIds := collext.PickCombine(musics, getter.MusicArtistIds)
	musIds := collext.Pick(musics, getter.MusicId)

	arts, err := l.ArtistDal.GetByIds(ctx, artIds...)
	if err != nil {
		return nil, nil, err
	}
	artMap := collext.Map(arts, getter.ArtistId)

	tags, err := l.TagDal.GetByResource(ctx, musIds...)
	if err != nil {
		return nil, nil, err
	}

	artGroup := make(map[int64][]*model.Artist)
	for mId, artIds := range musicArtistsMap {
		artGroup[mId] = collext.Pick(artIds, func(id int64) *model.Artist { return artMap[id] })
	}

	return artGroup, collext.Group(tags, getter.TagResourceId), nil
}
