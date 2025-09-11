package getter

import "onij/model"

func AlbumMusicMusicId(a *model.AlbumMusic) int64 {
	return a.MusicId
}
func AlbumMusicAlbumId(a *model.AlbumMusic) int64 {
	return a.AlbumId
}
