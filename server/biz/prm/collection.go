package prm

import (
	"onij/model"
	"onij/util"
	"strings"
)

type CollectionDTO struct {
	Id        int64   `json:"id,string"`
	Name      string  `json:"name"`
	CoverUrl  string  `json:"cover_url,omitempty"`
	SongIds   []int64 `json:"song_ids,omitempty"`
	SongCount int     `json:"song_count"`
	UpdatedAt int64   `json:"updated_at"`
}

type CollectionSongMeta struct {
	SongId      int64  `json:"song_id"`
	SongName    string `json:"song_name"`
	ArtistNames string `json:"artist_names"`
}

type ListCollectionParam struct {
	Keyword string
}

type ListCollectionResult struct {
	Items []*CollectionDTO
}

func (r *ListCollectionResult) Resp() map[string]any {
	return map[string]any{
		"code":    util.BaseCodeOK,
		"message": util.BaseMsgOK,
		"items":   r.Items,
	}
}

type CreateCollectionParam struct {
	Name     string
	CoverUrl string
}

type UpdateCollectionParam struct {
	Id       int64
	Name     *string
	CoverUrl *string
}

type CollectionIdParam struct {
	Id int64
}

type CollectionDetailResult struct {
	Item *CollectionDTO
}

func (r *CollectionDetailResult) Resp() map[string]any {
	return map[string]any{
		"code":    util.BaseCodeOK,
		"message": util.BaseMsgOK,
		"item":    r.Item,
	}
}

type MutationCollectionResult struct {
	Item *CollectionDTO
}

func (r *MutationCollectionResult) Resp() map[string]any {
	return map[string]any{
		"code":    util.BaseCodeOK,
		"message": util.BaseMsgOK,
		"item":    r.Item,
	}
}

type AddCollectionSongsParam struct {
	Id    int64
	Songs []CollectionSongMeta
}

type RemoveCollectionSongParam struct {
	Id     int64
	SongId int64
}

type OkCollectionResult struct{}

func (r *OkCollectionResult) Resp() map[string]any {
	return map[string]any{
		"code":    util.BaseCodeOK,
		"message": util.BaseMsgOK,
	}
}

func ToCollectionDTO(c *model.MusicCollection, songIds []int64) *CollectionDTO {
	if c == nil {
		return nil
	}
	if songIds == nil {
		songIds = []int64{}
	}
	return &CollectionDTO{
		Id:        c.Id,
		Name:      c.Name,
		CoverUrl:  c.CoverUrl,
		SongIds:   songIds,
		SongCount: len(songIds),
		UpdatedAt: c.UpdatedAt.UnixMilli(),
	}
}

func BuildCollectionSearchText(name string, items []*model.MusicCollectionItem) string {
	parts := make([]string, 0, 1+len(items)*2)
	name = strings.TrimSpace(name)
	if name != "" {
		parts = append(parts, name)
	}
	for _, it := range items {
		if it == nil {
			continue
		}
		if s := strings.TrimSpace(it.SongName); s != "" {
			parts = append(parts, s)
		}
		if s := strings.TrimSpace(it.ArtistNames); s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, " ")
}
