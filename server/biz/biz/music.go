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

	Albums    []*AlbumPrime
	Artists   []*ArtistPrime
	Composers []*ArtistPrime
	Writers   []*ArtistPrime
	AudioFile *FilePrime
	LyricFile *FilePrime
	Tags      []*TagPrime
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
	return &api.Music{
		Id:            p.Id,
		Name:          exp.ValueOrZero(p.Name),
		Artists:       collext.Pick(collext.Combine(p.Artists, p.Composers, p.Writers), Artist),
		IssueTime:     exp.ValueOrZero(p.IssueTime),
		MvUrl:         exp.ValueOrZero(p.MvUrl),
		AudioFileUrl:  p.AudioFile.Url(),
		LyricsFileUrl: p.LyricFile.Url(),
		Album:         collext.Pick(p.Albums, Album),
		Tags:          collext.Pick(p.Tags, Tag),
	}
}
