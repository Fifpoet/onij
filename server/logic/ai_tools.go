package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"onij/biz/prm"
	"onij/model"
	"onij/util"
	"os"
	"strconv"
	"strings"
)

const (
	aiToolCallsMarker = "__tool_calls__"
	aiDefaultMaxSteps = 6
)

func aiAgentMaxSteps() int {
	v := strings.TrimSpace(os.Getenv("AI_AGENT_MAX_STEPS"))
	if v == "" {
		return aiDefaultMaxSteps
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return aiDefaultMaxSteps
	}
	if n > 20 {
		return 20
	}
	return n
}

func aiToolDefs() []util.DeepseekToolDef {
	return []util.DeepseekToolDef{
		fnTool("search_music", "搜索网易云单曲，返回精简曲目列表", map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{"type": "string", "description": "搜索关键词，优先「歌手 歌名」"},
				"limit": map[string]any{"type": "integer", "description": "条数，默认 8，最大 18"},
			},
			"required": []string{"query"},
		}),
		fnTool("get_song_detail", "按网易云 song_id 拉歌曲详情", map[string]any{
			"type": "object",
			"properties": map[string]any{
				"song_ids": map[string]any{
					"type":  "array",
					"items": map[string]any{"type": "integer"},
				},
			},
			"required": []string{"song_ids"},
		}),
		fnTool("collection_list", "列出站内音乐合集，可按关键词过滤", map[string]any{
			"type": "object",
			"properties": map[string]any{
				"keyword": map[string]any{"type": "string"},
			},
		}),
		fnTool("collection_create", "新建音乐合集", map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{"type": "string"},
			},
			"required": []string{"name"},
		}),
		fnTool("collection_detail", "查看合集详情与曲目 id 列表", map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id": map[string]any{"type": "integer", "description": "合集 id"},
			},
			"required": []string{"id"},
		}),
		fnTool("collection_add_songs", "向合集添加歌曲", map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id": map[string]any{"type": "integer"},
				"songs": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"song_id":      map[string]any{"type": "integer"},
							"song_name":    map[string]any{"type": "string"},
							"artist_names": map[string]any{"type": "string"},
						},
						"required": []string{"song_id"},
					},
				},
			},
			"required": []string{"id", "songs"},
		}),
		fnTool("collection_remove_song", "从合集移除一首歌", map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":      map[string]any{"type": "integer"},
				"song_id": map[string]any{"type": "integer"},
			},
			"required": []string{"id", "song_id"},
		}),
		fnTool("queue_add", "【前端】把歌曲加入播放队列；需 song_id", map[string]any{
			"type": "object",
			"properties": map[string]any{
				"song_id": map[string]any{"type": "integer"},
			},
			"required": []string{"song_id"},
		}),
		fnTool("queue_play_now", "【前端】立即播放指定歌曲；需 song_id", map[string]any{
			"type": "object",
			"properties": map[string]any{
				"song_id": map[string]any{"type": "integer"},
			},
			"required": []string{"song_id"},
		}),
		fnTool("player_toggle", "【前端】播放/暂停切换", map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		}),
		fnTool("player_next", "【前端】下一首", map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		}),
		fnTool("open_page", "【前端】打开站内页面", map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path":  map[string]any{"type": "string", "description": "如 /search /queue /collection /cloud"},
				"query": map[string]any{"type": "object", "description": "可选 query，如 {\"q\":\"林俊杰\",\"id\":\"123\"}"},
			},
			"required": []string{"path"},
		}),
	}
}

func fnTool(name, desc string, params map[string]any) util.DeepseekToolDef {
	return util.DeepseekToolDef{
		Type: "function",
		Function: map[string]any{
			"name":        name,
			"description": desc,
			"parameters":   params,
		},
	}
}

func aiSystemPrompt() string {
	return "你是 onij 站内 AI 助手，可通过 tools 操作音乐搜索、合集与播放器。" +
		"需要真实数据时必须调用工具，不要编造 song_id / 合集 id。" +
		"前端工具（queue_* / player_* / open_page）调用后会由浏览器执行，你可假设已派发成功并继续用简洁中文回复用户。" +
		"回答简洁，使用简体中文。"
}

func aiVoiceSystemPrompt() string {
	return aiSystemPrompt() +
		"当前是语音助手会话：口令短、可能有识别错字，可轻度纠正明显同音（如凌俊杰→林俊杰）。" +
		"回复一两句口语即可。不要让用户回复数字序号来选择。"
}

type aiToolExecResult struct {
	OK           bool
	Payload      any
	ClientAction map[string]any // 非空则需前端执行
	Error        string
}

func (l *aiLogic) executeAiTool(ctx context.Context, name string, argsJSON string) aiToolExecResult {
	args := map[string]any{}
	if strings.TrimSpace(argsJSON) != "" {
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return aiToolExecResult{OK: false, Error: "validation", Payload: map[string]any{"detail": "invalid json arguments"}}
		}
	}

	switch name {
	case "search_music":
		query := strings.TrimSpace(asString(args["query"]))
		if query == "" {
			return failValidation("query is required")
		}
		limit := asInt(args["limit"])
		if limit <= 0 {
			limit = 8
		}
		songs, err := util.NeteaseSearchSongs(query, limit)
		if err != nil {
			return failBiz(err.Error())
		}
		return aiToolExecResult{OK: true, Payload: map[string]any{"query": query, "count": len(songs), "songs": songs}}

	case "get_song_detail":
		ids := asInt64Slice(args["song_ids"])
		if len(ids) == 0 {
			return failValidation("song_ids is required")
		}
		songs, err := util.NeteaseSongDetails(ids)
		if err != nil {
			return failBiz(err.Error())
		}
		return aiToolExecResult{OK: true, Payload: map[string]any{"songs": songs}}

	case "collection_list":
		keyword := strings.TrimSpace(asString(args["keyword"]))
		res, err := l.CollectionLogicProxyList(ctx, keyword)
		if err != nil {
			return failBiz(err.Error())
		}
		return aiToolExecResult{OK: true, Payload: res}

	case "collection_create":
		nameArg := strings.TrimSpace(asString(args["name"]))
		if nameArg == "" {
			return failValidation("name is required")
		}
		res, err := l.CollectionLogicProxyCreate(ctx, nameArg)
		if err != nil {
			return failBiz(err.Error())
		}
		return aiToolExecResult{OK: true, Payload: res}

	case "collection_detail":
		id := asInt64(args["id"])
		if id <= 0 {
			return failValidation("id is required")
		}
		res, err := l.CollectionLogicProxyDetail(ctx, id)
		if err != nil {
			return failBiz(err.Error())
		}
		return aiToolExecResult{OK: true, Payload: res}

	case "collection_add_songs":
		id := asInt64(args["id"])
		if id <= 0 {
			return failValidation("id is required")
		}
		songsRaw, _ := args["songs"].([]any)
		if len(songsRaw) == 0 {
			return failValidation("songs is required")
		}
		metas := make([]prm.CollectionSongMeta, 0, len(songsRaw))
		for _, raw := range songsRaw {
			m, _ := raw.(map[string]any)
			if m == nil {
				continue
			}
			sid := asInt64(m["song_id"])
			if sid <= 0 {
				continue
			}
			metas = append(metas, prm.CollectionSongMeta{
				SongId:      sid,
				SongName:    strings.TrimSpace(asString(m["song_name"])),
				ArtistNames: strings.TrimSpace(asString(m["artist_names"])),
			})
		}
		if len(metas) == 0 {
			return failValidation("songs invalid")
		}
		res, err := l.CollectionLogicProxyAddSongs(ctx, id, metas)
		if err != nil {
			return failBiz(err.Error())
		}
		return aiToolExecResult{OK: true, Payload: res}

	case "collection_remove_song":
		id := asInt64(args["id"])
		songID := asInt64(args["song_id"])
		if id <= 0 || songID <= 0 {
			return failValidation("id and song_id required")
		}
		res, err := l.CollectionLogicProxyRemoveSong(ctx, id, songID)
		if err != nil {
			return failBiz(err.Error())
		}
		return aiToolExecResult{OK: true, Payload: res}

	case "queue_add", "queue_play_now":
		songID := asInt64(args["song_id"])
		if songID <= 0 {
			return failValidation("song_id is required")
		}
		action := map[string]any{"name": name, "args": map[string]any{"song_id": songID}}
		return aiToolExecResult{
			OK:           true,
			ClientAction: action,
			Payload:      map[string]any{"ok": true, "side": "frontend", "action": name, "song_id": songID},
		}

	case "player_toggle", "player_next":
		action := map[string]any{"name": name, "args": map[string]any{}}
		return aiToolExecResult{
			OK:           true,
			ClientAction: action,
			Payload:      map[string]any{"ok": true, "side": "frontend", "action": name},
		}

	case "open_page":
		path := strings.TrimSpace(asString(args["path"]))
		if path == "" || !strings.HasPrefix(path, "/") {
			return failValidation("path must start with /")
		}
		allowed := map[string]bool{
			"/search": true, "/queue": true, "/collection": true, "/cloud": true,
			"/monitor": true, "/ai": true, "/tran": true, "/prac": true, "/ktv": true, "/album": true, "/artist": true,
		}
		base := path
		if i := strings.IndexByte(path, '?'); i >= 0 {
			base = path[:i]
		}
		if !allowed[base] {
			return failValidation("path not allowed: " + base)
		}
		query, _ := args["query"].(map[string]any)
		actionArgs := map[string]any{"path": path}
		if query != nil {
			actionArgs["query"] = query
		}
		action := map[string]any{"name": name, "args": actionArgs}
		return aiToolExecResult{
			OK:           true,
			ClientAction: action,
			Payload:      map[string]any{"ok": true, "side": "frontend", "action": name, "path": path, "query": query},
		}

	default:
		return aiToolExecResult{OK: false, Error: "unknown_tool", Payload: map[string]any{"name": name}}
	}
}

func (l *aiLogic) CollectionLogicProxyList(ctx context.Context, keyword string) (any, error) {
	cl := NewCollectionLogic(l.AllInfra)
	res, err := cl.List(ctx, &prm.ListCollectionParam{Keyword: keyword})
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}

func (l *aiLogic) CollectionLogicProxyCreate(ctx context.Context, name string) (any, error) {
	cl := NewCollectionLogic(l.AllInfra)
	res, err := cl.Create(ctx, &prm.CreateCollectionParam{Name: name})
	if err != nil {
		return nil, err
	}
	return res.Item, nil
}

func (l *aiLogic) CollectionLogicProxyDetail(ctx context.Context, id int64) (any, error) {
	cl := NewCollectionLogic(l.AllInfra)
	res, err := cl.Detail(ctx, &prm.CollectionIdParam{Id: id})
	if err != nil {
		return nil, err
	}
	return res.Item, nil
}

func (l *aiLogic) CollectionLogicProxyAddSongs(ctx context.Context, id int64, songs []prm.CollectionSongMeta) (any, error) {
	cl := NewCollectionLogic(l.AllInfra)
	res, err := cl.AddSongs(ctx, &prm.AddCollectionSongsParam{Id: id, Songs: songs})
	if err != nil {
		return nil, err
	}
	return res.Item, nil
}

func (l *aiLogic) CollectionLogicProxyRemoveSong(ctx context.Context, id, songID int64) (any, error) {
	cl := NewCollectionLogic(l.AllInfra)
	res, err := cl.RemoveSong(ctx, &prm.RemoveCollectionSongParam{Id: id, SongId: songID})
	if err != nil {
		return nil, err
	}
	return res.Item, nil
}

func failValidation(detail string) aiToolExecResult {
	return aiToolExecResult{OK: false, Error: "validation", Payload: map[string]any{"detail": detail}}
}

func failBiz(detail string) aiToolExecResult {
	return aiToolExecResult{OK: false, Error: "business", Payload: map[string]any{"detail": detail}}
}

func toolResultJSON(r aiToolExecResult) string {
	body := map[string]any{"ok": r.OK}
	if r.Error != "" {
		body["error"] = r.Error
	}
	if r.Payload != nil {
		body["data"] = r.Payload
	}
	raw, _ := json.Marshal(body)
	return string(raw)
}

func asString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case fmt.Stringer:
		return t.String()
	case float64:
		return strconv.FormatInt(int64(t), 10)
	case json.Number:
		return t.String()
	default:
		return ""
	}
}

func asInt(v any) int {
	return int(asInt64(v))
}

func asInt64(v any) int64 {
	switch t := v.(type) {
	case float64:
		return int64(t)
	case int64:
		return t
	case int:
		return int64(t)
	case json.Number:
		n, _ := t.Int64()
		return n
	case string:
		n, _ := strconv.ParseInt(strings.TrimSpace(t), 10, 64)
		return n
	default:
		return 0
	}
}

func asInt64Slice(v any) []int64 {
	arr, ok := v.([]any)
	if !ok {
		if one := asInt64(v); one > 0 {
			return []int64{one}
		}
		return nil
	}
	out := make([]int64, 0, len(arr))
	for _, item := range arr {
		if n := asInt64(item); n > 0 {
			out = append(out, n)
		}
	}
	return out
}

func encodeAssistantToolCallsContent(text string, calls []util.DeepseekToolCall) string {
	raw, _ := json.Marshal(map[string]any{
		"content":    text,
		"tool_calls": calls,
	})
	return string(raw)
}

func decodeAssistantToolCallsContent(content string) (text string, calls []util.DeepseekToolCall, ok bool) {
	var envelope struct {
		Content   string                `json:"content"`
		ToolCalls []util.DeepseekToolCall `json:"tool_calls"`
	}
	if err := json.Unmarshal([]byte(content), &envelope); err != nil {
		return "", nil, false
	}
	if len(envelope.ToolCalls) == 0 {
		return "", nil, false
	}
	return envelope.Content, envelope.ToolCalls, true
}

func trimToLastUserTurns(rows []*model.AiMessage, maxUserTurns int) []*model.AiMessage {
	if maxUserTurns <= 0 || len(rows) == 0 {
		return rows
	}
	userIdx := make([]int, 0, 8)
	for i, r := range rows {
		if r != nil && r.Role == model.AiRoleUser {
			userIdx = append(userIdx, i)
		}
	}
	if len(userIdx) <= maxUserTurns {
		return rows
	}
	return rows[userIdx[len(userIdx)-maxUserTurns]:]
}

func buildDeepseekHistory(rows []*model.AiMessage, sess *model.AiSession) []util.DeepseekChatMessage {
	prompt := aiSystemPrompt()
	if sess != nil && sess.Kind == model.AiSessionKindVoice {
		prompt = aiVoiceSystemPrompt()
		rows = trimToLastUserTurns(rows, aiVoiceHistoryUserTurns)
	}
	out := []util.DeepseekChatMessage{{
		Role:    "system",
		Content: prompt,
	}}
	start := 0
	if len(rows) > aiHistoryMaxMessages {
		start = len(rows) - aiHistoryMaxMessages
	}
	for _, r := range rows[start:] {
		if r == nil {
			continue
		}
		switch r.Role {
		case model.AiRoleUser, model.AiRoleSystem:
			out = append(out, util.DeepseekChatMessage{Role: r.Role, Content: r.Content})
		case model.AiRoleAssistant:
			if r.ToolName == aiToolCallsMarker {
				text, calls, ok := decodeAssistantToolCallsContent(r.Content)
				if ok {
					out = append(out, util.DeepseekChatMessage{
						Role:      model.AiRoleAssistant,
						Content:   text,
						ToolCalls: calls,
					})
					continue
				}
			}
			out = append(out, util.DeepseekChatMessage{Role: model.AiRoleAssistant, Content: r.Content})
		case model.AiRoleTool:
			out = append(out, util.DeepseekChatMessage{
				Role:       model.AiRoleTool,
				Content:    r.Content,
				ToolCallID: r.ToolCallId,
				Name:       r.ToolName,
			})
		}
	}
	return out
}