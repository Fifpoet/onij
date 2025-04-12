package prm

import (
	"onij/model/api"
	"onij/util"
)

type UploadArtistParam struct {
	Name       string
	ArtistType api.ArtistType
}

func NewUploadArtistParam(req *api.UploadArtistReq) *UploadArtistParam {
	return &UploadArtistParam{
		Name:       req.Name,
		ArtistType: req.ArtistType,
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
