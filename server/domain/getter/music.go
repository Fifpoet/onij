package getter

import (
	"onij/infra/mysql"
	"onij/util"
)

func MusicArtistIds(m *mysql.Music) []int64 {
	return util.StrList2Int64(m.ArtistIds)
}
func MusicId(m *mysql.Music) int64 {
	return m.Id
}