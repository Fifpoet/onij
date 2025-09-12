package logic

import (
	"context"
	"errors"
	"onij/biz/biz"
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
	GetList(ctx context.Context, param *prm.GetAlbumListParam) (*prm.GetAlbumListResult, error)
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
		Name:        param.Name,
		Profile:     param.Profile,
		AlbumType:   int32(param.AlbumType),
		ArtistIds:   tool.ToJson(param.ArtistIds),
		IssueTime:   time.Unix(int64(param.IssueTime), 0),
		CoverFileId: param.CoverFileId,
		LiveUrl:     param.LiveUrl,
	}
	if param.AlbumId == nil {
		album.Id = util.IdGen.Generate()
	} else {
		album.Id = *param.AlbumId
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
	albums, err := l.AlbumDal.GetByIds(ctx, param.AlbumId)
	if err != nil {
		return nil, err
	}
	if len(albums) == 0 {
		return nil, errors.New("album not found")
	}
	album := albums[0]
	albumPrime := biz.AlbumToBiz(album)

	albumMusics, err := l.AlbumMusicDal.GetByAlbumId(ctx, param.AlbumId)
	if err != nil {
		return nil, err
	}
	albumPrime.AlbumMusics = collext.Pick(albumMusics, biz.AlbumMusicToBiz)

	// 获取艺术家信息
	artistIds := *tool.LoadJson[[]int64](album.ArtistIds, true)
	artists, err := l.ArtistDal.GetByIds(ctx, artistIds...)
	if err != nil {
		return nil, err
	}
	if len(artists) == 0 {
		return nil, errors.New("artists not found")
	}
	albumPrime.Artists = collext.Pick(artists, biz.ArtistToBiz)

	// 获取封面文件
	coverFiles, err := l.FileDal.GetByIds(ctx, album.CoverFileId)
	if err != nil {
		return nil, err
	}
	if len(coverFiles) == 0 {
		return nil, errors.New("cover file not found")
	}
	albumPrime.CoverFile = biz.FileToBiz(coverFiles[0])

	tags, err := l.TagDal.GetByResource(ctx, append(collext.Pick(albumMusics, getter.AlbumMusicMusicId), param.AlbumId)...)
	if err != nil {
		return nil, err
	}
	albumPrime.TagPrime = collext.Pick(tags, biz.TagToBiz)

	return &prm.GetAlbumDetailResult{
		Album: albumPrime,
	}, nil
}

func (l *albumLogic) GetList(ctx context.Context, param *prm.GetAlbumListParam) (*prm.GetAlbumListResult, error) {
	albums, cnt, err := l.AlbumDal.GetByKeywordAndArtistId(ctx, param.Keyword, param.ArtistId, util.Page{
		Page: param.Page,
		Limit: param.Limit,
	})
	if err != nil {
		return nil, err
	}

	files , err  := l.FileDal.GetByIds(ctx, collext.Pick(albums, getter.AlbumCoverFileId)...)
	if err != nil {
		return nil, err
	}
	fileMap := collext.Map(files, getter.FileId)

	albumMusics, err := l.AlbumMusicDal.GetByAlbumIds(ctx, collext.Pick(albums, getter.AlbumId)...)
	if err != nil {
		return nil, err
	}
	artistIds := collext.Distinct(append(collext.PickCombine(albumMusics, getter.AlbumMusicArtistIds), 
	collext.PickCombine(albums, getter.AlbumArtistIds)...))

	artists, err := l.ArtistDal.GetByIds(ctx, artistIds...)
	if err != nil {
		return nil, err
	}

	albumPrimes := collext.Pick(albums, func(a *model.Album) *biz.AlbumPrime {
		prime := biz.AlbumToBiz(a)
		prime.CoverFile = biz.FileToBiz(fileMap[a.CoverFileId])
		prime.Artists = collext.Pick(artists, biz.ArtistToBiz)
		return prime
	})
	
	return &prm.GetAlbumListResult{
		Albums: albumPrimes,
		Total:  cnt,
	}, nil
}