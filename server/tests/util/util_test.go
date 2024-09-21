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

func TestUploadFile(t *testing.T) {
	local := "/Users/asen/Documents/wy1.jpeg"
	err := util.UploadFile(local, "wy1.jpeg")
	if err != nil {
		t.Fatalf("upload file failed %v", err)
	}
	t.Log("upload file success")
}
