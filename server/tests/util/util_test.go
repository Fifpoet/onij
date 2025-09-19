package util

import (
	"onij/util"
	"onij/util/boost/crypto"
	"testing"
	"time"
)

func TestGetThisWeekDates(t *testing.T) {
	dates := util.GetThisWeekDates()
	t.Log(dates)
}

func TestGetFirstDayScore(t *testing.T) {
	t.Log(util.GetDayScore(time.Now()))
}

// 上传文件
// 同一个key, 如果已存在默认不覆盖
func TestUploadFile(t *testing.T) {
	//local := "C:\\KwDownload\\song\\Beyond-AMANI.lrc"
}

func TestDownload(t *testing.T) {
	url := util.DownloadFile("wy2.jpeg")
	t.Log(url)
}

func TestDelFile(t *testing.T) {
	err := util.DeleteFile("wy1.jpeg")
	if err != nil {
		t.Fatalf("delete file failed %v", err)
	}
	t.Log("delete file success")
}

func TestMd5(t *testing.T) {
	t.Log(crypto.Md5([]byte{1, 3, 2, 5, 53, 62, 1, 52, 12, 4, 6, 13}))
}

func TestGetFileName(t *testing.T) {
	t.Log(util.IdGen.Generate())
}

func TestGetList(t *testing.T) {
	files, err := util.GetFile("music/a")
	if err != nil {
		t.Fatalf("get file list failed %v", err)
		return
	}
	t.Log(files)
}

func TestGetSuffix(t *testing.T) {
	ext := util.GetFileSuffix("https://p1.music.126.net/VnZiScyynLG7atLIZ2YPkw==/18686200114669622.jpg")
	print(ext)
}
