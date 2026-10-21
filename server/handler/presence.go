package handler

import (
	"context"
	"strings"

	"onij/biz/prm"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type presenceHeartbeatReq struct {
	ClientId           string                  `json:"client_id"`
	Nickname           string                  `json:"nickname"`
	NowPlaying         *prm.PresenceTrackDTO   `json:"now_playing"`
	Queue              []*prm.PresenceTrackDTO `json:"queue"`
	PlaybackCurrentSec float64                 `json:"playback_current_sec"`
	IsPlaying          bool                    `json:"is_playing"`
}

// PresenceHeartbeat @router /presence/heartbeat [POST]
func PresenceHeartbeat(ctx context.Context, c *app.RequestContext) {
	var req presenceHeartbeatReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(req.ClientId) == "" {
		c.String(consts.StatusBadRequest, "client_id is empty")
		return
	}
	if strings.TrimSpace(req.Nickname) == "" {
		c.String(consts.StatusBadRequest, "nickname is empty")
		return
	}
	queue := req.Queue
	if len(queue) > prm.PresenceQueueLimit {
		queue = queue[:prm.PresenceQueueLimit]
	}

	resp, err := ser.PresenceLogic.Heartbeat(ctx, &prm.HeartbeatParam{
		ClientId:           req.ClientId,
		Nickname:           req.Nickname,
		NowPlaying:         req.NowPlaying,
		Queue:              queue,
		PlaybackCurrentSec: req.PlaybackCurrentSec,
		IsPlaying:          req.IsPlaying,
	})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}
