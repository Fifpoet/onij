package handler

import (
	"context"

	"onij/biz/prm"
	api "onij/model/api"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func PlayFileVideo(ctx context.Context, c *app.RequestContext) {
	var req api.PlayFileVideoReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.FileId <= 0 {
		c.String(consts.StatusBadRequest, "invalid file_id")
		return
	}
	resp, err := ser.FileVideoLogic.Play(ctx, prm.NewPlayFileVideoParam(&req))
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

func ListFileVideoMarker(ctx context.Context, c *app.RequestContext) {
	var req api.ListFileVideoMarkerReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.FileId <= 0 {
		c.String(consts.StatusBadRequest, "invalid file_id")
		return
	}
	resp, err := ser.FileVideoLogic.ListMarkers(ctx, prm.NewListFileVideoMarkerParam(&req))
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

func SaveFileVideoMarker(ctx context.Context, c *app.RequestContext) {
	var req api.SaveFileVideoMarkerReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.FileId <= 0 {
		c.String(consts.StatusBadRequest, "invalid file_id")
		return
	}
	resp, err := ser.FileVideoLogic.SaveMarkers(ctx, prm.NewSaveFileVideoMarkerParam(&req))
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

func DeleteFileVideoMarker(ctx context.Context, c *app.RequestContext) {
	var req api.DeleteFileVideoMarkerReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Id <= 0 {
		c.String(consts.StatusBadRequest, "invalid id")
		return
	}
	resp, err := ser.FileVideoLogic.DeleteMarker(ctx, prm.NewDeleteFileVideoMarkerParam(&req))
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}
