package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"onij/model/api"
	"path/filepath"
	"strings"
)

func CalcFileHash(file []byte) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, bytes.NewReader(file)); err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)
	hashHex := hex.EncodeToString(hashBytes)

	return hashHex, nil
}

func GetFileSuffix(fileName string) string {
	ext := filepath.Ext(fileName)
	if ext != "" {
		return strings.TrimPrefix(ext, ".")
	}
	return ""
}

func GetFileType(fileName string) api.FileType {
	ext := GetFileSuffix(fileName)
	switch ext {
	case "pdf":
		return api.FileType_FT_Pdf
	case "doc", "docx":
		return api.FileType_FT_Word
	case "xls", "xlsx":
		return api.FileType_FT_Excel
	case "ppt", "pptx":
		return api.FileType_FT_PowerPoint
	case "jpg", "jpeg", "png", "gif", "bmp", "webp":
		return api.FileType_FT_Image
	case "mp4", "avi", "mov", "wmv", "flv":
		return api.FileType_FT_Video
	case "mp3", "wav", "flac", "aac":
		return api.FileType_FT_Audio
	case "txt", "md", "log", "json", "xml", "html", "css", "js":
		return api.FileType_FT_Text
	case "zip", "rar", "7z", "tar", "gz":
		return api.FileType_FT_Archive
	case "exe", "msi", "dmg":
		return api.FileType_FT_Executable
	case "go", "java", "py", "cpp", "c", "php":
		return api.FileType_FT_Code
	case "csv":
		return api.FileType_FT_CSV
	default:
		return api.FileType_FT_Unknown
	}
}
