package util

import (
	"onij/util"
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
	local := "/Users/asen/Documents/swagger在线文档.png"
	str, err := util.UploadFile(local)
	if err != nil {
		t.Fatalf("upload file failed %v", err)
	}
	t.Log(str)
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

func TestHashFromLocalAndForm(t *testing.T) {
	hash, x, y, err := util.GetLocalFileHash("/Users/asen/Documents/后端API启动日志.png")
	if err != nil {
		return
	}
	t.Log(hash, x, y)
}
