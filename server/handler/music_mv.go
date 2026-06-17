package handler

import (
	"context"
	"errors"

	"onij/biz/prm"
	api "onij/model/api"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const maxMusicMvBatch = 64

// UpdateMusicMv 更新歌曲 MV 内嵌地址
// @router /music/mv [POST]
func UpdateMusicMv(ctx context.Context, c *app.RequestContext) {
	var req api.UpdateMusicMvReq
	if err := c.BindJSON(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if err := checkUpdateMusicMvReq(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp, err := ser.MusicLogic.UpdateMv(ctx, prm.NewUpdateMusicMvParam(&req))
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

func checkUpdateMusicMvReq(req *api.UpdateMusicMvReq) error {
	hasMusicID := req.MusicId != nil && *req.MusicId > 0
	hasThirdID := req.ThirdId != nil && *req.ThirdId > 0
	if !hasMusicID && !hasThirdID {
		return errors.New("music_id or third_id required")
	}
	return nil
}

// GetMusicMvBatch 批量查询 MV（播放前预取）
// @router /music/mv/batch [POST]
func GetMusicMvBatch(ctx context.Context, c *app.RequestContext) {
	var req api.GetMusicMvBatchReq
	if err := c.BindJSON(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if len(req.ThirdIds) == 0 {
		c.String(consts.StatusBadRequest, "third_ids required")
		return
	}
	if len(req.ThirdIds) > maxMusicMvBatch {
		c.String(consts.StatusBadRequest, "third_ids too many")
		return
	}

	resp, err := ser.MusicLogic.GetMvBatch(ctx, prm.NewGetMusicMvBatchParam(&req))
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}
