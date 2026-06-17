package logic

import (
	"context"
	"errors"
	"fmt"
	"onij/biz/biz"
	"onij/biz/getter"
	"onij/biz/prm"
	"onij/infra"
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
	"onij/util/boost/exp"
	"onij/util/boost/tool"
	"time"
)

type MusicLogic interface {
	Upload(ctx context.Context, param *prm.UploadMusicParam) (*prm.UploadMusicResult, error)
	GetDetail(ctx context.Context, param *prm.GetMusicDetailParam) (*prm.GetMusicDetailResult, error)
	GetList(ctx context.Context, param *prm.GetMusicListParam) (*prm.GetMusicListResult, error)
	UpdateMv(ctx context.Context, param *prm.UpdateMusicMvParam) (*prm.UpdateMusicMvResult, error)
	GetMvBatch(ctx context.Context, param *prm.GetMusicMvBatchParam) (*prm.GetMusicMvBatchResult, error)
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
	music := &model.Music{
		Id:           util.IdGen.Generate(),
		Name:         param.Name,
		FullName:     param.Name + param.SubName,
		ArtistIds:    tool.ToJson(param.ArtistIds),
		ComposerIds:  tool.ToJson(param.ComposerIds),
		WriterIds:    tool.ToJson(param.WriterIds),
		IssueTime:    time.Unix(int64(exp.ValueOrZero(param.IssueTime)), 0),
		PerformType:  int32(param.PerformType),
		TimeLength:   param.TimeLength,
		MvUrl:        exp.ValueOrZero(param.MvUrl),
		AudioFileId:  param.AudioFileId,
		AudioQuality: int32(param.AudioQuality),
		LyricFileId:  param.LyricFileId,
		LyricContent: exp.ValueOrZero(param.LyricContent),
		RootId:       exp.ValueOrZero(param.RootMusicId),
		Priority:     exp.ValueOrZero(param.Priority),
	}
	if param.Id != nil {
		music.Id = *param.Id
	}

	_, err := l.MusicDal.Save(ctx, music)
	if err != nil {
		return nil, err
	}

	return &prm.UploadMusicResult{
		MusicId: music.Id,
	}, nil
}

func (l *musicLogic) GetDetail(ctx context.Context, param *prm.GetMusicDetailParam) (*prm.GetMusicDetailResult, error) {
	musics, err := l.MusicDal.GetByIds(ctx, param.Id)
	if err != nil {
		return nil, err
	}
	if len(musics) == 0 {
		return nil, errors.New("music not found")
	}
	music := musics[0]
	musicPrime := biz.MusicToBiz(music)

	// 获取艺术家信息
	artistIds := *tool.LoadJson[[]int64](music.ArtistIds, true)
	composerIds := *tool.LoadJson[[]int64](music.ComposerIds, true)
	writerIds := *tool.LoadJson[[]int64](music.WriterIds, true)

	allArtistIds := append(append(artistIds, composerIds...), writerIds...)
	artists, err := l.ArtistDal.GetByIds(ctx, allArtistIds...)
	if err != nil {
		return nil, err
	}
	musicPrime.Artists = collext.Pick(artists, biz.ArtistToBiz)

	// 获取专辑信息
	albumMusics, err := l.AlbumMusicDal.GetByMusicId(ctx, music.Id)
	if err != nil {
		return nil, err
	}
	albums, err := l.AlbumDal.GetByIds(ctx, collext.Pick(albumMusics, getter.AlbumMusicMusicId)...)
	if err != nil {
		return nil, err
	}
	musicPrime.Albums = collext.Pick(albums, biz.AlbumToBiz)

	// 获取文件信息
	files, err := l.FileDal.GetByIds(ctx, music.AudioFileId, music.LyricFileId)
	if err != nil {
		return nil, err
	}
	musicPrime.Files = collext.Pick(files, biz.FileToBiz)

	// 获取标签信息
	tags, err := l.TagDal.GetByResource(ctx, music.Id)
	if err != nil {
		return nil, err
	}
	musicPrime.Tags = collext.Pick(tags, biz.TagToBiz)

	return &prm.GetMusicDetailResult{
		Music: musicPrime,
	}, nil
}

func (l *musicLogic) GetList(ctx context.Context, param *prm.GetMusicListParam) (*prm.GetMusicListResult, error) {
	var musics []*model.Music
	var err error
	musics, err = l.MusicDal.Search(
		ctx, param.ArtistIds, param.TagTypes,
		param.Keyword, util.Page{Page: param.Page, Limit: param.Limit},
	)
	if err != nil {
		return nil, err
	}

	artGroup, tagGroup, err := l.getMusicArtistAndTags(ctx, musics...)
	if err != nil {
		return nil, err
	}
	return &prm.GetMusicListResult{
		Musics:          musics,
		MusicArtistsMap: artGroup,
		MusicTagsMap:    tagGroup,
	}, nil
}

func (l *musicLogic) getMusicArtistAndTags(ctx context.Context, musics ...*model.Music) (map[int64][]*model.Artist, map[int64][]*model.Tag, error) {
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

func (l *musicLogic) UpdateMv(ctx context.Context, param *prm.UpdateMusicMvParam) (*prm.UpdateMusicMvResult, error) {
	var music *model.Music
	var err error

	if param.MusicId > 0 {
		music, err = l.musicByID(ctx, param.MusicId)
		if err != nil {
			return nil, err
		}
	} else if param.ThirdId > 0 {
		musics, err := l.MusicDal.GetByThirdIds(ctx, param.ThirdId)
		if err != nil {
			return nil, err
		}
		if len(musics) > 0 {
			music = musics[0]
		} else {
			music = &model.Music{
				Id:        util.IdGen.Generate(),
				Name:      fmt.Sprintf("#%d", param.ThirdId),
				FullName:  fmt.Sprintf("#%d", param.ThirdId),
				ArtistIds: "[]",
				ThirdId:   param.ThirdId,
				MvUrl:     param.MvUrl,
			}
			if err := l.MusicDal.InsertMvStub(ctx, music); err != nil {
				return nil, err
			}
			return &prm.UpdateMusicMvResult{MusicId: music.Id}, nil
		}
	} else {
		return nil, errors.New("music_id or third_id required")
	}

	ok, err := l.MusicDal.UpdateById(ctx, func(u model.MusicUpdater) {
		u.MvUrl(param.MvUrl)
	}, music.Id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("music update failed")
	}
	return &prm.UpdateMusicMvResult{MusicId: music.Id}, nil
}

func (l *musicLogic) GetMvBatch(ctx context.Context, param *prm.GetMusicMvBatchParam) (*prm.GetMusicMvBatchResult, error) {
	musics, err := l.MusicDal.GetByThirdIds(ctx, param.ThirdIds...)
	if err != nil {
		return nil, err
	}
	byThird := collext.Map(musics, getter.MusicThirdId)
	items := make([]*api.MusicMvItem, 0, len(param.ThirdIds))
	for _, tid := range param.ThirdIds {
		item := &api.MusicMvItem{ThirdId: tid}
		if m := byThird[tid]; m != nil {
			item.MusicId = m.Id
			item.MvUrl = m.MvUrl
		}
		items = append(items, item)
	}
	return &prm.GetMusicMvBatchResult{Items: items}, nil
}

func (l *musicLogic) musicByID(ctx context.Context, id int64) (*model.Music, error) {
	musics, err := l.MusicDal.GetByIds(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(musics) == 0 {
		return nil, errors.New("music not found")
	}
	return musics[0], nil
}
