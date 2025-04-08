package convert

import (
	"onij/infra/mysql"
	"onij/model/api"
	"onij/util"
)

func SingerNameToArtist(name string) *mysql.Artist {
	return &mysql.Artist{
		Id:         util.IdGen.Generate(),
		Name:       name,
		ArtistType: int32(api.ArtistType_AT_Singer),
	}
}
func NameToArtist(name string, typ api.ArtistType) *mysql.Artist {
	return &mysql.Artist{
		Id:         util.IdGen.Generate(),
		Name:       name,
		ArtistType: int32(typ),
	}
}
