package getter

import (
	"onij/infra/mysql"
)

func AlbumCoverFileId(a *mysql.Album) int64 {return a.CoverFileId}
func AlbumMusicId(a *mysql.Album) int64 {return a.MusicId}
