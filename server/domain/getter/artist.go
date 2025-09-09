package getter

import (
	"onij/infra/mysql"
	"onij/model/api"
)

func ArtistId(a *mysql.Artist) int64 {return a.Id}
func ArtistName(a *mysql.Artist) string {return a.Name}
func ArtistToProfile(a *mysql.Artist) *api.ArtistProfile {
	if a == nil {
		return nil
	}
	return &api.ArtistProfile{
		Id:         a.Id,
		Name:       a.Name,
		ArtistType: api.ArtistType(a.ArtistType),
	}
}