package uvr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultBaseURL         = "http://172.26.250.1:5555"
	separateRequestTimeout = 30 * time.Minute
	downloadRequestTimeout = 10 * time.Minute
	defaultRequestTimeout  = 30 * time.Second
)

type APIError struct {
	StatusCode int
	Body       []byte
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("uvr api error: status=%d", e.StatusCode)
}

type Client interface {
	Health(ctx context.Context) ([]byte, error)
	ListModels(ctx context.Context) ([]byte, error)
	SeparateInstrumentalFromURL(ctx context.Context, body []byte) ([]byte, error)
	DownloadInstrumental(ctx context.Context, jobID string) (body io.ReadCloser, contentLength int64, contentType string, filename string, err error)
}

type client struct {
	baseURL string
	http    *http.Client
}

func NewClient() Client {
	base := strings.TrimSpace(os.Getenv("UVR_API_BASE"))
	if base == "" {
		base = defaultBaseURL
	}
	return &client{
		baseURL: strings.TrimRight(base, "/"),
		http:    &http.Client{Timeout: 0},
	}
}

func (c *client) Health(ctx context.Context) ([]byte, error) {
	return c.get(ctx, "/health", defaultRequestTimeout)
}

func (c *client) ListModels(ctx context.Context) ([]byte, error) {
	return c.get(ctx, "/api/v1/models", defaultRequestTimeout)
}

func (c *client) SeparateInstrumentalFromURL(ctx context.Context, body []byte) ([]byte, error) {
	return c.postJSON(ctx, "/api/v1/separate/instrumental/from-url", body, separateRequestTimeout)
}

func (c *client) DownloadInstrumental(ctx context.Context, jobID string) (io.ReadCloser, int64, string, string, error) {
	url := fmt.Sprintf("%s/api/v1/jobs/%s/instrumental", c.baseURL, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, "", "", err
	}

	httpClient := &http.Client{Timeout: downloadRequestTimeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, "", "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		defer resp.Body.Close()
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		return nil, 0, "", "", newAPIError(resp.StatusCode, data)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "audio/mpeg"
	}
	filename := parseFilename(resp.Header.Get("Content-Disposition"))
	return resp.Body, resp.ContentLength, contentType, filename, nil
}

func (c *client) get(ctx context.Context, path string, timeout time.Duration) ([]byte, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Timeout: timeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, newAPIError(resp.StatusCode, data)
	}
	return data, nil
}

func (c *client) postJSON(ctx context.Context, path string, body []byte, timeout time.Duration) ([]byte, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: timeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, newAPIError(resp.StatusCode, data)
	}
	return data, nil
}

func newAPIError(statusCode int, body []byte) *APIError {
	msg := strings.TrimSpace(string(body))
	var payload struct {
		Detail string `json:"detail"`
	}
	if json.Unmarshal(body, &payload) == nil && payload.Detail != "" {
		msg = payload.Detail
	}
	return &APIError{
		StatusCode: statusCode,
		Body:       body,
		Message:    msg,
	}
}

func parseFilename(contentDisposition string) string {
	if contentDisposition == "" {
		return "instrumental.mp3"
	}
	const marker = "filename="
	idx := strings.Index(strings.ToLower(contentDisposition), marker)
	if idx < 0 {
		return "instrumental.mp3"
	}
	name := strings.TrimSpace(contentDisposition[idx+len(marker):])
	name = strings.Trim(name, `"`)
	if name == "" {
		return "instrumental.mp3"
	}
	return name
}
