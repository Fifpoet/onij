package prm

import (
	"onij/model/api"
	"onij/util"
)

type UploadArtistParam struct {
	Name         string
	AvatarFileId *string
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

type SearchArtistParam struct {
	Keyword  string
	TagGroup *int32
	TagType  *int32
	Page     int32
	Limit    int32
}
