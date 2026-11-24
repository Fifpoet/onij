package logic

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"onij/biz/prm"
)

func TestPlaySourceListUsesBuiltinsWhenFileMissing(t *testing.T) {
	dir := t.TempDir()
	l := newPlaySourceLogic(filepath.Join(dir, "play_sources.json"))
	got, err := l.List(context.Background(), &prm.ListPlaySourceParam{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	assertBuiltinDefaults(t, got.Sources)
}

func TestPlaySourceListUsesBuiltinsWhenFileBadOrEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "play_sources.json")
	if err := os.WriteFile(path, []byte("{not json"), 0644); err != nil {
		t.Fatal(err)
	}
	l := newPlaySourceLogic(path)
	got, err := l.List(context.Background(), &prm.ListPlaySourceParam{})
	if err != nil {
		t.Fatalf("List bad json: %v", err)
	}
	assertBuiltinDefaults(t, got.Sources)

	if err := os.WriteFile(path, []byte(`{"sources":[]}`), 0644); err != nil {
		t.Fatal(err)
	}
	got, err = l.List(context.Background(), &prm.ListPlaySourceParam{})
	if err != nil {
		t.Fatalf("List empty: %v", err)
	}
	assertBuiltinDefaults(t, got.Sources)

	if err := os.WriteFile(path, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	got, err = l.List(context.Background(), &prm.ListPlaySourceParam{})
	if err != nil {
		t.Fatalf("List empty file: %v", err)
	}
	assertBuiltinDefaults(t, got.Sources)
}

func TestPlaySourceSaveRejectsInvalidKindEmptyAndDupID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "play_sources.json")
	l := newPlaySourceLogic(path)

	_, err := l.Save(context.Background(), &prm.SavePlaySourceParam{
		Sources: []*prm.PlaySourceDTO{{Id: "x", Title: "bad", Kind: "meting", Enabled: true}},
	})
	if err == nil {
		t.Fatal("want invalid kind error")
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Fatalf("invalid kind must not write file: %v", statErr)
	}

	_, err = l.Save(context.Background(), &prm.SavePlaySourceParam{
		Sources: []*prm.PlaySourceDTO{{Id: "  ", Title: "x", Kind: "gdapi"}},
	})
	if err == nil {
		t.Fatal("want empty id error")
	}

	_, err = l.Save(context.Background(), &prm.SavePlaySourceParam{
		Sources: []*prm.PlaySourceDTO{
			{Id: "a", Title: "one", Kind: "gdapi"},
			{Id: "a", Title: "two", Kind: "gdstudio"},
		},
	})
	if err == nil {
		t.Fatal("want duplicate id error")
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Fatalf("validation error must not write file: %v", statErr)
	}
}

func TestPlaySourceSaveThenListKeepsOverridesAndExtras(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "play_sources.json")
	l := newPlaySourceLogic(path)

	saved, err := l.Save(context.Background(), &prm.SavePlaySourceParam{
		Sources: []*prm.PlaySourceDTO{
			{Id: "gdapi", Title: "公开源", Kind: "gdapi", Origin: "https://example.com/gdapi", Enabled: false},
			{Id: "gdstudio", Title: "播放器", Kind: "gdstudio", Origin: "", Enabled: true},
			{Id: "custom-1", Title: "备用", Kind: "gdstudio", Origin: "https://music.example", Enabled: true},
		},
	})
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if len(saved.Sources) != 3 {
		t.Fatalf("save result len=%d want 3", len(saved.Sources))
	}
	if saved.Sources[0].Title != "公开源" || saved.Sources[0].Enabled || !saved.Sources[0].Builtin {
		t.Fatalf("gdapi override: %+v", saved.Sources[0])
	}
	if saved.Sources[2].Id != "custom-1" || saved.Sources[2].Builtin {
		t.Fatalf("extra: %+v", saved.Sources[2])
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var file playSourceFile
	if err := json.Unmarshal(raw, &file); err != nil {
		t.Fatal(err)
	}
	if len(file.Sources) != 3 {
		t.Fatalf("file sources=%d", len(file.Sources))
	}

	got, err := l.List(context.Background(), &prm.ListPlaySourceParam{})
	if err != nil {
		t.Fatalf("List after save: %v", err)
	}
	if len(got.Sources) != 3 {
		t.Fatalf("list len=%d want 3", len(got.Sources))
	}
	if got.Sources[0].Origin != "https://example.com/gdapi" || got.Sources[0].Enabled {
		t.Fatalf("list gdapi: %+v", got.Sources[0])
	}
	if got.Sources[1].Id != "gdstudio" || !got.Sources[1].Builtin {
		t.Fatalf("list gdstudio: %+v", got.Sources[1])
	}
	if got.Sources[2].Kind != "gdstudio" || got.Sources[2].Title != "备用" {
		t.Fatalf("list extra: %+v", got.Sources[2])
	}
}

func TestPlaySourceMergeAlwaysKeepsBuiltins(t *testing.T) {
	dir := t.TempDir()
	l := newPlaySourceLogic(filepath.Join(dir, "play_sources.json"))
	got, err := l.Save(context.Background(), &prm.SavePlaySourceParam{
		Sources: []*prm.PlaySourceDTO{
			{Id: "only-extra", Title: "extra", Kind: "gdapi", Origin: "https://x", Enabled: true},
		},
	})
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if len(got.Sources) != 3 || got.Sources[0].Id != "gdapi" || got.Sources[1].Id != "gdstudio" || got.Sources[2].Id != "only-extra" {
		t.Fatalf("merge order: %+v", idsOf(got.Sources))
	}

	empty, err := l.Save(context.Background(), &prm.SavePlaySourceParam{Sources: nil})
	if err != nil {
		t.Fatalf("Save empty: %v", err)
	}
	assertBuiltinDefaults(t, empty.Sources)
}

func TestPlaySourceListSkipsInvalidExtraFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "play_sources.json")
	raw := []byte(`{"sources":[{"id":"gdapi","title":"公开源","kind":"gdstudio","origin":"https://example.com/gdapi","enabled":false},{"id":"old-meting","title":"Meting","kind":"meting","origin":"https://api.qijieya.cn/meting","enabled":true}]}`)
	if err := os.WriteFile(path, raw, 0644); err != nil {
		t.Fatal(err)
	}
	l := newPlaySourceLogic(path)
	got, err := l.List(context.Background(), &prm.ListPlaySourceParam{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Sources) != 2 {
		t.Fatalf("len=%d want 2 builtins only, ids=%v", len(got.Sources), idsOf(got.Sources))
	}
	if got.Sources[0].Id != "gdapi" || got.Sources[0].Kind != "gdapi" || got.Sources[0].Title != "公开源" || got.Sources[0].Enabled {
		t.Fatalf("gdapi override kept, kind not overwritten: %+v", got.Sources[0])
	}
	if got.Sources[1].Id != "gdstudio" || !got.Sources[1].Builtin {
		t.Fatalf("gdstudio: %+v", got.Sources[1])
	}
}

func assertBuiltinDefaults(t *testing.T, sources []*prm.PlaySourceDTO) {
	t.Helper()
	if len(sources) != 2 {
		t.Fatalf("len=%d want 2", len(sources))
	}
	if sources[0].Id != "gdapi" || sources[0].Kind != "gdapi" || !sources[0].Builtin || sources[0].Origin != "https://music-api.gdstudio.xyz" {
		t.Fatalf("gdapi: %+v", sources[0])
	}
	if sources[1].Id != "gdstudio" || sources[1].Kind != "gdstudio" || !sources[1].Builtin || sources[1].Origin != "" {
		t.Fatalf("gdstudio: %+v", sources[1])
	}
}

func idsOf(sources []*prm.PlaySourceDTO) []string {
	out := make([]string, 0, len(sources))
	for _, s := range sources {
		out = append(out, s.Id)
	}
	return out
}

func TestGdstudioSignGetEmptyWhenMissing(t *testing.T) {
	dir := t.TempDir()
	l := newPlaySourceLogic(filepath.Join(dir, "play_sources.json"))
	got, err := l.GetSign(context.Background(), &prm.GetGdstudioSignParam{})
	if err != nil {
		t.Fatalf("GetSign: %v", err)
	}
	if got.Sign == nil || got.Sign.Version != "" || got.Sign.Crc32 != "" {
		t.Fatalf("empty sign: %+v", got.Sign)
	}
}

func TestGdstudioSignSaveRejectsInvalid(t *testing.T) {
	dir := t.TempDir()
	l := newPlaySourceLogic(filepath.Join(dir, "play_sources.json"))
	crc32 := strings.Repeat("a", 2048)

	_, err := l.SaveSign(context.Background(), &prm.SaveGdstudioSignParam{Version: "bad", Crc32: crc32})
	if err == nil {
		t.Fatal("want invalid version")
	}
	_, err = l.SaveSign(context.Background(), &prm.SaveGdstudioSignParam{Version: "2026.11.22", Crc32: "short"})
	if err == nil {
		t.Fatal("want invalid crc32 size")
	}
	_, err = l.SaveSign(context.Background(), &prm.SaveGdstudioSignParam{
		Version: "2026.11.22",
		Crc32:   "<html>" + crc32,
	})
	if err == nil {
		t.Fatal("want html crc32 rejected")
	}
	if _, statErr := os.Stat(filepath.Join(dir, "gdstudio_sign.json")); !os.IsNotExist(statErr) {
		t.Fatalf("invalid sign must not write file: %v", statErr)
	}
}

func TestGdstudioSignSaveThenGet(t *testing.T) {
	dir := t.TempDir()
	l := newPlaySourceLogic(filepath.Join(dir, "play_sources.json"))
	crc32 := strings.Repeat("md5fn", 400)
	saved, err := l.SaveSign(context.Background(), &prm.SaveGdstudioSignParam{
		Version: "2026.11.22",
		Crc32:   crc32,
	})
	if err != nil {
		t.Fatalf("SaveSign: %v", err)
	}
	if saved.Sign.Version != "2026.11.22" || saved.Sign.Crc32 != crc32 || saved.Sign.UpdatedAt == 0 {
		t.Fatalf("saved: %+v", saved.Sign)
	}

	got, err := l.GetSign(context.Background(), &prm.GetGdstudioSignParam{})
	if err != nil {
		t.Fatalf("GetSign: %v", err)
	}
	if got.Sign.Version != "2026.11.22" || got.Sign.Crc32 != crc32 {
		t.Fatalf("get: %+v", got.Sign)
	}
}
