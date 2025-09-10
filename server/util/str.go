package util

import (
	"onij/model/api"
	"path/filepath"
)

func GetFileExtension(file string) string {
	fullFilename := filepath.Base(file)
	// 去除后缀
	ext := filepath.Ext(fullFilename)
	return ext
}

func GetFileType(file string) api.FileType {
	ext := GetFileExtension(file)
	if len(ext) <= 1 {
		return api.FileType_FT_Unknown
	}
	return fileTypeMap[ext[1:]]
}
