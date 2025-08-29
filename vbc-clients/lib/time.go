package lib

import "time"

func TimeString(timeString string, format string) (time.Time, error) {
	return time.Parse(format, timeString)
}

func TimeParse(timeString string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, timeString)
	// 有 TZ后，不需要指定时区了
	//return time.ParseInLocation("2006-01-02T15:04:05.999Z", timeString, time.FixedZone("CST", 0))
}

func TimeEpoch1899() time.Time {
	return time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
}

func TimeEpoch1899Day(time time.Time) float64 {
	r := time.Sub(TimeEpoch1899())
	return r.Hours() / 24
}
