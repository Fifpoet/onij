package prm

import "onij/util"

const PresenceQueueLimit = 80

type PresenceArtistDTO struct {
	ArtistId   int64  `json:"artist_id"`
	ArtistName string `json:"artist_name"`
}

type PresenceTrackDTO struct {
	Id             int64               `json:"id"`
	Name           string              `json:"name"`
	Artists        []PresenceArtistDTO `json:"artists"`
	DisplayArtists []PresenceArtistDTO `json:"display_artists,omitempty"`
	TimeLong       int64               `json:"time_long"`
	ComposerId     int64               `json:"composer_id,omitempty"`
	ComposerName   string              `json:"composer_name,omitempty"`
	WriterId       int64               `json:"writer_id,omitempty"`
	WriterName     string              `json:"writer_name,omitempty"`
	IssueTime      int64               `json:"issue_time,omitempty"`
	MvUrl          string              `json:"mv_url,omitempty"`
	AlbumId        int64               `json:"album_id"`
	AlbumName      string              `json:"album_name"`
	CoverFileUrl   string              `json:"cover_file_url,omitempty"`
	Tags           []string            `json:"tags,omitempty"`
}

type PresencePeerDTO struct {
	ClientId           string              `json:"client_id"`
	Nickname           string              `json:"nickname"`
	NowPlaying         *PresenceTrackDTO   `json:"now_playing"`
	Queue              []*PresenceTrackDTO `json:"queue"`
	PlaybackCurrentSec float64             `json:"playback_current_sec"`
	IsPlaying          bool                `json:"is_playing"`
}

type HeartbeatParam struct {
	ClientId           string
	Nickname           string
	NowPlaying         *PresenceTrackDTO
	Queue              []*PresenceTrackDTO
	PlaybackCurrentSec float64
	IsPlaying          bool
}

type HeartbeatResult struct {
	Peers []*PresencePeerDTO
}

func (r *HeartbeatResult) Resp() map[string]any {
	peers := r.Peers
	if peers == nil {
		peers = []*PresencePeerDTO{}
	}
	return map[string]any{
		"code":    util.BaseCodeOK,
		"message": util.BaseMsgOK,
		"peers":   peers,
	}
}
