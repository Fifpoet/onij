package prm

import (
	"onij/biz/getter"
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
)

type UploadMusicParam struct {
	Id          *int64
	Name        string
	ArtistIds   []int64
	AudioFileId int64
	LyricFileId int64
	ComposerIds []int64
	WriterIds   []int64
	AlbumId     *int64
	MvUrl       *string
	RootMusicId *int64
	IssueTime   *int32
}

func NewUploadMusicParam(req *api.UploadMusicReq) *UploadMusicParam {
	return &UploadMusicParam{
		Id:          req.MusicId,
		Name:        req.Name,
		ArtistIds:   req.ArtistIds,
		AudioFileId: req.AudioFileId,
		LyricFileId: req.LyricFileId,
		ComposerIds: req.ComposerIds,
		WriterIds:   req.WriterIds,
		AlbumId:     req.AlbumId,
		MvUrl:       req.MvUrl,
		RootMusicId: req.RootMusicId,
		IssueTime:   req.IssueTime,
	}
}

type UploadMusicResult struct {
	MusicId int64
}

func (p *UploadMusicResult) Resp() *api.UploadMusicResp {
	return &api.UploadMusicResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		MusicId: p.MusicId,
	}
}

type GetMusicDetailParam struct {
	Id int64
}

func NewGetMusicDetailParam(req *api.GetMusicDetailReq) *GetMusicDetailParam {
	return &GetMusicDetailParam{
		Id: req.MusicId,
	}
}

type GetMusicDetailResult struct {
	Music      *model.Music
	Tags       []*model.Tag
	ArtistMap  map[int64]*model.Artist
	ArtistIds  []int64
	ComposerId int64
	WriterId   int64
	Albums     []*model.Album
	FileMap    map[int64]*model.File
}

func (p *GetMusicDetailResult) Resp() *api.GetMusicDetailResp {
	album := &model.Album{}
	if len(p.Albums) > 0 {
		album = p.Albums[0]
	}
	var mp3Url, lyrUrl, abmUrl string
	if p.FileMap[p.Music.AudioFileId] != nil {
		mp3Url = p.FileMap[p.Music.AudioFileId].StoreKey
	}
	if p.FileMap[p.Music.LyricFileId] != nil {
		lyrUrl = p.FileMap[p.Music.LyricFileId].StoreKey
	}
	if p.FileMap[album.CoverFileId] != nil {
		abmUrl = p.FileMap[album.CoverFileId].StoreKey
	}
	return &api.GetMusicDetailResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Music: &api.Music{
			Id:            p.Music.Id,
			Name:          p.Music.Name,
			ArtistIds:     collext.Pick(p.ArtistIds, func(id int64) int64 { return p.ArtistMap[id].Id }),
			ArtistNames:   collext.Pick(p.ArtistIds, func(id int64) string { return p.ArtistMap[id].Name }),
			ComposerId:    p.ComposerId,
			ComposerName:  p.ArtistMap[p.ComposerId].Name,
			WriterId:      p.WriterId,
			WriterName:    p.ArtistMap[p.WriterId].Name,
			IssueTime:     int32(p.Music.IssueTime.Unix()),
			MvUrl:         p.Music.MvUrl,
			Mp3FileUrl:    util.DownloadFile(mp3Url),
			LyricsFileUrl: util.DownloadFile(lyrUrl),
			Album: &api.Album{
				Id:           album.Id,
				Name:         album.Name,
				CoverFileUrl: util.DownloadFile(abmUrl),
			},
			Tags: collext.Pick(p.Tags, func(tag *model.Tag) *api.Tag {
				return &api.Tag{
					ResourceId:   tag.ResourceId,
					ResourceType: api.ResourceType(tag.ResourceType),
					TagBiz:       api.TagBiz(tag.TagBiz),
					TagGroup:     api.TagGroup(tag.TagGroup),
					TagType:      api.TagType(tag.TagType),
					TargetId:     tag.TargetId,
					TargetType:   tag.TargetType,
					Extra:        tag.Extra,
				}
			}),
		},
	}
}

type GetMusicListParam struct {
	Keywords    []string
	ArtistIds   []int64
	WriterIds   []int64
	ComposerIds []int64
	TagTypes    []int32
	Page        int32
	Limit       int32
}

func NewGetMusicListParam(req *api.GetMusicListReq) *GetMusicListParam {
	return &GetMusicListParam{
		Keywords:    req.Keywords,
		ArtistIds:   req.ArtistIds,
		WriterIds:   req.WriterIds,
		ComposerIds: req.ComposerIds,
		TagTypes:    collext.Pick(req.TagTypes, func(tt api.TagType) int32 { return int32(tt) }),
		Page:        req.Page,
		Limit:       req.Limit,
	}
}

type GetMusicListResult struct {
	Musics          []*model.Music
	MusicArtistsMap map[int64][]*model.Artist
	MusicTagsMap    map[int64][]*model.Tag
}

func (p *GetMusicListResult) Resp() *api.GetMusicListResp {
	return &api.GetMusicListResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Musics: collext.Pick(p.Musics, func(music *model.Music) *api.Music {
			return &api.Music{
				Id:   music.Id,
				Name: music.Name,
				Tags: collext.Pick(p.MusicTagsMap[music.Id], getter.TagToDetail),
			}
		}),
	}
}
