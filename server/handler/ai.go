package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"onij/biz/prm"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/protocol/http1/resp"
)

type aiSessionIdReq struct {
	SessionId int64 `json:"session_id,string"`
}

type aiChatReq struct {
	SessionId int64  `json:"session_id,string"`
	Content   string `json:"content"`
}

func CreateAiSession(ctx context.Context, c *app.RequestContext) {
	respBody, err := ser.AiLogic.CreateSession(ctx)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(consts.StatusOK, respBody.Resp())
}

func ListAiSessions(ctx context.Context, c *app.RequestContext) {
	respBody, err := ser.AiLogic.ListSessions(ctx)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(consts.StatusOK, respBody.Resp())
}

func ListAiMessages(ctx context.Context, c *app.RequestContext) {
	var req aiSessionIdReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	respBody, err := ser.AiLogic.ListMessages(ctx, req.SessionId)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, respBody.Resp())
}

func DeleteAiSession(ctx context.Context, c *app.RequestContext) {
	var req aiSessionIdReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	respBody, err := ser.AiLogic.DeleteSession(ctx, req.SessionId)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, respBody.Resp())
}

func AiChat(ctx context.Context, c *app.RequestContext) {
	var req aiChatReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	respBody, err := ser.AiLogic.Chat(ctx, &prm.AiChatParam{
		SessionId: req.SessionId,
		Content:   req.Content,
	})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, respBody.Resp())
}

func AiChatStream(ctx context.Context, c *app.RequestContext) {
	var req aiChatReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	c.SetStatusCode(consts.StatusOK)
	c.Response.Header.Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Response.Header.Set("Cache-Control", "no-cache")
	c.Response.Header.Set("Connection", "keep-alive")
	c.Response.Header.Set("X-Accel-Buffering", "no")

	w := resp.NewChunkedBodyWriter(&c.Response, c.GetWriter())
	c.Response.HijackWriter(w)

	writeSSE := func(event string, payload any) error {
		raw, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if _, err = c.Write([]byte(fmt.Sprintf("event: %s\ndata: %s\n\n", event, raw))); err != nil {
			return err
		}
		return c.Flush()
	}

	err := ser.AiLogic.ChatStream(ctx, &prm.AiChatParam{
		SessionId: req.SessionId,
		Content:   req.Content,
	}, writeSSE)
	if err != nil {
		_ = writeSSE("error", map[string]string{"message": err.Error()})
	}
}

func AiConfirmStream(ctx context.Context, c *app.RequestContext) {
	var req struct {
		ConfirmId int64                 `json:"confirm_id,string"`
		Decision  prm.AiConfirmDecision `json:"decision"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	c.SetStatusCode(consts.StatusOK)
	c.Response.Header.Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Response.Header.Set("Cache-Control", "no-cache")
	c.Response.Header.Set("Connection", "keep-alive")
	c.Response.Header.Set("X-Accel-Buffering", "no")

	w := resp.NewChunkedBodyWriter(&c.Response, c.GetWriter())
	c.Response.HijackWriter(w)

	writeSSE := func(event string, payload any) error {
		raw, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if _, err = c.Write([]byte(fmt.Sprintf("event: %s\ndata: %s\n\n", event, raw))); err != nil {
			return err
		}
		return c.Flush()
	}

	err := ser.AiLogic.ConfirmStream(ctx, &prm.AiConfirmParam{
		ConfirmId: req.ConfirmId,
		Decision:  req.Decision,
	}, writeSSE)
	if err != nil {
		_ = writeSSE("error", map[string]string{"message": err.Error()})
	}
}
