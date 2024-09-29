package util

import (
	"onij/boost/collection/collext"
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
