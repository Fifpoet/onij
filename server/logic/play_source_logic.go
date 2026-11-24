package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"onij/biz/prm"
)

const (
	playSourcesFileName  = "play_sources.json"
	gdstudioSignFileName = "gdstudio_sign.json"
	minCrc32Runes        = 1024
	maxCrc32Bytes        = 512 * 1024
)

var gdstudioVersionRe = regexp.MustCompile(`^\d+(\.\d+)+$`)

type PlaySourceLogic interface {
	List(ctx context.Context, param *prm.ListPlaySourceParam) (*prm.PlaySourceListResult, error)
	Save(ctx context.Context, param *prm.SavePlaySourceParam) (*prm.PlaySourceListResult, error)
	GetSign(ctx context.Context, param *prm.GetGdstudioSignParam) (*prm.GdstudioSignResult, error)
	SaveSign(ctx context.Context, param *prm.SaveGdstudioSignParam) (*prm.GdstudioSignResult, error)
}

type playSourceFile struct {
	Sources []*prm.PlaySourceDTO `json:"sources"`
}

type playSourceLogic struct {
	mu   sync.Mutex
	path string
}

func NewPlaySourceLogic() PlaySourceLogic {
	return newPlaySourceLogic(playSourcesFileName)
}

func newPlaySourceLogic(path string) *playSourceLogic {
	return &playSourceLogic{path: path}
}

func builtinPlaySources() []*prm.PlaySourceDTO {
	return []*prm.PlaySourceDTO{
		{
			Id:      "gdapi",
			Title:   "GD 公开 API",
			Kind:    "gdapi",
			Origin:  "https://music-api.gdstudio.xyz",
			Enabled: true,
			Builtin: true,
		},
		{
			Id:      "gdstudio",
			Title:   "GD 播放器",
			Kind:    "gdstudio",
			Origin:  "",
			Enabled: true,
			Builtin: true,
		},
	}
}

func (l *playSourceLogic) List(_ context.Context, _ *prm.ListPlaySourceParam) (*prm.PlaySourceListResult, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return &prm.PlaySourceListResult{Sources: mergePlaySources(l.readUnlocked())}, nil
}

func (l *playSourceLogic) Save(_ context.Context, param *prm.SavePlaySourceParam) (*prm.PlaySourceListResult, error) {
	var incoming []*prm.PlaySourceDTO
	if param != nil {
		incoming = param.Sources
	}
	if err := validatePlaySources(incoming); err != nil {
		return nil, err
	}
	merged := mergePlaySources(incoming)

	l.mu.Lock()
	defer l.mu.Unlock()
	if err := l.writeUnlocked(merged); err != nil {
		return nil, err
	}
	return &prm.PlaySourceListResult{Sources: clonePlaySources(merged)}, nil
}

func (l *playSourceLogic) readUnlocked() []*prm.PlaySourceDTO {
	raw, err := os.ReadFile(l.path)
	if err != nil {
		return nil
	}
	var file playSourceFile
	if err := json.Unmarshal(raw, &file); err != nil || len(file.Sources) == 0 {
		return nil
	}
	return file.Sources
}

func (l *playSourceLogic) writeUnlocked(sources []*prm.PlaySourceDTO) error {
	return l.writeJSONFile(l.path, playSourceFile{Sources: sources})
}

func validatePlaySources(sources []*prm.PlaySourceDTO) error {
	seen := make(map[string]struct{}, len(sources))
	for _, src := range sources {
		if src == nil {
			return fmt.Errorf("source is empty")
		}
		id := strings.TrimSpace(src.Id)
		if id == "" {
			return fmt.Errorf("id is empty")
		}
		if _, ok := seen[id]; ok {
			return fmt.Errorf("duplicate id")
		}
		seen[id] = struct{}{}
		kind := strings.TrimSpace(src.Kind)
		if kind != "gdstudio" && kind != "gdapi" {
			return fmt.Errorf("kind is invalid")
		}
	}
	return nil
}

func mergePlaySources(saved []*prm.PlaySourceDTO) []*prm.PlaySourceDTO {
	byID := make(map[string]*prm.PlaySourceDTO, len(saved))
	order := make([]string, 0, len(saved))
	for _, src := range saved {
		if src == nil {
			continue
		}
		id := strings.TrimSpace(src.Id)
		if id == "" {
			continue
		}
		if _, exists := byID[id]; exists {
			continue
		}
		byID[id] = src
		order = append(order, id)
	}

	out := make([]*prm.PlaySourceDTO, 0, len(builtinPlaySources())+len(saved))
	used := make(map[string]struct{})
	for _, def := range builtinPlaySources() {
		item := clonePlaySource(def)
		if prev, ok := byID[def.Id]; ok {
			if title := strings.TrimSpace(prev.Title); title != "" {
				item.Title = title
			}
			item.Origin = strings.TrimSpace(prev.Origin)
			item.Enabled = prev.Enabled
		}
		out = append(out, item)
		used[def.Id] = struct{}{}
	}
	for _, id := range order {
		if _, keep := used[id]; keep {
			continue
		}
		prev := byID[id]
		kind := strings.TrimSpace(prev.Kind)
		if kind != "gdstudio" && kind != "gdapi" {
			continue
		}
		title := strings.TrimSpace(prev.Title)
		if title == "" {
			title = "自定义源"
		}
		out = append(out, &prm.PlaySourceDTO{
			Id:      id,
			Title:   title,
			Kind:    kind,
			Origin:  strings.TrimSpace(prev.Origin),
			Enabled: prev.Enabled,
			Builtin: false,
		})
	}
	return out
}

func clonePlaySources(in []*prm.PlaySourceDTO) []*prm.PlaySourceDTO {
	out := make([]*prm.PlaySourceDTO, 0, len(in))
	for _, src := range in {
		out = append(out, clonePlaySource(src))
	}
	return out
}

func clonePlaySource(src *prm.PlaySourceDTO) *prm.PlaySourceDTO {
	if src == nil {
		return nil
	}
	cp := *src
	return &cp
}

func (l *playSourceLogic) signPath() string {
	dir := filepath.Dir(l.path)
	if dir == "." || dir == "" {
		return gdstudioSignFileName
	}
	return filepath.Join(dir, gdstudioSignFileName)
}

func (l *playSourceLogic) GetSign(_ context.Context, _ *prm.GetGdstudioSignParam) (*prm.GdstudioSignResult, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return &prm.GdstudioSignResult{Sign: l.readSignUnlocked()}, nil
}

func (l *playSourceLogic) SaveSign(_ context.Context, param *prm.SaveGdstudioSignParam) (*prm.GdstudioSignResult, error) {
	version := ""
	crc32 := ""
	if param != nil {
		version = strings.TrimSpace(param.Version)
		crc32 = strings.TrimSpace(param.Crc32)
	}
	if err := validateGdstudioSign(version, crc32); err != nil {
		return nil, err
	}
	dto := &prm.GdstudioSignDTO{
		Version:   version,
		Crc32:     crc32,
		UpdatedAt: time.Now().Unix(),
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := l.writeJSONFile(l.signPath(), dto); err != nil {
		return nil, err
	}
	cp := *dto
	return &prm.GdstudioSignResult{Sign: &cp}, nil
}

func (l *playSourceLogic) readSignUnlocked() *prm.GdstudioSignDTO {
	raw, err := os.ReadFile(l.signPath())
	if err != nil {
		return &prm.GdstudioSignDTO{}
	}
	var dto prm.GdstudioSignDTO
	if err := json.Unmarshal(raw, &dto); err != nil {
		return &prm.GdstudioSignDTO{}
	}
	if strings.TrimSpace(dto.Version) == "" || strings.TrimSpace(dto.Crc32) == "" {
		return &prm.GdstudioSignDTO{}
	}
	return &dto
}

func (l *playSourceLogic) writeJSONFile(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(path)
		if err2 := os.Rename(tmp, path); err2 != nil {
			_ = os.Remove(tmp)
			return err2
		}
	}
	return nil
}

func validateGdstudioSign(version, crc32 string) error {
	if !gdstudioVersionRe.MatchString(version) {
		return fmt.Errorf("version is invalid")
	}
	if utf8.RuneCountInString(crc32) < minCrc32Runes || len(crc32) > maxCrc32Bytes {
		return fmt.Errorf("crc32 size is invalid")
	}
	lower := strings.ToLower(crc32)
	if strings.Contains(lower, "<html") || strings.Contains(lower, "<!doctype") {
		return fmt.Errorf("crc32 is invalid")
	}
	return nil
}
