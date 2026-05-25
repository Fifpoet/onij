package handler

import (
	"context"

	"onij/biz/prm"
	api "onij/model/api"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// UploadPractice .
// @router /practice/upload [POST]
func UploadPractice(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.UploadPracticeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if err = checkUploadPracticeReq(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp, err := ser.PracticeLogic.Upload(ctx, prm.NewUploadPracticeParam(&req))
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

func checkUploadPracticeReq(req *api.UploadPracticeReq) error {
	if req.PracticeAt <= 0 {
		return errInvalidPracticeAt
	}
	return nil
}
