package util

import (
	"time"
)

// GetThisWeekDates 获取本周的日期
func GetThisWeekDates() (res []time.Time) {
	now := time.Now()
	weekday := int(now.Weekday())

	lastSunday := now.AddDate(0, 0, -weekday)
	// 打印从上周日到本周六的日期
	for i := 0; i < 7; i++ {
		day := lastSunday.AddDate(0, 0, i)
		res = append(res, day)
	}
	return res
}

// GetDayScore 获取距离2024-08-01的天数（db存储）
func GetDayScore(day time.Time) int {
	startDate := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	duration := day.Sub(startDate)
	// 计算天数
	score := int(duration.Hours() / 24)
	return score
}
