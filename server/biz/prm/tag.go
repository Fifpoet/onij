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
}

func NewUploadTagParam(req *api.UploadTagReq) *UploadTagParam {
	return &UploadTagParam{
		ResourceId:   req.Tag.ResourceId,
		ResourceType: int32(req.Tag.ResourceType),
		TagBiz:       int32(req.Tag.TagBiz),
		TagGroup:     int32(req.Tag.TagGroup),
		TagType:      int32(req.Tag.TagType),
		TargetId:     req.Tag.TargetId,
		TargetType:   (*int32)(req.Tag.TargetType),
		Extra:        req.Tag.Extra,
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

type DeleteTagParam struct {
	ResourceId int64
	TagType    int32
}

func NewDeleteTagParam(req *api.DeleteTagReq) *DeleteTagParam {
	return &DeleteTagParam{
		ResourceId: req.ResourceId,
		TagType:    int32(req.TagType),
	}
}

type DeleteTagResult struct{}

func (r *DeleteTagResult) Resp() *api.DeleteTagResp {
	return &api.DeleteTagResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
	}
}
