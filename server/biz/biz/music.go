package biz

import "onij/model/api"

type MusicPrime struct {
	Id           *int64
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
