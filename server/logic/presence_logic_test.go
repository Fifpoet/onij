package logic

import (
	"context"
	"fmt"
	"testing"
	"time"

	"onij/biz/prm"
)

func TestPresenceHeartbeatExcludesSelf(t *testing.T) {
	l := NewPresenceLogic()
	ctx := context.Background()

	self, err := l.Heartbeat(ctx, &prm.HeartbeatParam{
		ClientId: "self",
		Nickname: "安静的磁带",
	})
	if err != nil {
		t.Fatalf("Heartbeat(self): %v", err)
	}
	if len(self.Peers) != 0 {
		t.Fatalf("alone: want 0 peers, got %d", len(self.Peers))
	}

	other, err := l.Heartbeat(ctx, &prm.HeartbeatParam{
		ClientId: "other",
		Nickname: "温柔的唱片",
	})
	if err != nil {
		t.Fatalf("Heartbeat(other): %v", err)
	}
	if len(other.Peers) != 1 || other.Peers[0].ClientId != "self" {
		t.Fatalf("other: want only self, got %+v", peerIDs(other))
	}

	again, err := l.Heartbeat(ctx, &prm.HeartbeatParam{
		ClientId: "self",
		Nickname: "安静的磁带",
	})
	if err != nil {
		t.Fatalf("Heartbeat(self again): %v", err)
	}
	if len(again.Peers) != 1 || again.Peers[0].ClientId != "other" {
		t.Fatalf("self: want only other, got %+v", peerIDs(again))
	}
}

func TestPresenceHeartbeatKeepsIdlePeer(t *testing.T) {
	l := NewPresenceLogic()
	ctx := context.Background()

	idleTrack := &prm.PresenceTrackDTO{Id: 7, Name: "未播这首"}
	if _, err := l.Heartbeat(ctx, &prm.HeartbeatParam{
		ClientId:   "idle",
		Nickname:   "沉默的耳机",
		NowPlaying: nil,
		Queue:      []*prm.PresenceTrackDTO{idleTrack},
		IsPlaying:  false,
	}); err != nil {
		t.Fatalf("Heartbeat(idle): %v", err)
	}

	viewer, err := l.Heartbeat(ctx, &prm.HeartbeatParam{
		ClientId: "viewer",
		Nickname: "快乐的晚风",
		NowPlaying: &prm.PresenceTrackDTO{
			Id:   1,
			Name: "正在听",
		},
		IsPlaying: true,
	})
	if err != nil {
		t.Fatalf("Heartbeat(viewer): %v", err)
	}
	if len(viewer.Peers) != 1 {
		t.Fatalf("want idle peer in list, got %+v", peerIDs(viewer))
	}
	p := viewer.Peers[0]
	if p.ClientId != "idle" {
		t.Fatalf("peer id = %q, want idle", p.ClientId)
	}
	if p.NowPlaying != nil {
		t.Fatalf("idle NowPlaying = %+v, want nil", p.NowPlaying)
	}
	if p.IsPlaying {
		t.Fatalf("idle IsPlaying = true, want false")
	}
	if len(p.Queue) != 1 || p.Queue[0].Id != 7 {
		t.Fatalf("idle queue = %+v", p.Queue)
	}
}

func TestPresenceHeartbeatCapsQueueAt80(t *testing.T) {
	l := NewPresenceLogic()
	ctx := context.Background()

	n := prm.PresenceQueueLimit + 1
	queue := make([]*prm.PresenceTrackDTO, n)
	for i := 0; i < n; i++ {
		queue[i] = &prm.PresenceTrackDTO{Id: int64(i + 1), Name: fmt.Sprintf("t%d", i+1)}
	}
	if _, err := l.Heartbeat(ctx, &prm.HeartbeatParam{
		ClientId: "long",
		Nickname: "辽阔的山谷",
		Queue:    queue,
	}); err != nil {
		t.Fatalf("Heartbeat(long): %v", err)
	}

	viewer, err := l.Heartbeat(ctx, &prm.HeartbeatParam{
		ClientId: "viewer",
		Nickname: "小小的口袋",
	})
	if err != nil {
		t.Fatalf("Heartbeat(viewer): %v", err)
	}
	if len(viewer.Peers) != 1 {
		t.Fatalf("want 1 peer, got %d", len(viewer.Peers))
	}
	got := viewer.Peers[0].Queue
	if len(got) != prm.PresenceQueueLimit {
		t.Fatalf("queue len = %d, want %d", len(got), prm.PresenceQueueLimit)
	}
	if got[0].Id != 1 || got[len(got)-1].Id != int64(prm.PresenceQueueLimit) {
		t.Fatalf("capped queue first=%d last=%d, want 1..%d", got[0].Id, got[len(got)-1].Id, prm.PresenceQueueLimit)
	}
}

func TestPresenceHeartbeatTTLSweep(t *testing.T) {
	l := newPresenceLogic()
	t0 := time.Date(2026, 10, 13, 19, 0, 0, 0, time.UTC)
	now := t0
	l.now = func() time.Time { return now }
	ctx := context.Background()

	hb := func(id, nick string) *prm.HeartbeatResult {
		t.Helper()
		r, err := l.Heartbeat(ctx, &prm.HeartbeatParam{ClientId: id, Nickname: nick})
		if err != nil {
			t.Fatalf("Heartbeat(%s): %v", id, err)
		}
		return r
	}

	hb("a", "安静的磁带")
	hb("b", "温柔的唱片")

	now = t0.Add(presenceTTL)
	atTTL := hb("a", "安静的磁带")
	if len(atTTL.Peers) != 1 || atTTL.Peers[0].ClientId != "b" {
		t.Fatalf("at TTL: want b still present, got %+v", peerIDs(atTTL))
	}

	now = t0.Add(presenceTTL + time.Nanosecond)
	after := hb("a", "安静的磁带")
	if len(after.Peers) != 0 {
		t.Fatalf("after TTL: want b gone, got %+v", peerIDs(after))
	}
}

func TestPresenceHeartbeatRejectsEmptyIdOrNickname(t *testing.T) {
	l := NewPresenceLogic()
	ctx := context.Background()
	cases := []struct {
		name string
		p    *prm.HeartbeatParam
	}{
		{"nil", nil},
		{"blank id", &prm.HeartbeatParam{ClientId: "  ", Nickname: "安静的磁带"}},
		{"blank nick", &prm.HeartbeatParam{ClientId: "self", Nickname: "  "}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := l.Heartbeat(ctx, c.p); err == nil {
				t.Fatalf("want error, got nil")
			}
		})
	}
}

func peerIDs(r *prm.HeartbeatResult) []string {
	if r == nil {
		return nil
	}
	ids := make([]string, len(r.Peers))
	for i, p := range r.Peers {
		ids[i] = p.ClientId
	}
	return ids
}
