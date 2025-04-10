package util

import (
	"code.chenji.com/pkg/boost/collection/collext"
	"onij/model/api"
	"path/filepath"
	"strconv"
	"strings"
)

const delimiter = ","

func StrList2Int64(s string) []int64 {
	trimmed := strings.Trim(s, delimiter)
	ss := strings.Split(trimmed, delimiter)
	return collext.Pick(ss, func(str string) int64 {
		i, _ := strconv.ParseInt(str, 10, 64)
		return i
	})
}
func Int64List2Str(s []int64) string {
	return strings.Join(collext.Pick(s, func(i int64) string {
		return strconv.FormatInt(i, 10)
	}), delimiter)
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
