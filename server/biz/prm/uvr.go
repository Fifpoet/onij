package prm

import (
	"encoding/json"
	"io"
)

type SeparateInstrumentalFromURLParam struct {
	ID           string
	SourceURL    string
	Architecture string
	Model        string
	UseGPU       bool
	OutputFormat string
	SingleStem   string
}

type SeparateInstrumentalFromURLResult struct {
	Raw json.RawMessage
}

type DownloadInstrumentalParam struct {
	JobID string
}

type DownloadInstrumentalResult struct {
	Body          io.ReadCloser
	ContentLength int64 // <0 表示未知长度，流式 chunked
	ContentType   string
	Filename      string
}
