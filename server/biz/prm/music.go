package prm

import (
	"onij/infra/mysql"
	"onij/model/api"
)

type UploadMusicParam struct {
	Id           *int64
	Name         string
	ArtistIds    []int64
	Mp3FileId    int64
	LyricsFileId int64
	ComposerId   *int64
	WriterId     *int64
	AlbumId      *int64
	MvUrl        *string
	RootMusicId  *int64
	IssueTime    *int32
}

func NewUploadMusicParam(req *api.UploadMusicReq) *UploadMusicParam {
	return &UploadMusicParam{
		Id:           req.MusicId,
		Name:         req.Name,
		ArtistIds:    req.ArtistIds,
		Mp3FileId:    req.Mp3FileId,
		LyricsFileId: req.LyricsFileId,
		ComposerId:   req.ComposerId,
		WriterId:     req.WriterId,
		AlbumId:      req.AlbumId,
		MvUrl:        req.MvUrl,
		RootMusicId:  req.RootMusicId,
		IssueTime:    req.IssueTime,
	}
}

type UploadMusicResult struct {
}

func (p *UploadMusicResult) Resp() *api.UploadMusicResp {
	return &api.UploadMusicResp{}
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
	Music         *mysql.Music
	Artists       []*mysql.Artist
	Composer      *mysql.Artist
	Writer        *mysql.Artist
	Mp3FileUrl    string
	LyricsFileUrl string
	Albums        []*mysql.Album
}

func (p *GetMusicDetailResult) Resp() *api.GetMusicDetailResp {
	return &api.GetMusicDetailResp{}
}

type GetMusicListParam struct {
	SortType api.MusicSortType
	Keyword  *string
	AlbumId  *int64
	ArtistId *int64
	Page     int32
	Limit    int32
}

func NewGetMusicListParam(req *api.GetMusicListReq) *GetMusicListParam {
	return &GetMusicListParam{
		SortType: req.SortType,
		Keyword:  req.Keyword,
		AlbumId:  req.AlbumId,
		ArtistId: req.ArtistId,
		Page:     req.Page,
		Limit:    req.Limit,
	}
}

type GetMusicListResult struct {
	Musics []*mysql.Music
}

func (p *GetMusicListResult) Resp() *api.GetMusicListResp {
	return &api.GetMusicListResp{}
}
