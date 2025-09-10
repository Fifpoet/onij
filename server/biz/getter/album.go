package getter

import "onij/model"

func AlbumCoverFileId(a *model.Album) int64 { return a.CoverFileId }
