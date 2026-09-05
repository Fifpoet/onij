package biz

import (
	"encoding/json"
	"strings"
)

type fileExtraMeta struct {
	OriginAt int64  `json:"origin_at,omitempty"`
	Exif     string `json:"exif,omitempty"`
	HlsPid   string `json:"hls_pid,omitempty"`
}

func BuildFileExtra(originAt int64, exifRaw string) string {
	if originAt <= 0 && exifRaw == "" {
		return ""
	}
	b, err := json.Marshal(fileExtraMeta{
		OriginAt: originAt,
		Exif:     exifRaw,
	})
	if err != nil {
		return ""
	}
	return string(b)
}

func ParseFileOriginAt(extra string) int64 {
	if extra == "" {
		return 0
	}
	var m fileExtraMeta
	if err := json.Unmarshal([]byte(extra), &m); err != nil {
		return 0
	}
	return m.OriginAt
}

func ParseHlsPid(extra string) string {
	if extra == "" {
		return ""
	}
	var m fileExtraMeta
	if err := json.Unmarshal([]byte(extra), &m); err != nil {
		return ""
	}
	return strings.TrimSpace(m.HlsPid)
}

func SetHlsPid(extra, pid string) string {
	var m fileExtraMeta
	if extra != "" {
		_ = json.Unmarshal([]byte(extra), &m)
	}
	m.HlsPid = strings.TrimSpace(pid)
	b, err := json.Marshal(m)
	if err != nil {
		return extra
	}
	return string(b)
}

func ParseOriginAtFromExifField(exif string) int64 {
	if exif == "" {
		return 0
	}
	var m struct {
		OriginAt int64 `json:"origin_at"`
	}
	if err := json.Unmarshal([]byte(exif), &m); err != nil {
		return 0
	}
	return m.OriginAt
}
