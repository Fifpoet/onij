package logic

import (
	"context"
	"errors"
	"onij/biz/biz"
	"onij/biz/getter"
	"onij/biz/prm"
	"onij/infra"
	"onij/util"
	"onij/util/boost/collection/collext"
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
		ArtistNames:  collext.Pick(artistNames, getter.ArtistName),
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
	albums, err := l.AlbumDal.GetByMusicId(music.Id)
	if err != nil {
		return nil, err
	}
	files, err := l.FileDal.GetByIds(append(
		collext.Pick(albums, getter.AlbumCoverFileId),
		[]int64{music.Mp3FileId, music.LyricFileId}...)...,
	)
	if err != nil {
		return nil, err
	}
	tags, err := l.TagDal.GetByResource(append(collext.Pick(artists, getter.ArtistId), music.Id)...)
	if err != nil {
		return nil, err
	}
	fileMap := collext.Map(files, getter.FileId)
	return &prm.GetMusicDetailResult{
		Music:      music,
		Tags:       tags,
		ArtistMap:  collext.Map(artists, getter.ArtistId),
		ArtistIds:  util.StrList2Int64(music.ArtistIds),
		ComposerId: music.ComposerId,
		WriterId:   music.WriterId,
		Albums:     albums,
		FileMap:    fileMap,
	}, nil
}

func (l *musicLogic) GetList(ctx context.Context, param *prm.GetMusicListParam) (*prm.GetMusicListResult, error) {
	var musics []*model.Music
	var err error
	if param.AlbumId != nil {
		// 专辑范围内
		album, err := l.AlbumDal.GetById(*param.AlbumId)
		if err != nil {
			return nil, err
		}
		relatedAlbum, err := l.AlbumDal.GetByRelatedId(album.Id)
		if err != nil {
			return nil, err
		}
		musics, err = l.MusicDal.GetByIds(collext.Pick(relatedAlbum, getter.AlbumMusicId)...)
		if err != nil {
			return nil, err
		}
		if len(param.Keywords) > 0 {
			musics = collext.Select(musics, func(m *model.Music) (*model.Music, bool) {
				for _, kw := range param.Keywords {
					if strings.Contains(m.Name, kw) {
						return m, true
					}
				}
				return nil, false
			})
		}
	} else {
		// 不指定专辑
		musics, err = l.MusicDal.SearchByArtistAndNameAndTag(
			param.ArtistIds, param.WriterIds, param.ComposerIds, param.TagTypes,
			param.Keywords, util.Page{Page: param.Page, Limit: param.Limit},
		)
		if err != nil {
			return nil, err
		}
	}

	artGroup, tagGroup, err := l.getMusicArtistAndTags(musics...)
	if err != nil {
		return nil, err
	}
	return &prm.GetMusicListResult{
		Musics:          musics,
		MusicArtistsMap: artGroup,
		MusicTagsMap:    tagGroup,
	}, nil
}

func (l *musicLogic) getMusicArtistAndTags(musics ...*model.Music) (map[int64][]*model.Artist, map[int64][]*model.Tag, error) {
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
	artGroup := make(map[int64][]*model.Artist)
	for mId, artIds := range musicArtistsMap {
		artGroup[mId] = collext.Pick(artIds, func(id int64) *model.Artist { return artMap[id] })
	}
	return artGroup, collext.Group(tags, getter.TagResourceId), nil
}
