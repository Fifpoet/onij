package util

import (
	"net/url"
	"strings"
	"testing"
)

func TestRewriteM3U8(t *testing.T) {
	t.Setenv("QINIU_SK", "unit-test-sk-not-real")
	t.Setenv("QINIU_PUBLIC_DOMAIN", "http://cloud.onij.fun")

	playlist := strings.Join([]string{
		"#EXTM3U",
		"#EXT-X-VERSION:3",
		"#EXTINF:10.0,",
		"seg0.ts",
		"#EXTINF:10.0,",
		"http://cloud.onij.fun/videos/concert/seg1.ts",
		"#EXT-X-ENDLIST",
		"",
	}, "\n")

	got := RewriteM3U8(playlist, "videos/concert/index.m3u8")
	lines := strings.Split(got, "\n")

	if lines[0] != "#EXTM3U" || lines[1] != "#EXT-X-VERSION:3" || lines[2] != "#EXTINF:10.0," {
		t.Fatalf("comment/header lines mutated: %q", lines[:3])
	}
	if lines[4] != "#EXTINF:10.0," || lines[6] != "#EXT-X-ENDLIST" {
		t.Fatalf("remaining comment lines mutated: %v", lines)
	}

	assertSignedCloudURL(t, "relative ts", lines[3], "videos/concert/seg0.ts")
	assertSignedCloudURL(t, "absolute ts", lines[5], "videos/concert/seg1.ts")
}

func TestRewriteM3U8TrimsLeadingSlashes(t *testing.T) {
	t.Setenv("QINIU_SK", "unit-test-sk-not-real")
	t.Setenv("QINIU_PUBLIC_DOMAIN", "http://cloud.onij.fun")

	got := RewriteM3U8("http://iovip.qbox.me//abc/000000.ts\n", "cloud/live/a.mp4.m3u8")
	assertSignedCloudURL(t, "iodomain ts", strings.TrimSpace(got), "abc/000000.ts")
}

func TestIsQiniuErrorJSON(t *testing.T) {
	if !isQiniuErrorJSON([]byte(`{"error":"no such domain"}`)) {
		t.Fatal("should detect qiniu error json")
	}
	if isQiniuErrorJSON([]byte("#EXTM3U\nseg0.ts\n")) {
		t.Fatal("m3u8 is not error json")
	}
	if isQiniuErrorJSON([]byte(`{"format":{"duration":"1.0"}}`)) {
		t.Fatal("avinfo json without error field is not an error body")
	}
}

func TestIsM3U8Playlist(t *testing.T) {
	if !isM3U8Playlist("\n#EXTM3U\n#EXTINF:10,\nseg.ts\n") {
		t.Fatal("want m3u8")
	}
	if isM3U8Playlist(`{"error":"no such domain"}`) {
		t.Fatal("error json must not pass as m3u8")
	}
}

func TestRewriteM3U8IgnoresErrorJSON(t *testing.T) {
	t.Setenv("QINIU_SK", "unit-test-sk-not-real")
	t.Setenv("QINIU_PUBLIC_DOMAIN", "http://cloud.onij.fun")

	got := RewriteM3U8(`{"error":"no such domain"}`, "cloud/live/concert.mp4.m3u8")
	if strings.Contains(got, "token=") || strings.Contains(got, "cloud/live/%7B") {
		t.Fatalf("error json should not become a signed ts url: %q", got)
	}
}

func TestHLSKey(t *testing.T) {
	if got := HLSKey("videos/a.mp4"); got != "ts/videos/a.mp4.m3u8" {
		t.Fatalf("HLSKey = %q", got)
	}
}

func TestIsHlsHashRootPrefix(t *testing.T) {
	if !IsHlsHashRootPrefix("cmIt1FPt0ihZ7qw8ztQZlSOPRco=") {
		t.Fatal("want hash prefix")
	}
	if IsHlsHashRootPrefix("cloud") || IsHlsHashRootPrefix("ts") || IsHlsHashRootPrefix("tran") {
		t.Fatal("app folders must be kept")
	}
}

func TestIsLegacyHlsPlaylistKey(t *testing.T) {
	if !IsLegacyHlsPlaylistKey("cloud/live/a.mp4.m3u8") {
		t.Fatal("old sidecar m3u8")
	}
	if IsLegacyHlsPlaylistKey("ts/cloud/live/a.mp4.m3u8") {
		t.Fatal("new ts/ playlist is not legacy")
	}
}

func TestPrefopStatus(t *testing.T) {
	if PrefopOK(0) != true || PrefopBusy(0) {
		t.Fatal("code 0 should be OK, not busy")
	}
	if !PrefopBusy(1) || !PrefopBusy(2) || PrefopOK(1) {
		t.Fatal("code 1/2 should be busy")
	}
}

func assertSignedCloudURL(t *testing.T, name, raw, wantKey string) {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("%s: parse %q: %v", name, raw, err)
	}
	if u.Scheme != "http" || u.Host != "cloud.onij.fun" {
		t.Fatalf("%s: want http://cloud.onij.fun host, got %s://%s", name, u.Scheme, u.Host)
	}
	if u.Host == "onij.fun" {
		t.Fatalf("%s: signed onto site host onij.fun: %q", name, raw)
	}
	gotKey := strings.TrimPrefix(u.Path, "/")
	if gotKey != wantKey {
		t.Fatalf("%s: key = %q, want %q", name, gotKey, wantKey)
	}
	if u.Query().Get("e") == "" || u.Query().Get("token") == "" {
		t.Fatalf("%s: missing private sign query: %q", name, raw)
	}
}
