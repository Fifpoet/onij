package handler

import (
	"context"
	"onij/biz/prm"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type downloadFileReq struct {
	FileIds []int64 `json:"file_ids"`
}

// DownloadFile 返回七牛私有下载链接
// @router /file/download [POST]
func DownloadFile(ctx context.Context, c *app.RequestContext) {
	var req downloadFileReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if len(req.FileIds) == 0 {
		c.String(consts.StatusBadRequest, "file_ids is empty")
		return
	}

	resp, err := ser.FileLogic.Download(ctx, &prm.DownloadFileParam{FileIds: req.FileIds})
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}
