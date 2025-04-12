package util

import "onij/model/api"

var fileTypeMap map[string]api.FileType = map[string]api.FileType{
	"mp3": api.FileType_FT_Mp3,
	"lyc": api.FileType_FT_Lyrics,
	"png": api.FileType_FT_Png,
}

var (
	BaseCodeOK    = int32(200)
	BaseCodeError = int32(500)

	BaseMsgOK = "ok"
)
