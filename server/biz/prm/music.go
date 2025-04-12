package prm

import (
	"onij/infra/mysql"
	"onij/model/api"
	"onij/util/boost/collection/collext"
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
	MusicId int64
}

func (p *UploadMusicResult) Resp() *api.UploadMusicResp {
	return &api.UploadMusicResp{
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
	Music         *mysql.Music
	Artists       []*mysql.Artist
	Composer      *mysql.Artist
	Writer        *mysql.Artist
	Mp3FileUrl    string
	LyricsFileUrl string
	Albums        []*mysql.Album
}

func (p *GetMusicDetailResult) Resp() *api.GetMusicDetailResp {
	album := &mysql.Album{}
	if len(p.Albums) > 0 {
		album = p.Albums[0]
	}
	return &api.GetMusicDetailResp{
		Detail: &api.MusicDetail{
			Id:            p.Music.Id,
			Name:          p.Music.Name,
			ArtistIds:     collext.Pick(p.Artists, func(artist *mysql.Artist) int64 { return artist.Id }),
			ArtistNames:   collext.Pick(p.Artists, func(artist *mysql.Artist) string { return artist.Name }),
			ComposerId:    p.Composer.Id,
			ComposerName:  p.Composer.Name,
			WriterId:      p.Writer.Id,
			WriterName:    p.Writer.Name,
			IssueTime:     p.Music.IssueTime,
			MvUrl:         p.Music.MvUrl,
			Mp3FileUrl:    p.Mp3FileUrl,
			LyricsFileUrl: p.LyricsFileUrl,
			AlbumId:       album.Id,
			AlbumName:     album.Name,
		},
	}
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
