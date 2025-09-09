package getter

import (
	"onij/model/api"
)

func ArtistId(a *model.Artist) int64    { return a.Id }
func ArtistName(a *model.Artist) string { return a.Name }
func ArtistToProfile(a *model.Artist) *api.ArtistProfile {
	if a == nil {
		return nil
	}
	return &api.ArtistProfile{
		Id:         a.Id,
		Name:       a.Name,
		ArtistType: api.ArtistType(a.ArtistType),
	}
}
