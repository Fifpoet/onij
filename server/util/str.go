package util

import (
	"code.chenji.com/pkg/boost/collection/collext"
	"onij/model/api"
	"path/filepath"
	"strconv"
	"strings"
)

const delimiter = ","

func DbToList(s string) []int {
	trimmed := strings.Trim(s, delimiter)
	ss := strings.Split(trimmed, delimiter)
	return collext.Pick(ss, func(str string) int {
		i, _ := strconv.Atoi(str)
		return i
	})
}

func GetFilenameWithoutExtension(localFilePath string) string {
	fullFilename := filepath.Base(localFilePath)
	// 去除后缀
	ext := filepath.Ext(fullFilename)
	return strings.TrimSuffix(fullFilename, ext)
}

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
