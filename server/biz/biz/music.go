package biz

import (
	"onij/model"
	"onij/model/api"
	"onij/util/boost/collection/collext"
	"onij/util/boost/exp"
	"onij/util/boost/tool"
	"time"
)

type MusicPrime struct {
	Id           int64
	Name         *string
	FullName     *string
	ArtistIds    []int64
	ComposerIds  []int64
	WriterIds    []int64
	IssueTime    *int32
	PerformType  *api.PerformType
	TimeLength   *int32
	MvUrl        *string
	AudioFileId  *int64
	AudioQuality *api.AudioQuality
	LyricFileId  *int64
	LyricContent *string
	RootId       *int64
	Priority     *int32
	CreatedAt    *int64
	UpdatedAt    *int64

	Albums  []*AlbumPrime
	Artists []*ArtistPrime // 包含artist composer writer
	Files   []*FilePrime   // 包含audio lyric
	Tags    []*TagPrime
}

func (m *MusicPrime) Model() *model.Music {
	return &model.Music{
		Id:           m.Id,
		Name:         exp.ValueOrZero(m.Name),
		FullName:     exp.ValueOrZero(m.FullName),
		ArtistIds:    tool.ToJson(m.ArtistIds),
		ComposerIds:  tool.ToJson(m.ComposerIds),
		WriterIds:    tool.ToJson(m.WriterIds),
		IssueTime:    time.Unix(int64(exp.ValueOrZero(m.IssueTime)), 0),
		PerformType:  int32(exp.ValueOrZero(m.PerformType)),
		TimeLength:   exp.ValueOrZero(m.TimeLength),
		MvUrl:        exp.ValueOrZero(m.MvUrl),
		AudioFileId:  exp.ValueOrZero(m.AudioFileId),
		AudioQuality: int32(exp.ValueOrZero(m.AudioQuality)),
		LyricFileId:  exp.ValueOrZero(m.LyricFileId),
		LyricContent: exp.ValueOrZero(m.LyricContent),
		RootId:       exp.ValueOrZero(m.RootId),
		Priority:     exp.ValueOrZero(m.Priority),
	}
}

func Music(p *MusicPrime) *api.Music {
	artistMap := collext.Map(p.ArtistIds, func(a int64) int64 {
		return a
	})
	composerMap := collext.Map(p.ComposerIds, func(a int64) int64 {
		return a
	})
	writerMap := collext.Map(p.WriterIds, func(a int64) int64 {
		return a
	})
	return &api.Music{
		Id:   p.Id,
		Name: exp.ValueOrZero(p.Name),
		Artists: collext.Pick(collext.Select(p.Artists, func(a *ArtistPrime) (*ArtistPrime, bool) {
			return a, artistMap[a.Id] > 0
		}), Artist),
		Composers: collext.Pick(collext.Select(p.Artists, func(a *ArtistPrime) (*ArtistPrime, bool) {
			return a, composerMap[a.Id] > 0
		}), Artist),
		Writers: collext.Pick(collext.Select(p.Artists, func(a *ArtistPrime) (*ArtistPrime, bool) {
			return a, writerMap[a.Id] > 0
		}), Artist),
		IssueTime: exp.ValueOrZero(p.IssueTime),
		MvUrl:     exp.ValueOrZero(p.MvUrl),
		AudioFileUrl: collext.SelectOne(p.Files, func(a *FilePrime) (string, bool) {
			return a.Url(), a.Id == exp.ValueOrZero(p.AudioFileId)
		}),
		LyricsFileUrl: collext.SelectOne(p.Files, func(a *FilePrime) (string, bool) {
			return a.Url(), a.Id == exp.ValueOrZero(p.LyricFileId)
		}),
		Album: collext.Pick(p.Albums, Album),
		Tags:  collext.Pick(p.Tags, Tag),
	}
}

func MusicToBiz(m *model.Music) *MusicPrime {
	return &MusicPrime{
		Id:           m.Id,
		Name:         exp.Ptr(m.Name),
		FullName:     exp.Ptr(m.FullName),
		ArtistIds:    *tool.LoadJson[[]int64](m.ArtistIds, true),
		ComposerIds:  *tool.LoadJson[[]int64](m.ComposerIds, true),
		WriterIds:    *tool.LoadJson[[]int64](m.WriterIds, true),
		IssueTime:    exp.Ptr(int32(m.IssueTime.Unix())),
		PerformType:  exp.Ptr(api.PerformType(m.PerformType)),
		TimeLength:   exp.Ptr(m.TimeLength),
		MvUrl:        exp.Ptr(m.MvUrl),
		AudioFileId:  exp.Ptr(m.AudioFileId),
		AudioQuality: exp.Ptr(api.AudioQuality(m.AudioQuality)),
		LyricFileId:  exp.Ptr(m.LyricFileId),
		LyricContent: exp.Ptr(m.LyricContent),
		RootId:       exp.Ptr(m.RootId),
		Priority:     exp.Ptr(m.Priority),
		CreatedAt:    exp.Ptr(m.CreatedAt.Unix()),
		UpdatedAt:    exp.Ptr(m.UpdatedAt.Unix()),
	}
}
