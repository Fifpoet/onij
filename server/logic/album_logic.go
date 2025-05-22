package logic

import (
	"context"
	"errors"
	"onij/biz/biz"
	"onij/biz/getter"
	"onij/biz/prm"
	"onij/infra"
	"onij/infra/mysql"
	"onij/util"
	"onij/util/boost/collection/collext"
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
	albumPrime := &biz.AlbumPrime{
		Name:        param.Name,
		ArtistId:    param.ArtistId,
		CoverFileId: param.CoverFileId,
		IssueTime:   param.IssueTime,
		MusicId:     param.MusicId,
	}
	album := albumPrime.Model()
	err := l.AlbumDal.Save(album)
	if err != nil {
		return nil, err
	}
	return &prm.UploadAlbumResult{
		AlbumId: album.Id,
	}, nil
}

func (l *albumLogic) GetDetail(ctx context.Context, param *prm.GetAlbumDetailParam) (*prm.GetAlbumDetailResult, error) {
	mainAlbum, err := l.AlbumDal.GetById(param.AlbumId)
	if err != nil {
		return nil, err
	}
	relatedAlbums, err := l.AlbumDal.GetByRelatedId(mainAlbum.RelatedAlbumId)
	if err != nil {
		return nil, err
	}
	musics, err := l.MusicDal.GetByIds(collext.Pick(relatedAlbums, getter.AlbumMusicId)...)
	if err != nil {
		return nil, err
	}
	artists, err := l.ArtistDal.GetByIds(mainAlbum.ArtistId)
	if err != nil {
		return nil, err
	}
	if len(artists) == 0{
		return nil, errors.New("artists not found")
	}
	coverFiles, err := l.FileDal.GetByIds(mainAlbum.CoverFileId)
	if err != nil {
		return nil, err
	}
	if len(coverFiles) == 0{
		return nil, errors.New("cover file not found")
	}
	artGroup, tagGroup, err := l.getMusicArtistAndTags(musics...)
	if err != nil {
		return nil, err
	}
	return &prm.GetAlbumDetailResult{
		Album:           mainAlbum,
		Musics:          musics,
		CoverUrl:        util.DownloadFile(coverFiles[0].StoreKey),
		Artist:          artists[0],
		MusicArtistsMap: artGroup,
		MusicTagsMap:    tagGroup,
	}, nil
}

func (l *albumLogic) getMusicArtistAndTags(musics ...*mysql.Music) (map[int64][]*mysql.Artist, map[int64][]*mysql.Tag, error) {
	musicArtistsMap := collext.MapKV(musics, getter.MusicId, getter.MusicArtistIds)
	artIds := collext.PickCombine(musics, getter.MusicArtistIds)
	musIds := collext.Pick(musics, getter.MusicId)
	arts, err := l.ArtistDal.GetByIds(artIds...)
	if err != nil {
		return nil, nil, err
	}
	artMap := collext.Map(arts, getter.ArtistId)
	tags, err := l.TagDal.GetByResource(musIds...)
	if err != nil {
		return nil, nil, err
	}
	artGroup := make(map[int64][]*mysql.Artist)
	for mId, artIds := range musicArtistsMap {
		artGroup[mId] = collext.Pick(artIds, func(id int64) *mysql.Artist { return artMap[id] })
	}
	return collext.Group(arts, getter.ArtistId),
		collext.Group(tags, getter.TagResourceId), nil
}