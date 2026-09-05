package prm

import (
	"onij/model"
	"onij/model/api"
	"onij/util"
)

const (
	VideoPlayModeMP4         = "mp4"
	VideoPlayModeHLS         = "hls"
	VideoPlayModeProcessing  = "processing"
	VideoPlayModeUnsupported = "unsupported"
)

type PlayFileVideoParam struct {
	FileId int64
}

func NewPlayFileVideoParam(req *api.PlayFileVideoReq) *PlayFileVideoParam {
	return &PlayFileVideoParam{FileId: req.FileId}
}

type PlayFileVideoResult struct {
	Mode       string
	URL        string
	Playlist   string
	DurationMs int64
	Hint       string
}

func (r *PlayFileVideoResult) Resp() *api.PlayFileVideoResp {
	msg := util.BaseMsgOK
	if r.Hint != "" {
		msg = r.Hint
	}
	return &api.PlayFileVideoResp{
		Code:       util.BaseCodeOK,
		Message:    msg,
		Mode:       r.Mode,
		Url:        r.URL,
		Playlist:   r.Playlist,
		DurationMs: r.DurationMs,
	}
}

func ToFileVideoMarkerDTO(m *model.FileVideoMarker) *api.FileVideoMarkerItem {
	if m == nil {
		return nil
	}
	return &api.FileVideoMarkerItem{
		Id:     m.Id,
		FileId: m.FileId,
		TimeMs: m.TimeMs,
		Label:  m.Label,
	}
}

type ListFileVideoMarkerParam struct {
	FileId int64
}

func NewListFileVideoMarkerParam(req *api.ListFileVideoMarkerReq) *ListFileVideoMarkerParam {
	return &ListFileVideoMarkerParam{FileId: req.FileId}
}

type ListFileVideoMarkerResult struct {
	Markers []*api.FileVideoMarkerItem
}

func (r *ListFileVideoMarkerResult) Resp() *api.ListFileVideoMarkerResp {
	markers := r.Markers
	if markers == nil {
		markers = []*api.FileVideoMarkerItem{}
	}
	return &api.ListFileVideoMarkerResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Markers: markers,
	}
}

type SaveFileVideoMarkerParam struct {
	FileId  int64
	Markers []*api.FileVideoMarkerItem
}

func NewSaveFileVideoMarkerParam(req *api.SaveFileVideoMarkerReq) *SaveFileVideoMarkerParam {
	return &SaveFileVideoMarkerParam{FileId: req.FileId, Markers: req.Markers}
}

type SaveFileVideoMarkerResult struct {
	Markers []*api.FileVideoMarkerItem
}

func (r *SaveFileVideoMarkerResult) Resp() *api.SaveFileVideoMarkerResp {
	markers := r.Markers
	if markers == nil {
		markers = []*api.FileVideoMarkerItem{}
	}
	return &api.SaveFileVideoMarkerResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Markers: markers,
	}
}

type DeleteFileVideoMarkerParam struct {
	Id int64
}

func NewDeleteFileVideoMarkerParam(req *api.DeleteFileVideoMarkerReq) *DeleteFileVideoMarkerParam {
	return &DeleteFileVideoMarkerParam{Id: req.Id}
}

type DeleteFileVideoMarkerResult struct{}

func (r *DeleteFileVideoMarkerResult) Resp() *api.DeleteFileVideoMarkerResp {
	return &api.DeleteFileVideoMarkerResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
	}
}
