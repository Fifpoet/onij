package getter

import (
	"onij/model"
	"onij/util/boost/tool"
)

func AlbumId(a *model.Album) int64 { return a.Id }
func AlbumCoverFileId(a *model.Album) int64 { return a.CoverFileId }
func AlbumArtistIds(a *model.Album) []int64 { return *tool.LoadJson[[]int64](a.ArtistIds, true) }