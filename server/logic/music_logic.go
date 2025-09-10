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
	// 获取艺术家信息（用于验证艺术家是否存在）
	_, err := l.ArtistDal.GetByIds(ctx, param.ArtistIds...)
	if err != nil {
		return nil, err
	}

	// 构建音乐对象
	music := &model.Music{
		Id:           util.IdGen.Generate(),
		Name:         param.Name,
		FullName:     param.Name, // TODO: 根据艺术家名称构建全名
		ArtistIds:    tool.ToJson(param.ArtistIds),
		ComposerIds:  tool.ToJson(param.ComposerIds),
		WriterIds:    tool.ToJson(param.WriterIds),
		IssueTime:    time.Now(), // TODO: 根据 param.IssueTime 设置
		PerformType:  0,          // TODO: 根据实际需求设置
		TimeLength:   0,          // TODO: 根据实际需求设置
		MvUrl:        "",         // TODO: 根据 param.MvUrl 设置
		AudioFileId:  param.AudioFileId,
		AudioQuality: 0, // TODO: 根据实际需求设置
		LyricFileId:  param.LyricFileId,
		LyricContent: "", // TODO: 根据实际需求设置
		RootId:       0,  // TODO: 根据 param.RootMusicId 设置
		Priority:     0,  // TODO: 根据实际需求设置
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 如果传入了 ID，则使用传入的 ID
	if param.Id != nil {
		music.Id = *param.Id
	}

	_, err = l.MusicDal.Save(ctx, music)
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

	// 获取艺术家信息
	artistIds := *tool.LoadJson[[]int64](music.ArtistIds, true)
	composerIds := *tool.LoadJson[[]int64](music.ComposerIds, true)
	writerIds := *tool.LoadJson[[]int64](music.WriterIds, true)

	allArtistIds := append(append(artistIds, composerIds...), writerIds...)
	artists, err := l.ArtistDal.GetByIds(ctx, allArtistIds...)
	if err != nil {
		return nil, err
	}

	// TODO: 需要实现 GetByMusicId 方法
	// albums, err := l.AlbumDal.GetByMusicId(ctx, music.Id)
	// if err != nil {
	// 	return nil, err
	// }

	// 获取文件信息
	files, err := l.FileDal.GetByIds(ctx, music.AudioFileId, music.LyricFileId)
	if err != nil {
		return nil, err
	}

	// 获取标签信息
	tags, err := l.TagDal.GetByResource(ctx, music.Id)
	if err != nil {
		return nil, err
	}

	fileMap := collext.Map(files, getter.FileId)
	return &prm.GetMusicDetailResult{
		Music:      music,
		Tags:       tags,
		ArtistMap:  collext.Map(artists, getter.ArtistId),
		ArtistIds:  artistIds,
		ComposerId: 0,                // TODO: 从 composerIds 中获取
		WriterId:   0,                // TODO: 从 writerIds 中获取
		Albums:     []*model.Album{}, // TODO: 获取实际专辑数据
		FileMap:    fileMap,
	}, nil
}

func (l *musicLogic) GetList(ctx context.Context, param *prm.GetMusicListParam) (*prm.GetMusicListResult, error) {
	var musics []*model.Music
	var err error

	// TODO: 需要实现 AlbumId 相关的逻辑
	// if param.AlbumId != nil {
	// 	// 专辑范围内
	// 	album, err := l.AlbumDal.GetById(ctx, *param.AlbumId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	relatedAlbum, err := l.AlbumDal.GetByRelatedId(ctx, album.Id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	musics, err = l.MusicDal.GetByIds(ctx, collext.Pick(relatedAlbum, getter.AlbumMusicId)...)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if len(param.Keywords) > 0 {
	// 		musics = collext.Select(musics, func(m *model.Music) (*model.Music, bool) {
	// 			for _, kw := range param.Keywords {
	// 				if strings.Contains(m.Name, kw) {
	// 					return m, true
	// 				}
	// 			}
	// 			return nil, false
	// 		})
	// 	}
	// } else {
	// 不指定专辑，使用搜索功能
	musics, err = l.MusicDal.SearchByArtistAndNameAndTag(
		ctx, param.ArtistIds, param.WriterIds, param.ComposerIds, param.TagTypes,
		param.Keywords, util.Page{Page: param.Page, Limit: param.Limit},
	)
	if err != nil {
		return nil, err
	}
	// }

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
