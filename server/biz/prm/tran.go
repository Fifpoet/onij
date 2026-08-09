package prm

import (
	"onij/model"
	"onij/util"
	"sort"
	"time"
)

type TranItemDTO struct {
	Id        int64  `json:"id,string"`
	Kind      string `json:"kind"`
	Content   string `json:"content,omitempty"`
	FileId    int64  `json:"file_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Size      int64  `json:"size,omitempty"`
	Format    int32  `json:"format,omitempty"`
	Pinned    bool   `json:"pinned"`
	PinnedAt  int64  `json:"pinned_at,omitempty"`
	CreatedAt int64  `json:"created_at"`
}

type ListTranParam struct{}

type ListTranResult struct {
	Items []*TranItemDTO
}

func (r *ListTranResult) Resp() map[string]any {
	return map[string]any{
		"code":    util.BaseCodeOK,
		"message": util.BaseMsgOK,
		"items":   r.Items,
	}
}

type AddTranTextParam struct {
	Content string
}

type AddTranFileParam struct {
	FileId int64
	Name   string
	Size   int64
	Format int32
}

type AddTranResult struct {
	Item *TranItemDTO
}

func (r *AddTranResult) Resp() map[string]any {
	return map[string]any{
		"code":    util.BaseCodeOK,
		"message": util.BaseMsgOK,
		"item":    r.Item,
	}
}

type PinTranParam struct {
	Id int64
}

type PinTranResult struct {
	Item *TranItemDTO
}

func (r *PinTranResult) Resp() map[string]any {
	return map[string]any{
		"code":    util.BaseCodeOK,
		"message": util.BaseMsgOK,
		"item":    r.Item,
	}
}

type DeleteTranParam struct {
	Id int64
}

type DeleteTranResult struct{}

func (r *DeleteTranResult) Resp() map[string]any {
	return map[string]any{
		"code":    util.BaseCodeOK,
		"message": util.BaseMsgOK,
	}
}

func ToTranItemDTO(it *model.TranItem) *TranItemDTO {
	if it == nil {
		return nil
	}
	kind := "text"
	if it.Kind == model.TranKindFile {
		kind = "file"
	}
	createdAt := it.CreatedAt.UnixMilli()
	if createdAt <= 0 {
		createdAt = time.Now().UnixMilli()
	}
	dto := &TranItemDTO{
		Id:        it.Id,
		Kind:      kind,
		Pinned:    it.Pinned,
		PinnedAt:  it.PinnedAt,
		CreatedAt: createdAt,
	}
	if kind == "text" {
		dto.Content = it.Content
	} else {
		dto.FileId = it.FileId
		dto.Name = it.Name
		dto.Size = it.Size
		dto.Format = it.Format
	}
	return dto
}

func SortTranDTOs(items []*TranItemDTO) {
	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i], items[j]
		if a.Pinned != b.Pinned {
			return a.Pinned
		}
		if a.Pinned {
			return a.PinnedAt > b.PinnedAt
		}
		return a.CreatedAt > b.CreatedAt
	})
}
