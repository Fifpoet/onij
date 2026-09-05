package logic

import "testing"

func TestShouldUseHLS(t *testing.T) {
	const (
		min20 = 20 * 60 * 1000
		mb512 = 512 * 1024 * 1024
	)
	cases := []struct {
		name       string
		durationMs int64
		size       int64
		want       bool
	}{
		{"short small", 3 * 60 * 1000, 8 << 20, false},
		{"exactly 20min", min20, 8 << 20, false},
		{"just over 20min", min20 + 1, 8 << 20, true},
		{"exactly 512MB", 60 * 1000, mb512, false},
		{"just over 512MB", 60 * 1000, mb512 + 1, true},
		{"both over", min20 + 1, mb512 + 1, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := shouldUseHLS(c.durationMs, c.size); got != c.want {
				t.Fatalf("shouldUseHLS(%d, %d) = %v, want %v", c.durationMs, c.size, got, c.want)
			}
		})
	}
}

func TestIsMp4Name(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"lower", "concert.mp4", true},
		{"upper", "CONCERT.MP4", true},
		{"mixed", "Clip.Mp4", true},
		{"mkv", "concert.mkv", false},
		{"no ext", "concert", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isMp4Name(c.in); got != c.want {
				t.Fatalf("isMp4Name(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}
