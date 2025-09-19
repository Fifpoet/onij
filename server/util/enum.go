package util

import "onij/model/api"

var (
	BaseCodeOK    = int32(200)
	BaseCodeError = int32(500)

	BaseMsgOK = "ok"
)

func NameToAlbumType(name string) api.AlbumType {
	switch name {
	case "专辑":
		return api.AlbumType_ALT_Studio
	case "a":
		return api.AlbumType_ALT_Live
	}
	return api.AlbumType_ALT_Studio
}
