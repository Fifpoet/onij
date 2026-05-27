package handler

import (
	"context"
	"errors"
	"fmt"
	"onij/biz/prm"
	"onij/infra/uvr"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type separateInstrumentalFromURLReq struct {
	ID           string `json:"id"`
	SourceURL    string `json:"source_url" vd:"len($)>0"`
	Architecture string `json:"architecture"`
	Model        string `json:"model"`
	UseGPU       *bool  `json:"use_gpu"`
	OutputFormat string `json:"output_format"`
	SingleStem   string `json:"single_stem"`
}

func writeUvrAPIError(ctx context.Context, c *app.RequestContext, err error) {
	var apiErr *uvr.APIError
	if errors.As(err, &apiErr) {
		if len(apiErr.Body) > 0 {
			c.Data(apiErr.StatusCode, "application/json", apiErr.Body)
			return
		}
		c.String(apiErr.StatusCode, apiErr.Message)
		return
	}
	c.String(consts.StatusBadGateway, err.Error())
}

// UvrHealth 代理 UVR 健康检查
// @router /uvr/health [GET]
func UvrHealth(ctx context.Context, c *app.RequestContext) {
	data, err := ser.UvrLogic.Health(ctx)
	if err != nil {
		writeUvrAPIError(ctx, c, err)
		return
	}
	c.Data(consts.StatusOK, "application/json", data)
}

// UvrListModels 代理 UVR 模型列表
// @router /uvr/models [GET]
func UvrListModels(ctx context.Context, c *app.RequestContext) {
	data, err := ser.UvrLogic.ListModels(ctx)
	if err != nil {
		writeUvrAPIError(ctx, c, err)
		return
	}
	c.Data(consts.StatusOK, "application/json", data)
}

// SeparateInstrumentalFromURL 代理 UVR 从 URL 提取伴奏
// @router /uvr/instrumental/from-url [POST]
func SeparateInstrumentalFromURL(ctx context.Context, c *app.RequestContext) {
	var req separateInstrumentalFromURLReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	useGPU := true
	if req.UseGPU != nil {
		useGPU = *req.UseGPU
	}
	architecture := req.Architecture
	if architecture == "" {
		architecture = "mdx"
	}
	model := req.Model
	if model == "" {
		model = "UVR-MDX-NET-Inst_HQ_3.onnx"
	}
	outputFormat := req.OutputFormat
	if outputFormat == "" {
		outputFormat = "MP3"
	}
	singleStem := req.SingleStem
	if singleStem == "" {
		singleStem = "Instrumental"
	}

	resp, err := ser.UvrLogic.SeparateInstrumentalFromURL(ctx, &prm.SeparateInstrumentalFromURLParam{
		ID:           req.ID,
		SourceURL:    req.SourceURL,
		Architecture: architecture,
		Model:        model,
		UseGPU:       useGPU,
		OutputFormat: outputFormat,
		SingleStem:   singleStem,
	})
	if err != nil {
		writeUvrAPIError(ctx, c, err)
		return
	}
	c.Data(consts.StatusOK, "application/json", resp.Raw)
}

// DownloadInstrumental 代理 UVR 伴奏下载
// @router /uvr/jobs/:job_id/instrumental [GET]
func DownloadInstrumental(ctx context.Context, c *app.RequestContext) {
	jobID := c.Param("job_id")
	if jobID == "" {
		c.String(consts.StatusBadRequest, "job_id is empty")
		return
	}

	resp, err := ser.UvrLogic.DownloadInstrumental(ctx, &prm.DownloadInstrumentalParam{JobID: jobID})
	if err != nil {
		writeUvrAPIError(ctx, c, err)
		return
	}
	c.Header("Content-Type", resp.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, resp.Filename))
	if resp.ContentLength >= 0 {
		c.Header("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	}
	c.SetBodyStream(resp.Body, int(resp.ContentLength))
}
