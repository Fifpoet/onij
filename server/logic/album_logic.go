package logic

import (
	"context"
	"onij/biz/biz"
	"onij/biz/prm"
	"onij/infra"
)

type AlbumLogic interface {
	Upload(ctx context.Context, param *prm.UploadAlbumParam) (*prm.UploadAlbumResult, error)
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
	err := l.AlbumDal.Save(albumPrime.Model())
	if err != nil {
		return nil, err
	}
	return &prm.UploadAlbumResult{}, nil
}
