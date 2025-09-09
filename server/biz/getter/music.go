package getter

import (
	"onij/util"
)

func MusicArtistIds(m *model.Music) []int64 {
	return util.StrList2Int64(m.ArtistIds)
}
func MusicId(m *model.Music) int64 {
	return m.Id
}
