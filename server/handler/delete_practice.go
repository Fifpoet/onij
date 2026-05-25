package handler

import (
	"context"
	"errors"

	"onij/biz/prm"
	api "onij/model/api"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var errInvalidPracticeAt = errors.New("invalid practice_at")

// DeletePractice .
// @router /practice/delete [POST]
func DeletePractice(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DeletePracticeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Id <= 0 {
		c.String(consts.StatusBadRequest, "invalid id")
		return
	}

	resp, err := ser.PracticeLogic.Delete(ctx, prm.NewDeletePracticeParam(&req))
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}
