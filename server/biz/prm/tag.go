package prm

import (
	"onij/biz/getter"
	"onij/infra/mysql"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
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
	ListShow     bool
}

func NewUploadTagParam(req *api.UploadTagReq) *UploadTagParam {
	return &UploadTagParam{
		ResourceId:   req.TagDetail.ResourceId,
		ResourceType: int32(req.TagDetail.ResourceType),
		TagBiz:       int32(req.TagDetail.TagBiz),
		TagGroup:     int32(req.TagDetail.TagGroup),
		TagType:      int32(req.TagDetail.TagType),
		TargetId:     req.TagDetail.TargetId,
		TargetType:   (*int32)(req.TagDetail.TargetType),
		Extra:        req.TagDetail.Extra,
		ListShow:     req.TagDetail.ListShow,
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

type GetTagListParam struct {
	ResourceIds []int64
	TagGroups   []int32
	TagTypes    []int32
}

func NewGetTagListParam(req *api.GetTagListReq) *GetTagListParam {
	return &GetTagListParam{
		ResourceIds: req.ResourceIds,
		TagGroups:   collext.Pick(req.TagGroups, func(t api.TagGroup) int32{return int32(t)}),
		TagTypes:    collext.Pick(req.TagTypes, func(t api.TagType) int32{return int32(t)}),
	}
}

type GetTagListResult struct {
	Tags []*mysql.Tag
}

func (r *GetTagListResult) Resp() *api.GetTagListResp {
	return &api.GetTagListResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Tags:    collext.Pick(r.Tags, getter.TagToDetail),
	}
}

type DeleteTagParam struct {
	ResourceId   int64
	TagType       int32
}
func NewDeleteTagParam(req *api.DeleteTagReq) *DeleteTagParam {
	return &DeleteTagParam{
		ResourceId:   req.ResourceId,
		TagType:       int32(req.TagType),
	}
}
type DeleteTagResult struct {}
func (r *DeleteTagResult) Resp() *api.DeleteTagResp {
	return &api.DeleteTagResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
	}
}
