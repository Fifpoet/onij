package logic

import (
	"context"
	"encoding/json"
	"onij/biz/prm"
	"onij/infra/uvr"
)

type UvrLogic interface {
	Health(ctx context.Context) (json.RawMessage, error)
	ListModels(ctx context.Context) (json.RawMessage, error)
	SeparateInstrumentalFromURL(ctx context.Context, param *prm.SeparateInstrumentalFromURLParam) (*prm.SeparateInstrumentalFromURLResult, error)
	DownloadInstrumental(ctx context.Context, param *prm.DownloadInstrumentalParam) (*prm.DownloadInstrumentalResult, error)
}

type uvrLogic struct {
	client uvr.Client
}

func NewUvrLogic(client uvr.Client) UvrLogic {
	return &uvrLogic{client: client}
}

func (l *uvrLogic) Health(ctx context.Context) (json.RawMessage, error) {
	data, err := l.client.Health(ctx)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(data), nil
}

func (l *uvrLogic) ListModels(ctx context.Context) (json.RawMessage, error) {
	data, err := l.client.ListModels(ctx)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(data), nil
}

func (l *uvrLogic) SeparateInstrumentalFromURL(ctx context.Context, param *prm.SeparateInstrumentalFromURLParam) (*prm.SeparateInstrumentalFromURLResult, error) {
	body, err := json.Marshal(map[string]any{
		"id":             param.ID,
		"source_url":     param.SourceURL,
		"architecture":   param.Architecture,
		"model":          param.Model,
		"use_gpu":        param.UseGPU,
		"output_format":  param.OutputFormat,
		"single_stem":    param.SingleStem,
	})
	if err != nil {
		return nil, err
	}

	data, err := l.client.SeparateInstrumentalFromURL(ctx, body)
	if err != nil {
		return nil, err
	}
	return &prm.SeparateInstrumentalFromURLResult{Raw: data}, nil
}

func (l *uvrLogic) DownloadInstrumental(ctx context.Context, param *prm.DownloadInstrumentalParam) (*prm.DownloadInstrumentalResult, error) {
	data, contentType, filename, err := l.client.DownloadInstrumental(ctx, param.JobID)
	if err != nil {
		return nil, err
	}
	return &prm.DownloadInstrumentalResult{
		Data:        data,
		ContentType: contentType,
		Filename:    filename,
	}, nil
}
