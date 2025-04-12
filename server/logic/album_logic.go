package logic

import (
	"context"
	"onij/biz/biz"
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
	musics, err := l.MusicDal.GetByIds(collext.Pick(relatedAlbums, func(a *mysql.Album) int64 {
		return a.MusicId
	})...)
	if err != nil {
		return nil, err
	}
	artists, err := l.ArtistDal.GetByIds(mainAlbum.ArtistId)
	if err != nil {
		return nil, err
	}
	coverFiles, err := l.FileDal.GetByIds(mainAlbum.CoverFileId)
	if err != nil {
		return nil, err
	}
	return &prm.GetAlbumDetailResult{
		Album:    mainAlbum,
		Musics:   musics,
		CoverUrl: util.DownloadFile(coverFiles[0].StoreKey),
		Artist:   artists[0],
	}, nil
}
