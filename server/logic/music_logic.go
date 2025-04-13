package logic

import (
	"context"
	"errors"
	"onij/biz/biz"
	"onij/biz/prm"
	"onij/infra"
	"onij/infra/mysql"
	"onij/util"
	"onij/util/boost/collection/collext"
	"onij/util/boost/exp"
	"strings"
)

type MusicLogic interface {
	Upload(ctx context.Context, param *prm.UploadMusicParam) (*prm.UploadMusicResult, error)
	GetDetail(ctx context.Context, param *prm.GetMusicDetailParam) (*prm.GetMusicDetailResult, error)
	GetList(ctx context.Context, param *prm.GetMusicListParam) (*prm.GetMusicListResult, error)
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
	artistNames, err := l.ArtistDal.GetByIds(param.ArtistIds...)
	if err != nil {
		return nil, err
	}
	musicPrime := &biz.MusicPrime{
		Id:           param.Id,
		Name:         param.Name,
		ArtistIds:    param.ArtistIds,
		ArtistNames:  collext.Pick(artistNames, func(a *mysql.Artist) string { return a.Name }),
		Mp3FileId:    param.Mp3FileId,
		LyricsFileId: param.LyricsFileId,
		ComposerId:   param.ComposerId,
		WriterId:     param.WriterId,
		AlbumId:      param.AlbumId,
		MvUrl:        param.MvUrl,
		RootMusicId:  param.RootMusicId,
		IssueTime:    param.IssueTime,
	}
	music, toUpdate := musicPrime.Upsert()
	err = l.MusicDal.Upsert(music, toUpdate)
	if err != nil {
		return nil, err
	}
	return &prm.UploadMusicResult{
		MusicId: music.Id,
	}, nil
}

func (l *musicLogic) GetDetail(ctx context.Context, param *prm.GetMusicDetailParam) (*prm.GetMusicDetailResult, error) {
	musics, err := l.MusicDal.GetByIds(param.Id)
	if err != nil {
		return nil, err
	}
	if len(musics) == 0 {
		return nil, errors.New("music not found")
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
	albums, err := l.AlbumDal.GetByMusicId(music.Id)
	if err != nil {
		return nil, err
	}
	return &prm.GetMusicDetailResult{
		Music:         music,
		ArtistMap:     collext.Map(artists, func(art *mysql.Artist) int64 { return art.Id }),
		ArtistIds:     util.StrList2Int64(music.ArtistIds),
		ComposerId:    music.ComposerId,
		WriterId:      music.WriterId,
		Mp3FileUrl:    util.DownloadFile(files[0].StoreKey),
		LyricsFileUrl: util.DownloadFile(files[1].StoreKey),
		Albums:        albums,
	}, nil
}

func (l *musicLogic) GetList(ctx context.Context, param *prm.GetMusicListParam) (*prm.GetMusicListResult, error) {
	// 专辑范围内
	if param.AlbumId != nil {
		album, err := l.AlbumDal.GetById(*param.AlbumId)
		if err != nil {
			return nil, err
		}
		relatedAlbum, err := l.AlbumDal.GetByRelatedId(album.Id)
		if err != nil {
			return nil, err
		}
		musics, err := l.MusicDal.GetByIds(collext.Pick(relatedAlbum, func(a *mysql.Album) int64 {
			return a.MusicId
		})...)
		if err != nil {
			return nil, err
		}
		if param.Keyword != nil {
			return &prm.GetMusicListResult{Musics: musics}, nil
		}
		return &prm.GetMusicListResult{
			Musics: collext.Select(musics, func(m *mysql.Music) (*mysql.Music, bool) {
				return m, strings.Contains(m.FullName, *param.Keyword)
			}),
		}, nil
	}

	// 不指定专辑
	musics, err := l.MusicDal.GetByArtistAndName(
		exp.ValueOrZero(param.ArtistId), exp.ValueOrZero(param.Keyword),
		util.Page{Page: param.Page, Limit: param.Limit},
	)
	if err != nil {
		return nil, err
	}
	return &prm.GetMusicListResult{
		Musics: musics,
	}, nil
}
