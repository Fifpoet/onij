package logic

import (
	"context"
	"onij/biz/biz"
	"onij/biz/prm"
	"onij/infra"
	"onij/util"
)

type MusicLogic interface {
	Upload(ctx context.Context, param *prm.UploadMusicParam) (*prm.UploadMusicResult, error)
	GetDetail(ctx context.Context, param *prm.GetMusicDetailParam) (*prm.GetMusicDetailResult, error)
}

type musicLogic struct {
	*infra.AllInfra
}

func NewMusicLogic(i *infra.AllInfra) MusicLogic {
	return &musicLogic{
		AllInfra: i,
	}
}

func (l *musicLogic) Upload(ctx context.Context, param *prm.UploadMusicParam) (*prm.UploadMusicResult, error) {
	musicPrime := &biz.MusicPrime{
		Id:           param.Id,
		Name:         param.Name,
		ArtistIds:    param.ArtistIds,
		Mp3FileId:    param.Mp3FileId,
		LyricsFileId: param.LyricsFileId,
		ComposerId:   param.ComposerId,
		WriterId:     param.WriterId,
		AlbumId:      param.AlbumId,
		MvUrl:        param.MvUrl,
		RootMusicId:  param.RootMusicId,
		IssueTime:    param.IssueTime,
	}
	err := l.MusicDal.Upsert(musicPrime.Upsert())
	if err != nil {
		return nil, err
	}
	return &prm.UploadMusicResult{}, nil
}

func (l *musicLogic) GetDetail(ctx context.Context, param *prm.GetMusicDetailParam) (*prm.GetMusicDetailResult, error) {
	musics, err := l.MusicDal.GetByIds(param.Id)
	if err != nil {
		return nil, err
	}
	music := musics[0]
	artists, err := l.ArtistDal.GetByIds(append([]int64{music.ComposerId, music.WriterId}, util.StrList2Int64(music.ArtistIds)...)...)
	if err != nil {
		return nil, err
	}
	files, err := l.FileDal.GetByIds(append([]int64{music.Mp3FileId}, music.LyricFileId)...)
	if err != nil {
		return nil, err
	}
	album, err := l.AlbumDal.GetById(music.AlbumId)
	if err != nil {
		return nil, err
	}
	return &prm.GetMusicDetailResult{
		Music:         music,
		Artists:       artists[2:],
		Composer:      artists[0],
		Writer:        artists[1],
		Mp3FileUrl:    util.DownloadFile(files[0].StoreKey),
		LyricsFileUrl: util.DownloadFile(files[1].StoreKey),
		Album:         album,
	}, nil
}
