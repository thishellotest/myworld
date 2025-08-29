package biz

import (
	"time"
	"vbc/configs"
)

func TimestampToDate(val int64) string {
	if val == 0 {
		return ""
	}
	return time.Unix(val, 0).In(configs.GetVBCDefaultLocation()).Format(time.DateOnly)
}

func TimeDateOnlyToTimestamp(date string) int64 {
	t, _ := time.ParseInLocation(time.DateOnly, date, configs.GetVBCDefaultLocation())
	return t.Unix()
}

func TimestampToString(timeLocation *time.Location, val int64) string {
	if val > 0 {
		a := time.Unix(val, 0).In(timeLocation)
		r := a.Format(configs.TimeFormatDateTime)
		return r
	}
	return ""
}

func TimestampToStringByUserFacade(userFacade *UserFacade, TimezonesUsecase *TimezonesUsecase, val int64) (string, error) {

	var timeLocation *time.Location
	var err error
	if userFacade != nil {
		timeLocation, err = userFacade.GetTimeLocation(TimezonesUsecase)
		if err != nil {
			return "", err
		}
	}
	r := TimestampToString(timeLocation, val)
	return r, nil
}
