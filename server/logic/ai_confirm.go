package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"onij/biz/prm"
	"onij/model"
	"onij/util"
	"strings"
	"time"
)

const aiConfirmTTL = 30 * time.Minute

func toolRequiresConfirm(name string) bool {
	switch name {
	case "collection_add_songs", "collection_remove_song":
		return true
	default:
		return false
	}
}

func (l *aiLogic) buildConfirmUI(ctx context.Context, toolName string, args map[string]any) (prm.AiConfirmUI, error) {
	switch toolName {
	case "collection_add_songs":
		collID := asInt64(args["id"])
		collName := fmt.Sprintf("合集#%d", collID)
		if collID > 0 {
			if item, err := l.CollectionLogicProxyDetail(ctx, collID); err == nil {
				if m, ok := item.(*prm.CollectionDTO); ok && m != nil && m.Name != "" {
					collName = m.Name
				}
			}
		}
		songs := parseSongMetas(args["songs"])
		if len(songs) == 0 {
			return prm.AiConfirmUI{}, fmt.Errorf("songs is empty")
		}
		if len(songs) == 1 {
			label := songs[0].SongName
			if label == "" {
				label = fmt.Sprintf("song_id=%d", songs[0].SongId)
			}
			if songs[0].ArtistNames != "" {
				label = label + " · " + songs[0].ArtistNames
			}
			return prm.AiConfirmUI{
				Type:     "yes_no",
				Title:    fmt.Sprintf("确认将「%s」加入合集「%s」？", label, collName),
				YesLabel: "加入",
				NoLabel:  "取消",
			}, nil
		}
		opts := make([]prm.AiConfirmOption, 0, len(songs))
		for _, s := range songs {
			label := s.SongName
			if label == "" {
				label = fmt.Sprintf("song_id=%d", s.SongId)
			}
			if s.ArtistNames != "" {
				label = label + " · " + s.ArtistNames
			}
			opts = append(opts, prm.AiConfirmOption{
				Id:    fmt.Sprintf("%d", s.SongId),
				Label: label,
				Meta:  s,
			})
		}
		return prm.AiConfirmUI{
			Type:    "multi_select",
			Title:   fmt.Sprintf("选择要加入合集「%s」的歌曲", collName),
			Options: opts,
			Min:     1,
			Max:     len(opts),
		}, nil

	case "collection_remove_song":
		collID := asInt64(args["id"])
		songID := asInt64(args["song_id"])
		collName := fmt.Sprintf("合集#%d", collID)
		if collID > 0 {
			if item, err := l.CollectionLogicProxyDetail(ctx, collID); err == nil {
				if m, ok := item.(*prm.CollectionDTO); ok && m != nil && m.Name != "" {
					collName = m.Name
				}
			}
		}
		return prm.AiConfirmUI{
			Type:     "yes_no",
			Title:    fmt.Sprintf("确认从合集「%s」移除歌曲 %d？", collName, songID),
			YesLabel: "移除",
			NoLabel:  "取消",
		}, nil

	default:
		return prm.AiConfirmUI{
			Type:     "yes_no",
			Title:    "确认执行该操作？",
			YesLabel: "确认",
			NoLabel:  "取消",
		}, nil
	}
}

func parseSongMetas(v any) []prm.CollectionSongMeta {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]prm.CollectionSongMeta, 0, len(arr))
	for _, raw := range arr {
		m, _ := raw.(map[string]any)
		if m == nil {
			continue
		}
		sid := asInt64(m["song_id"])
		if sid <= 0 {
			continue
		}
		out = append(out, prm.CollectionSongMeta{
			SongId:      sid,
			SongName:    strings.TrimSpace(asString(m["song_name"])),
			ArtistNames: strings.TrimSpace(asString(m["artist_names"])),
		})
	}
	return out
}

func (l *aiLogic) createConfirmPending(
	ctx context.Context,
	sess *model.AiSession,
	toolCallID, toolName, argsJSON string,
	args map[string]any,
) (*prm.AiConfirmDTO, error) {
	ui, err := l.buildConfirmUI(ctx, toolName, args)
	if err != nil {
		return nil, err
	}
	uiRaw, _ := json.Marshal(ui)
	now := time.Now()
	row := &model.AiConfirmPending{
		Id:         util.IdGen.Generate(),
		SessionId:  sess.Id,
		ToolCallId: toolCallID,
		ToolName:   toolName,
		DraftArgs:  argsJSON,
		UiJson:     string(uiRaw),
		Status:     model.AiConfirmStatusPending,
		ExpiresAt:  now.Add(aiConfirmTTL),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if _, err = l.AiConfirmPendingDal.Save(ctx, row); err != nil {
		return nil, err
	}
	var draft any
	_ = json.Unmarshal([]byte(argsJSON), &draft)
	return &prm.AiConfirmDTO{
		ConfirmId:  row.Id,
		SessionId:  sess.Id,
		ToolCallId: toolCallID,
		ToolName:   toolName,
		DraftArgs:  draft,
		UI:         ui,
		ExpiresAt:  row.ExpiresAt.UnixMilli(),
	}, nil
}

func applyConfirmDecision(toolName string, draftArgs map[string]any, ui prm.AiConfirmUI, d prm.AiConfirmDecision) (finalArgs map[string]any, cancelled bool, err error) {
	if d.Type == "cancel" || (d.Type == "yes_no" && d.Value != nil && !*d.Value) {
		return nil, true, nil
	}
	finalArgs = map[string]any{}
	for k, v := range draftArgs {
		finalArgs[k] = v
	}

	switch ui.Type {
	case "yes_no":
		if d.Type != "yes_no" || d.Value == nil || !*d.Value {
			return nil, true, nil
		}
		return finalArgs, false, nil

	case "single_select":
		opt := strings.TrimSpace(d.OptionId)
		if opt == "" {
			return nil, false, fmt.Errorf("option_id required")
		}
		if !optionAllowed(ui.Options, opt) {
			return nil, false, fmt.Errorf("option_id invalid")
		}
		if toolName == "collection_add_songs" {
			finalArgs["songs"] = filterSongsByIDs(draftArgs["songs"], []string{opt})
		}
		return finalArgs, false, nil

	case "multi_select":
		ids := d.OptionIds
		if len(ids) == 0 {
			return nil, false, fmt.Errorf("option_ids required")
		}
		for _, id := range ids {
			if !optionAllowed(ui.Options, id) {
				return nil, false, fmt.Errorf("option_id invalid: %s", id)
			}
		}
		if ui.Min > 0 && len(ids) < ui.Min {
			return nil, false, fmt.Errorf("select at least %d", ui.Min)
		}
		if ui.Max > 0 && len(ids) > ui.Max {
			return nil, false, fmt.Errorf("select at most %d", ui.Max)
		}
		if toolName == "collection_add_songs" {
			finalArgs["songs"] = filterSongsByIDs(draftArgs["songs"], ids)
		}
		return finalArgs, false, nil

	default:
		return nil, false, fmt.Errorf("unsupported ui type")
	}
}

func optionAllowed(opts []prm.AiConfirmOption, id string) bool {
	for _, o := range opts {
		if o.Id == id {
			return true
		}
	}
	return false
}

func filterSongsByIDs(songsRaw any, ids []string) []any {
	want := map[string]bool{}
	for _, id := range ids {
		want[id] = true
	}
	arr, _ := songsRaw.([]any)
	out := make([]any, 0, len(ids))
	for _, raw := range arr {
		m, _ := raw.(map[string]any)
		if m == nil {
			continue
		}
		sid := fmt.Sprintf("%d", asInt64(m["song_id"]))
		if want[sid] {
			out = append(out, m)
		}
	}
	return out
}

func (l *aiLogic) ConfirmStream(ctx context.Context, param *prm.AiConfirmParam, emit AiStreamEmitter) error {
	if emit == nil {
		return fmt.Errorf("emitter is nil")
	}
	if param == nil || param.ConfirmId <= 0 {
		return fmt.Errorf("confirm_id is empty")
	}
	key := util.DeepseekAPIKey()
	if key == "" {
		return fmt.Errorf("dskey not configured")
	}

	pending, err := l.AiConfirmPendingDal.GetById(ctx, param.ConfirmId)
	if err != nil {
		return err
	}
	if pending == nil {
		return fmt.Errorf("confirm not found")
	}
	if pending.Status != model.AiConfirmStatusPending {
		return fmt.Errorf("confirm already %s", pending.Status)
	}
	if time.Now().After(pending.ExpiresAt) {
		_ = l.AiConfirmPendingDal.UpdateStatus(ctx, pending.Id, model.AiConfirmStatusExpired)
		return fmt.Errorf("confirm expired")
	}

	sess, err := l.AiSessionDal.GetById(ctx, pending.SessionId)
	if err != nil {
		return err
	}
	if sess == nil {
		return fmt.Errorf("session not found")
	}
	_ = emit("session", toAiSessionDTO(sess))

	var ui prm.AiConfirmUI
	if err = json.Unmarshal([]byte(pending.UiJson), &ui); err != nil {
		return fmt.Errorf("invalid ui_json")
	}
	draftArgs := map[string]any{}
	if err = json.Unmarshal([]byte(pending.DraftArgs), &draftArgs); err != nil {
		return fmt.Errorf("invalid draft_args")
	}

	finalArgs, cancelled, err := applyConfirmDecision(pending.ToolName, draftArgs, ui, param.Decision)
	if err != nil {
		return err
	}

	var exec aiToolExecResult
	if cancelled {
		_ = l.AiConfirmPendingDal.UpdateStatus(ctx, pending.Id, model.AiConfirmStatusCancelled)
		exec = aiToolExecResult{
			OK:      false,
			Error:   "cancelled_by_user",
			Payload: map[string]any{"detail": "用户取消"},
		}
	} else {
		argsJSON, _ := json.Marshal(finalArgs)
		exec = l.executeAiTool(ctx, pending.ToolName, string(argsJSON))
		if exec.ClientAction != nil {
			_ = emit("client_action", exec.ClientAction)
		}
		st := model.AiConfirmStatusConfirmed
		if !exec.OK {
			st = model.AiConfirmStatusConfirmed // 仍算确认过，只是业务失败
		}
		_ = l.AiConfirmPendingDal.UpdateStatus(ctx, pending.Id, st)
	}

	resultContent := toolResultJSON(exec)
	toolMsg := &model.AiMessage{
		Id:         util.IdGen.Generate(),
		SessionId:  sess.Id,
		Role:       model.AiRoleTool,
		Content:    resultContent,
		ToolName:   pending.ToolName,
		ToolCallId: pending.ToolCallId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if _, err = l.AiMessageDal.Save(ctx, toolMsg); err != nil {
		return err
	}
	_ = emit("tool_result", map[string]any{
		"id":      pending.ToolCallId,
		"name":    pending.ToolName,
		"ok":      exec.OK,
		"content": resultContent,
	})

	outcome, err := l.runAgentLoop(ctx, key, sess, emit)
	if err != nil {
		return err
	}
	if outcome.Paused != nil {
		_ = emit("needs_confirm", outcome.Paused)
		_ = emit("paused", map[string]any{
			"session":    toAiSessionDTO(sess),
			"confirm_id": outcome.Paused.ConfirmId,
		})
		return nil
	}
	return emit("done", map[string]any{
		"session": toAiSessionDTO(sess),
		"message": toAiMessageDTO(outcome.Message),
		"reply":   outcome.Message.Content,
	})
}
