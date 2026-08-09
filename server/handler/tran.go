package handler

import (
	"context"
	"fmt"

	"onij/biz/prm"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type listTranReq struct{}

type addTranTextReq struct {
	Content string `json:"content"`
}

type addTranFileReq struct {
	FileId int64  `json:"file_id"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Format int32  `json:"format"`
}

type tranIdReq struct {
	Id int64 `json:"id,string"`
}

// ListTran @router /tran/list [POST]
func ListTran(ctx context.Context, c *app.RequestContext) {
	var req listTranReq
	_ = c.BindAndValidate(&req)
	resp, err := ser.TranLogic.List(ctx, &prm.ListTranParam{})
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

// AddTranText @router /tran/text [POST]
func AddTranText(ctx context.Context, c *app.RequestContext) {
	var req addTranTextReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp, err := ser.TranLogic.AddText(ctx, &prm.AddTranTextParam{Content: req.Content})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

// AddTranFile @router /tran/file [POST]
func AddTranFile(ctx context.Context, c *app.RequestContext) {
	var req addTranFileReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.FileId <= 0 {
		c.String(consts.StatusBadRequest, "invalid file_id")
		return
	}
	resp, err := ser.TranLogic.AddFile(ctx, &prm.AddTranFileParam{
		FileId: req.FileId,
		Name:   req.Name,
		Size:   req.Size,
		Format: req.Format,
	})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

// PinTran @router /tran/pin [POST]
func PinTran(ctx context.Context, c *app.RequestContext) {
	var req tranIdReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Id <= 0 {
		c.String(consts.StatusBadRequest, fmt.Sprintf("invalid id: %d", req.Id))
		return
	}
	resp, err := ser.TranLogic.TogglePin(ctx, &prm.PinTranParam{Id: req.Id})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

// DeleteTran @router /tran/delete [POST]
func DeleteTran(ctx context.Context, c *app.RequestContext) {
	var req tranIdReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Id <= 0 {
		c.String(consts.StatusBadRequest, "invalid id")
		return
	}
	resp, err := ser.TranLogic.Delete(ctx, &prm.DeleteTranParam{Id: req.Id})
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}
