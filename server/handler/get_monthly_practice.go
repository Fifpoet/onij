package handler

import (
	"context"

	"onij/biz/prm"
	api "onij/model/api"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetMonthlyPractice .
// @router /practice/monthly [POST]
func GetMonthlyPractice(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetMonthlyPracticeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Year <= 0 || req.Month < 1 || req.Month > 12 {
		c.String(consts.StatusBadRequest, "invalid year or month")
		return
	}

	resp, err := ser.PracticeLogic.GetMonthly(ctx, prm.NewGetMonthlyPracticeParam(&req))
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}
