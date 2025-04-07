package util

import (
	"code.chenji.com/pkg/boost/collection/collext"
	"path/filepath"
	"strconv"
	"strings"
)

const delimiter = ","

func ListToDb(l []int) string {
	s := collext.Pick(l, func(i int) string { return strconv.Itoa(i) })
	return delimiter + strings.Join(s, delimiter) + delimiter
}

func DbToList(s string) []int {
	trimmed := strings.Trim(s, delimiter)
	ss := strings.Split(trimmed, delimiter)
	return collext.Pick(ss, func(str string) int {
		i, _ := strconv.Atoi(str)
		return i
	})
}

// GetFilenameWithoutExtension 从路径提取无后缀的文件名
func GetFilenameWithoutExtension(localFilePath string) string {
	// 获取带后缀的文件名
	fullFilename := filepath.Base(localFilePath)
	// 去除后缀
	ext := filepath.Ext(fullFilename)
	return strings.TrimSuffix(fullFilename, ext)
}
