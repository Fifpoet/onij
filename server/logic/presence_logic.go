package logic

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"onij/biz/prm"
)

const (
	presenceTTL         = 90 * time.Second
	presenceNicknameMax = 32
)

type PresenceLogic interface {
	Heartbeat(ctx context.Context, param *prm.HeartbeatParam) (*prm.HeartbeatResult, error)
}

type presenceEntry struct {
	clientId           string
	nickname           string
	nowPlaying         *prm.PresenceTrackDTO
	queue              []*prm.PresenceTrackDTO
	playbackCurrentSec float64
	isPlaying          bool
	lastSeen           time.Time
}

type presenceLogic struct {
	mu    sync.Mutex
	peers map[string]*presenceEntry
	ttl   time.Duration
	now   func() time.Time
}

func NewPresenceLogic() PresenceLogic {
	return newPresenceLogic()
}

func newPresenceLogic() *presenceLogic {
	return &presenceLogic{
		peers: make(map[string]*presenceEntry),
		ttl:   presenceTTL,
		now:   time.Now,
	}
}

func (l *presenceLogic) Heartbeat(_ context.Context, param *prm.HeartbeatParam) (*prm.HeartbeatResult, error) {
	if param == nil {
		return nil, fmt.Errorf("client_id is empty")
	}
	clientId := strings.TrimSpace(param.ClientId)
	if clientId == "" {
		return nil, fmt.Errorf("client_id is empty")
	}
	nickname := strings.TrimSpace(param.Nickname)
	if nickname == "" {
		return nil, fmt.Errorf("nickname is empty")
	}
	if utf8.RuneCountInString(nickname) > presenceNicknameMax {
		nickname = string([]rune(nickname)[:presenceNicknameMax])
	}

	now := l.now()
	queue := capPresenceQueue(param.Queue)
	sec := param.PlaybackCurrentSec
	if sec < 0 || !finiteFloat(sec) {
		sec = 0
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.peers[clientId] = &presenceEntry{
		clientId:           clientId,
		nickname:           nickname,
		nowPlaying:         param.NowPlaying,
		queue:              queue,
		playbackCurrentSec: sec,
		isPlaying:          param.IsPlaying,
		lastSeen:           now,
	}
	l.sweepExpiredLocked(now)

	others := make([]*prm.PresencePeerDTO, 0, len(l.peers))
	for id, e := range l.peers {
		if id == clientId {
			continue
		}
		others = append(others, entryToPeerDTO(e))
	}
	sort.SliceStable(others, func(i, j int) bool {
		if others[i].Nickname != others[j].Nickname {
			return others[i].Nickname < others[j].Nickname
		}
		return others[i].ClientId < others[j].ClientId
	})
	return &prm.HeartbeatResult{Peers: others}, nil
}

func (l *presenceLogic) sweepExpiredLocked(now time.Time) {
	for id, e := range l.peers {
		if now.Sub(e.lastSeen) > l.ttl {
			delete(l.peers, id)
		}
	}
}

func capPresenceQueue(q []*prm.PresenceTrackDTO) []*prm.PresenceTrackDTO {
	if len(q) == 0 {
		return []*prm.PresenceTrackDTO{}
	}
	if len(q) > prm.PresenceQueueLimit {
		q = q[:prm.PresenceQueueLimit]
	}
	out := make([]*prm.PresenceTrackDTO, len(q))
	copy(out, q)
	return out
}

func entryToPeerDTO(e *presenceEntry) *prm.PresencePeerDTO {
	queue := e.queue
	if queue == nil {
		queue = []*prm.PresenceTrackDTO{}
	}
	return &prm.PresencePeerDTO{
		ClientId:           e.clientId,
		Nickname:           e.nickname,
		NowPlaying:         e.nowPlaying,
		Queue:              queue,
		PlaybackCurrentSec: e.playbackCurrentSec,
		IsPlaying:          e.isPlaying,
	}
}

func finiteFloat(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
