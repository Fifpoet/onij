package getter

import (
	"onij/model"
)

func ArtistId(a *model.Artist) int64    { return a.Id }
func ArtistName(a *model.Artist) string { return a.Name }
