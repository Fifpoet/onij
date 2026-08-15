package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func NeteaseAPIBase() string {
	if v := strings.TrimSpace(os.Getenv("NETEASE_API_BASE")); v != "" {
		return strings.TrimRight(v, "/")
	}
	return "http://onij.fun:3000"
}

type NeteaseSearchSong struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Duration int64  `json:"duration"`
	Artists  []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"artists"`
	Album struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"album"`
}

func NeteaseSearchSongs(query string, limit int) ([]map[string]any, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("query is empty")
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 30 {
		limit = 30
	}
	u := fmt.Sprintf("%s/search?keywords=%s&limit=%d&type=1",
		NeteaseAPIBase(), url.QueryEscape(query), limit)
	raw, err := neteaseGet(u)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Code   int `json:"code"`
		Result struct {
			Songs []NeteaseSearchSong `json:"songs"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(parsed.Result.Songs))
	for _, s := range parsed.Result.Songs {
		artists := make([]string, 0, len(s.Artists))
		for _, a := range s.Artists {
			artists = append(artists, a.Name)
		}
		out = append(out, map[string]any{
			"song_id":      s.ID,
			"song_name":    s.Name,
			"artist_names": strings.Join(artists, " / "),
			"album_name":   s.Album.Name,
			"duration_ms":  s.Duration,
		})
	}
	return out, nil
}

func NeteaseSongDetails(ids []int64) ([]map[string]any, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("song_ids is empty")
	}
	if len(ids) > 20 {
		ids = ids[:20]
	}
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		if id > 0 {
			parts = append(parts, strconv.FormatInt(id, 10))
		}
	}
	if len(parts) == 0 {
		return nil, fmt.Errorf("song_ids invalid")
	}
	u := fmt.Sprintf("%s/song/detail?ids=%s", NeteaseAPIBase(), url.QueryEscape(strings.Join(parts, ",")))
	raw, err := neteaseGet(u)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Code  int `json:"code"`
		Songs []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Dt   int64  `json:"dt"`
			Ar   []struct {
				ID   int64  `json:"id"`
				Name string `json:"name"`
			} `json:"ar"`
			Al struct {
				ID     int64  `json:"id"`
				Name   string `json:"name"`
				PicURL string `json:"picUrl"`
			} `json:"al"`
		} `json:"songs"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(parsed.Songs))
	for _, s := range parsed.Songs {
		artists := make([]map[string]any, 0, len(s.Ar))
		names := make([]string, 0, len(s.Ar))
		for _, a := range s.Ar {
			artists = append(artists, map[string]any{"artist_id": a.ID, "artist_name": a.Name})
			names = append(names, a.Name)
		}
		out = append(out, map[string]any{
			"song_id":      s.ID,
			"song_name":    s.Name,
			"artist_names": strings.Join(names, " / "),
			"artists":      artists,
			"album_id":     s.Al.ID,
			"album_name":   s.Al.Name,
			"cover_url":    s.Al.PicURL,
			"duration_ms":  s.Dt,
		})
	}
	return out, nil
}

func neteaseGet(u string) ([]byte, error) {
	cli := &http.Client{Timeout: 20 * time.Second}
	resp, err := cli.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("netease http %d: %s", resp.StatusCode, string(raw))
	}
	return raw, nil
}
