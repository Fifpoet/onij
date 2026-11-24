package handler

import (
	"context"

	"onij/biz/prm"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type playSourceSaveReq struct {
	Sources []*prm.PlaySourceDTO `json:"sources"`
}

type playSourceSignSaveReq struct {
	Version string `json:"version"`
	Crc32   string `json:"crc32"`
}

// ListPlaySource @router /play-source/list [POST]
func ListPlaySource(ctx context.Context, c *app.RequestContext) {
	_ = c.BindAndValidate(&struct{}{})
	resp, err := ser.PlaySourceLogic.List(ctx, &prm.ListPlaySourceParam{})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

// SavePlaySource @router /play-source/save [POST]
func SavePlaySource(ctx context.Context, c *app.RequestContext) {
	var req playSourceSaveReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp, err := ser.PlaySourceLogic.Save(ctx, &prm.SavePlaySourceParam{Sources: req.Sources})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

// GetGdstudioSign @router /play-source/gdstudio-sign [POST]
func GetGdstudioSign(ctx context.Context, c *app.RequestContext) {
	_ = c.BindAndValidate(&struct{}{})
	resp, err := ser.PlaySourceLogic.GetSign(ctx, &prm.GetGdstudioSignParam{})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

// SaveGdstudioSign @router /play-source/gdstudio-sign/save [POST]
func SaveGdstudioSign(ctx context.Context, c *app.RequestContext) {
	var req playSourceSignSaveReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp, err := ser.PlaySourceLogic.SaveSign(ctx, &prm.SaveGdstudioSignParam{
		Version: req.Version,
		Crc32:   req.Crc32,
	})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}
