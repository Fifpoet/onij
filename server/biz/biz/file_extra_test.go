package biz

import "testing"

func TestSetHlsPidKeepsOriginAt(t *testing.T) {
	const origin int64 = 1693800000
	extra := BuildFileExtra(origin, "raw-exif")
	if extra == "" {
		t.Fatal("BuildFileExtra returned empty")
	}
	if ParseFileOriginAt(extra) != origin {
		t.Fatalf("origin_at before set: %d", ParseFileOriginAt(extra))
	}
	if ParseHlsPid(extra) != "" {
		t.Fatalf("hls_pid should be empty initially, got %q", ParseHlsPid(extra))
	}

	updated := SetHlsPid(extra, " pid-abc ")
	if ParseHlsPid(updated) != "pid-abc" {
		t.Fatalf("ParseHlsPid = %q", ParseHlsPid(updated))
	}
	if ParseFileOriginAt(updated) != origin {
		t.Fatalf("origin_at lost after SetHlsPid: %d", ParseFileOriginAt(updated))
	}
}

func TestParseHlsPidEmptyAndOverwrite(t *testing.T) {
	if ParseHlsPid("") != "" {
		t.Fatal("empty extra should yield empty pid")
	}
	once := SetHlsPid("", "first")
	twice := SetHlsPid(once, "second")
	if ParseHlsPid(twice) != "second" {
		t.Fatalf("overwrite pid = %q", ParseHlsPid(twice))
	}
	if ParseFileOriginAt(twice) != 0 {
		t.Fatalf("empty origin should stay 0, got %d", ParseFileOriginAt(twice))
	}
}
