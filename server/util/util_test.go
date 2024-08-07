package util

import (
	"testing"
	"time"
)

func TestGetThisWeekDates(t *testing.T) {
	dates := GetThisWeekDates()
	t.Log(dates)
}

func TestGetFirstDayScore(t *testing.T) {
	t.Log(GetDayScore(time.Now()))
}
