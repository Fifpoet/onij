// 手写 API 结构体：新增接口优先在此定义，避免依赖 hz 重新生成路由/handler。
// 命名与 proto/pb 对齐（snake_case json tag），仅保留 handler 所需字段。

package api

// --- MV ---

type UpdateMusicMvReq struct {
	MusicId *int64 `json:"music_id,omitempty"`
	ThirdId *int64 `json:"third_id,omitempty"`
	MvUrl   string `json:"mv_url"`
}

type UpdateMusicMvResp struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	MusicId int64  `json:"music_id,omitempty"`
}

type GetMusicMvBatchReq struct {
	ThirdIds []int64 `json:"third_ids"`
}

type MusicMvItem struct {
	ThirdId int64  `json:"third_id"`
	MusicId int64  `json:"music_id"`
	MvUrl   string `json:"mv_url"`
}

type GetMusicMvBatchResp struct {
	Code    int32          `json:"code"`
	Message string         `json:"message"`
	Items   []*MusicMvItem `json:"items,omitempty"`
}
