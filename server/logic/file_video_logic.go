package logic

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"onij/biz/biz"
	"onij/biz/prm"
	"onij/infra"
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/logs"
)

const (
	videoLongDurationMs = 20 * 60 * 1000
	videoLongSizeBytes  = 512 * 1024 * 1024
)

type FileVideoLogic interface {
	Play(ctx context.Context, param *prm.PlayFileVideoParam) (*prm.PlayFileVideoResult, error)
	ListMarkers(ctx context.Context, param *prm.ListFileVideoMarkerParam) (*prm.ListFileVideoMarkerResult, error)
	SaveMarkers(ctx context.Context, param *prm.SaveFileVideoMarkerParam) (*prm.SaveFileVideoMarkerResult, error)
	DeleteMarker(ctx context.Context, param *prm.DeleteFileVideoMarkerParam) (*prm.DeleteFileVideoMarkerResult, error)
}

type fileVideoLogic struct {
	*infra.AllInfra
}

func NewFileVideoLogic(i *infra.AllInfra) FileVideoLogic {
	return &fileVideoLogic{AllInfra: i}
}

func (l *fileVideoLogic) Play(ctx context.Context, param *prm.PlayFileVideoParam) (*prm.PlayFileVideoResult, error) {
	if param.FileId <= 0 {
		return nil, fmt.Errorf("invalid file_id")
	}
	fis, err := l.FileDal.GetByIds(ctx, param.FileId)
	if err != nil {
		return nil, err
	}
	if len(fis) == 0 {
		return nil, fmt.Errorf("file not found")
	}
	fi := fis[0]
	if fi.Format != int32(api.FileType_FT_Video) {
		return &prm.PlayFileVideoResult{Mode: prm.VideoPlayModeUnsupported, Hint: "不是视频文件"}, nil
	}
	if !isMp4Name(fi.Name) {
		return &prm.PlayFileVideoResult{Mode: prm.VideoPlayModeUnsupported, Hint: "仅支持 MP4"}, nil
	}
	if len(fi.StoreKey) == 0 {
		return nil, fmt.Errorf("empty store_key")
	}

	var durationMs int64
	info, avErr := util.Avinfo(fi.StoreKey)
	if avErr != nil {
		logs.Error("fileVideo Play avinfo: %v", avErr)
	} else {
		durationMs = info.DurationMs()
	}

	if shouldUseHLS(durationMs, fi.Size) {
		return l.playHLS(ctx, fi, durationMs)
	}
	return l.playMP4(fi, durationMs)
}

func (l *fileVideoLogic) playMP4(fi *model.File, durationMs int64) (*prm.PlayFileVideoResult, error) {
	ok, err := util.ObjectHeadHasMoovFirst(fi.StoreKey)
	if err != nil {
		logs.Error("fileVideo playMP4 head probe: %v", err)
		ok = true
	}
	if !ok {
		if _, pfErr := util.PfopFaststart(fi.StoreKey); pfErr != nil {
			logs.Error("fileVideo playMP4 pfop faststart: %v", pfErr)
		}
		return &prm.PlayFileVideoResult{
			Mode:       prm.VideoPlayModeProcessing,
			DurationMs: durationMs,
			Hint:       "正在处理以便拖动进度，请稍后重试",
		}, nil
	}
	return &prm.PlayFileVideoResult{
		Mode:       prm.VideoPlayModeMP4,
		URL:        util.SignKeyURL(fi.StoreKey),
		DurationMs: durationMs,
	}, nil
}

func (l *fileVideoLogic) playHLS(ctx context.Context, fi *model.File, durationMs int64) (*prm.PlayFileVideoResult, error) {
	m3u8Key := util.HLSKey(fi.StoreKey)
	if util.ObjectExists(m3u8Key) {
		return l.hlsReady(m3u8Key, durationMs)
	}

	pid := biz.ParseHlsPid(fi.Extra)
	if pid != "" {
		st, err := util.Prefop(pid)
		if err != nil {
			logs.Error("fileVideo playHLS prefop: %v", err)
		} else if util.PrefopOK(st.Code) {
			if util.ObjectExists(m3u8Key) {
				return l.hlsReady(m3u8Key, durationMs)
			}
		} else if util.PrefopBusy(st.Code) {
			return &prm.PlayFileVideoResult{
				Mode:       prm.VideoPlayModeProcessing,
				DurationMs: durationMs,
				Hint:       "切片处理中，请稍后重试",
			}, nil
		}
	}

	newPid, err := util.PfopHLS(fi.StoreKey)
	if err != nil {
		return nil, err
	}
	if newPid != "" && newPid != pid {
		fi.Extra = biz.SetHlsPid(fi.Extra, newPid)
		if _, saveErr := l.FileDal.Save(ctx, fi); saveErr != nil {
			logs.Error("fileVideo playHLS save hls_pid: %v", saveErr)
		}
	}
	return &prm.PlayFileVideoResult{
		Mode:       prm.VideoPlayModeProcessing,
		DurationMs: durationMs,
		Hint:       "已提交切片，请稍后重试",
	}, nil
}

func (l *fileVideoLogic) hlsReady(m3u8Key string, durationMs int64) (*prm.PlayFileVideoResult, error) {
	text, err := util.FetchM3U8Text(m3u8Key)
	if err != nil {
		return nil, err
	}
	return &prm.PlayFileVideoResult{
		Mode:       prm.VideoPlayModeHLS,
		Playlist:   util.RewriteM3U8(text, m3u8Key),
		DurationMs: durationMs,
	}, nil
}

func (l *fileVideoLogic) ListMarkers(ctx context.Context, param *prm.ListFileVideoMarkerParam) (*prm.ListFileVideoMarkerResult, error) {
	if param.FileId <= 0 {
		return nil, fmt.Errorf("invalid file_id")
	}
	rows, err := l.FileVideoMarkerDal.ListByFileId(ctx, param.FileId)
	if err != nil {
		return nil, err
	}
	items := make([]*api.FileVideoMarkerItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, prm.ToFileVideoMarkerDTO(row))
	}
	return &prm.ListFileVideoMarkerResult{Markers: items}, nil
}

func (l *fileVideoLogic) SaveMarkers(ctx context.Context, param *prm.SaveFileVideoMarkerParam) (*prm.SaveFileVideoMarkerResult, error) {
	if param.FileId <= 0 {
		return nil, fmt.Errorf("invalid file_id")
	}
	now := time.Now()
	rows := make([]*model.FileVideoMarker, 0, len(param.Markers))
	for _, m := range param.Markers {
		if m == nil {
			continue
		}
		if m.TimeMs < 0 {
			return nil, fmt.Errorf("invalid time_ms")
		}
		id := m.Id
		if id <= 0 {
			id = util.IdGen.Generate()
		}
		rows = append(rows, &model.FileVideoMarker{
			Id:        id,
			FileId:    param.FileId,
			TimeMs:    m.TimeMs,
			Label:     strings.TrimSpace(m.Label),
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	if len(rows) > 0 {
		if _, err := l.FileVideoMarkerDal.Save(ctx, rows...); err != nil {
			return nil, err
		}
	}
	listed, err := l.ListMarkers(ctx, &prm.ListFileVideoMarkerParam{FileId: param.FileId})
	if err != nil {
		return nil, err
	}
	return &prm.SaveFileVideoMarkerResult{Markers: listed.Markers}, nil
}

func (l *fileVideoLogic) DeleteMarker(ctx context.Context, param *prm.DeleteFileVideoMarkerParam) (*prm.DeleteFileVideoMarkerResult, error) {
	if param.Id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}
	if err := l.FileVideoMarkerDal.DeleteById(ctx, param.Id); err != nil {
		return nil, err
	}
	return &prm.DeleteFileVideoMarkerResult{}, nil
}

func shouldUseHLS(durationMs, size int64) bool {
	return durationMs > videoLongDurationMs || size > videoLongSizeBytes
}

func isMp4Name(name string) bool {
	return strings.EqualFold(path.Ext(name), ".mp4")
}
