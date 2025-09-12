package getter

import (
	"onij/model"
	"onij/util/boost/tool"
)

func AlbumMusicMusicId(a *model.AlbumMusic) int64 {
	return a.MusicId
}
func AlbumMusicAlbumId(a *model.AlbumMusic) int64 {
	return a.AlbumId
}
func AlbumMusicArtistIds(a *model.AlbumMusic) []int64 {
	return *tool.LoadJson[[]int64](a.ArtistIds, true)
}
