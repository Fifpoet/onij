package getter

import (
	"onij/model"
	"onij/util/boost/tool"
)

func MusicArtistIds(m *model.Music) []int64 {
	return *tool.LoadJson[[]int64](m.ArtistIds, true)
}
func MusicId(m *model.Music) int64 {
	return m.Id
}
func MusicThirdId(m *model.Music) int64 {
	return m.ThirdId
}
