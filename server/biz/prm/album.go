package prm

import (
	"onij/infra/mysql"
	"onij/model/api"
	"onij/util/boost/collection/collext"
)

type UploadAlbumParam struct {
	Name           string
	ArtistId       int64
	CoverFileId    int64
	IssueTime      int32
	MusicId        *int64
	RelatedAlbumId *int64
}

func NewUploadAlbumParam(req *api.UploadAlbumReq) *UploadAlbumParam {
	return &UploadAlbumParam{
		Name:           req.Name,
		ArtistId:       req.ArtistId,
		CoverFileId:    req.CoverFileId,
		IssueTime:      req.IssueTime,
		MusicId:        req.MusicId,
		RelatedAlbumId: req.RelatedAlbumId,
	}
}

type UploadAlbumResult struct {
}

func (r *UploadAlbumResult) Resp() *api.UploadAlbumResp {
	return &api.UploadAlbumResp{}
}

type GetAlbumDetailParam struct {
	AlbumId int64
}

func NewGetAlbumParam(req *api.GetAlbumDetailReq) *GetAlbumDetailParam {
	return &GetAlbumDetailParam{
		AlbumId: req.AlbumId,
	}
}

type GetAlbumDetailResult struct {
	Album    *mysql.Album
	Musics   []*mysql.Music
	CoverUrl string
	Artist   *mysql.Artist
}

func (r *GetAlbumDetailResult) Resp() *api.GetAlbumDetailResp {
	return &api.GetAlbumDetailResp{
		Detail: &api.AlbumDetail{
			Id:           r.Album.Id,
			Name:         r.Album.Name,
			ArtistId:     r.Artist.Id,
			ArtistName:   r.Artist.Name,
			IssueTime:    r.Album.IssueTime,
			CoverFileUrl: r.CoverUrl,
			MvUrl:        r.Album.MvUrl,
			RelatedMusics: collext.Pick(r.Musics, func(music *mysql.Music) *api.MusicProfile {
				return &api.MusicProfile{
					Id:   music.Id,
					Name: music.Name,
				}
			}),
		},
	}
}
