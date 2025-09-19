package prm

import (
	"onij/model"
	"onij/model/api"
	"onij/util"
)

type UploadArtistParam struct {
	Name         string
	AvatarFileId *int64
	ArtistId     *string
}

func NewUploadArtistParam(req *api.UploadArtistReq) *UploadArtistParam {
	return &UploadArtistParam{
		Name:         req.Name,
		AvatarFileId: req.AvatarFileId,
		ArtistId:     req.ArtistId,
	}
}

type UploadArtistResult struct {
	ArtistId int64
}

func (r *UploadArtistResult) Resp() *api.UploadArtistResp {
	return &api.UploadArtistResp{
		Code:     util.BaseCodeOK,
		Message:  util.BaseMsgOK,
		ArtistId: r.ArtistId,
	}
}

type GetArtistListParam struct {
	Keyword string
	Name    string
	Page    int32
	Limit   int32
}

func NewGetArtistListParam(req *api.GetArtistListReq) *GetArtistListParam {
	return &GetArtistListParam{
		Keyword: req.Keyword,
		Name:    req.Name,
		Page:    req.Page,
		Limit:   req.Limit,
	}
}

type GetArtistListResult struct {
	Artists []*model.Artist
	Total   int32
}

func (r *GetArtistListResult) Resp() *api.GetArtistListResp {
	return &api.GetArtistListResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Artists: model.ArtistList(r.Artists),
	}
}
