package prm

import (
	"onij/biz/biz"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
)

type UploadAlbumParam struct {
	Name        string
	Profile     string
	ArtistIds   []int64
	CoverFileId int64
	IssueTime   int32
	AlbumType   api.AlbumType
	LiveUrl     string
	AlbumMusics []*api.UploadAlbumReq_AlbumMusic
	AlbumId     *int64
}

func NewUploadAlbumParam(req *api.UploadAlbumReq) *UploadAlbumParam {
	return &UploadAlbumParam{
		Name:        req.Name,
		Profile:     req.Profile,
		ArtistIds:   req.ArtistIds,
		CoverFileId: req.CoverFileId,
		IssueTime:   req.IssueTime,
		AlbumType:   req.AlbumType,
		LiveUrl:     req.LiveUrl,
		AlbumMusics: req.AlbumMusics,
		AlbumId:     req.AlbumId,
	}
}

type UploadAlbumResult struct {
	AlbumId int64
}

func (r *UploadAlbumResult) Resp() *api.UploadAlbumResp {
	return &api.UploadAlbumResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		AlbumId: r.AlbumId,
	}
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
	Album *biz.AlbumPrime
}

func (r *GetAlbumDetailResult) Resp() *api.GetAlbumDetailResp {
	return &api.GetAlbumDetailResp{
		Code:        util.BaseCodeOK,
		Message:     util.BaseMsgOK,
		Album:       biz.Album(r.Album),
		AlbumMusics: collext.Pick(r.Album.AlbumMusics, biz.AlbumMusic),
	}
}
