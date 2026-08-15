package logic

import (
	"context"
	"encoding/json"
	"io"
	"onij/biz/prm"
	"onij/infra/uvr"
)

type UvrLogic interface {
	Health(ctx context.Context) (json.RawMessage, error)
	ListModels(ctx context.Context) (json.RawMessage, error)
	SeparateInstrumentalFromURL(ctx context.Context, param *prm.SeparateInstrumentalFromURLParam) (*prm.SeparateInstrumentalFromURLResult, error)
	DownloadInstrumental(ctx context.Context, param *prm.DownloadInstrumentalParam) (*prm.DownloadInstrumentalResult, error)
	GetPitch(ctx context.Context, jobID string) (json.RawMessage, error)
	WhisperHealth(ctx context.Context) (json.RawMessage, error)
	WhisperTranscribe(ctx context.Context, filename string, file io.Reader, language string, vadFilter bool) (json.RawMessage, error)
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

func (l *uvrLogic) WhisperHealth(ctx context.Context) (json.RawMessage, error) {
	data, err := l.client.WhisperHealth(ctx)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(data), nil
}

func (l *uvrLogic) WhisperTranscribe(ctx context.Context, filename string, file io.Reader, language string, vadFilter bool) (json.RawMessage, error) {
	data, err := l.client.WhisperTranscribe(ctx, filename, file, language, vadFilter)
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
		"id":            param.ID,
		"source_url":    param.SourceURL,
		"architecture":  param.Architecture,
		"model":         param.Model,
		"use_gpu":       param.UseGPU,
		"output_format": param.OutputFormat,
		"single_stem":   param.SingleStem,
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
	body, contentLength, contentType, filename, err := l.client.DownloadInstrumental(ctx, param.JobID)
	if err != nil {
		return nil, err
	}
	return &prm.DownloadInstrumentalResult{
		Body:          body,
		ContentLength: contentLength,
		ContentType:   contentType,
		Filename:      filename,
	}, nil
}

func (l *uvrLogic) GetPitch(ctx context.Context, jobID string) (json.RawMessage, error) {
	data, err := l.client.GetPitch(ctx, jobID)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(data), nil
}
