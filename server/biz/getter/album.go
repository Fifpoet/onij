package getter

func AlbumCoverFileId(a *model.Album) int64 { return a.CoverFileId }
func AlbumMusicId(a *model.Album) int64     { return a.MusicId }
