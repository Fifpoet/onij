package conv

import "time"

func StringToDateTime(datetime string) time.Time {
	return stringToDateTime(datetime, nil)
}

func StringToLocalDateTime(datetime string) time.Time {
	loc, _ := time.LoadLocation("Local")
	return stringToDateTime(datetime, loc)
}

func StringToDateTimeUnixMilli(datetime string) int64 {
	return StringToDateTime(datetime).UnixMilli()
}

func StringToLocalDateTimeUnixMilli(datetime string) int64 {
	return StringToLocalDateTime(datetime).UnixMilli()
}

func stringToDateTime(datetime string, loc *time.Location) time.Time {
	var t time.Time
	for _, format := range []string{time.DateTime, time.RFC3339} {
		if v, err := time.ParseInLocation(format, datetime, loc); err == nil {
			t = v
			break
		}
	}
	return t
}
