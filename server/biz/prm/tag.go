package prm

import (
	"onij/model/api"
	"onij/util"
)

type UploadTagParam struct {
	ResourceId   int64
	ResourceType int32
	TagBiz       int32
	TagGroup     int32
	TagType      int32
	TargetId     *int64
	TargetType   *int32
	Extra        string
	ListShow bool
}

func NewUploadTagParam(req *api.UploadTagReq) *UploadTagParam {
	return &UploadTagParam{
		ResourceId:   req.ResourceId,
		ResourceType: req.ResourceType,
		TagBiz:       int32(req.TagBiz),
		TagGroup:     int32(req.TagGroup),
		TagType:      int32(req.TagType),
		TargetId:     req.TargetId,
		TargetType:   (*int32)(req.TargetType),
		Extra:        req.Extra,
		ListShow: req.ListShow,
	}
}

type UploadTagResult struct {
}

func (r *UploadTagResult) Resp() *api.UploadTagResp {
	return &api.UploadTagResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
	}
}
