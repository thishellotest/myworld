package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
	"vbc/lib"
)

// CalDelayDayTime 当执行时间在 08:00 PT- 18:00 PT之间时，立即执行。 不在这个区间时改为第二天 08：00 PT执行？
func CalDelayDayTime(executeTime time.Time, timeLocation time.Location) (dest time.Time) {

	dayHour := 8
	DayTimeAt := fmt.Sprintf("0%d:00", dayHour)
	executeTime = executeTime.In(&timeLocation)
	hour := executeTime.Hour()
	minute := executeTime.Minute()
	if hour >= 8 && hour <= 18 {
		return executeTime
	} else if hour == 19 && minute == 0 {
		return executeTime
	}
	if hour > dayHour {
		dest = executeTime.AddDate(0, 0, 1)
	} else {
		dest = executeTime
	}
	date := dest.In(&timeLocation).Format("2006-01-02")
	destDate := fmt.Sprintf("%s %s:00", date, DayTimeAt)
	lib.DPrintln(destDate)
	dest, _ = time.ParseInLocation(time.DateTime, destDate, &timeLocation)
	lib.DPrintln(dest.Format(time.RFC3339), timeLocation.String())
	return
}

// CalIntervalDayTime 计算间隔时间 DayTimeAt：16:00  (下午4点)
func CalIntervalDayTime(currentTime time.Time, intervalDay int, DayTimeAt string, timeLocation time.Location) (dest time.Time, err error) {
	//var timeLocation *time.Location
	//timeLocation = lib.GetVBCDefaultLocation()
	if intervalDay <= 0 {
		return dest, errors.New("intervalDay is wrong")
	}
	if DayTimeAt == "" {
		return dest, errors.New("DayTimeAt is wrong")
	}

	dest = currentTime.AddDate(0, 0, intervalDay)
	date := dest.In(&timeLocation).Format("2006-01-02")
	destDate := fmt.Sprintf("%s %s:00", date, DayTimeAt)
	dest, err = time.ParseInLocation(time.DateTime, destDate, &timeLocation)
	return
}
