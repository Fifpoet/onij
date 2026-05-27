package prm

import "encoding/json"

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
	Data        []byte
	ContentType string
	Filename    string
}
