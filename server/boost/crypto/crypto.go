package crypto

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func Md5(data []byte) string {
	return strings.ToUpper(fmt.Sprintf("%x", md5.Sum(data)))
}
