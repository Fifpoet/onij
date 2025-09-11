package prm

import (
	"onij/biz/biz"
	"onij/biz/getter"
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
)

type UploadMusicParam struct {
	Id           *int64
	Name         string
	SubName      string
	ArtistIds    []int64
	AudioFileId  int64
	LyricFileId  int64
	TimeLength   int32
	ComposerIds  []int64
	WriterIds    []int64
	AlbumId      *int64
	MvUrl        *string
	RootMusicId  *int64
	IssueTime    *int32
	Priority     *int32
	LyricContent *string
	AudioQuality api.AudioQuality
	PerformType  api.PerformType
}

func NewUploadMusicParam(req *api.UploadMusicReq) *UploadMusicParam {
	return &UploadMusicParam{
		Id:           req.MusicId,
		Name:         req.Name,
		SubName:      req.SubName,
		ArtistIds:    req.ArtistIds,
		AudioFileId:  req.AudioFileId,
		LyricFileId:  req.LyricFileId,
		TimeLength:   req.TimeLength,
		ComposerIds:  req.ComposerIds,
		WriterIds:    req.WriterIds,
		AlbumId:      req.AlbumId,
		MvUrl:        req.MvUrl,
		RootMusicId:  req.RootMusicId,
		IssueTime:    req.IssueTime,
		Priority:     req.Priority,
		LyricContent: req.LyricContent,
		AudioQuality: req.AudioQuality,
		PerformType:  req.PerformType,
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
	Music *biz.MusicPrime
}

func (p *GetMusicDetailResult) Resp() *api.GetMusicDetailResp {
	return &api.GetMusicDetailResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Music:   biz.Music(p.Music),
	}
}

type GetMusicListParam struct {
	Keyword   string
	ArtistIds []int64
	TagTypes  []int32
	Page      int32
	Limit     int32
}

func NewGetMusicListParam(req *api.GetMusicListReq) *GetMusicListParam {
	return &GetMusicListParam{
		Keyword:   req.Keyword,
		ArtistIds: req.ArtistIds,
		TagTypes:  collext.Pick(req.TagTypes, func(tt api.TagType) int32 { return int32(tt) }),
		Page:      req.Page,
		Limit:     req.Limit,
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
