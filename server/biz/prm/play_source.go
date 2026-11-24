package prm

import "onij/util"

type PlaySourceDTO struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Kind    string `json:"kind"`
	Origin  string `json:"origin"`
	Enabled bool   `json:"enabled"`
	Builtin bool   `json:"builtin"`
}

type ListPlaySourceParam struct{}

type SavePlaySourceParam struct {
	Sources []*PlaySourceDTO
}

type PlaySourceListResult struct {
	Sources []*PlaySourceDTO
}

func (r *PlaySourceListResult) Resp() map[string]any {
	sources := r.Sources
	if sources == nil {
		sources = []*PlaySourceDTO{}
	}
	return map[string]any{
		"code":    util.BaseCodeOK,
		"message": util.BaseMsgOK,
		"sources": sources,
	}
}

type GdstudioSignDTO struct {
	Version   string `json:"version"`
	Crc32     string `json:"crc32"`
	UpdatedAt int64  `json:"updatedAt"`
}

type GetGdstudioSignParam struct{}

type SaveGdstudioSignParam struct {
	Version string
	Crc32   string
}

type GdstudioSignResult struct {
	Sign *GdstudioSignDTO
}

func (r *GdstudioSignResult) Resp() map[string]any {
	sign := r.Sign
	if sign == nil {
		sign = &GdstudioSignDTO{}
	}
	return map[string]any{
		"code":      util.BaseCodeOK,
		"message":   util.BaseMsgOK,
		"version":   sign.Version,
		"crc32":     sign.Crc32,
		"updatedAt": sign.UpdatedAt,
	}
}
