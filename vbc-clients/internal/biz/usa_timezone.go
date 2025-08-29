package biz

import (
	"github.com/pkg/errors"
	"time"
	"vbc/lib"
)

/*

Time Zone Abbreviation & Name	Offset	Current Time
PT	Pacific Time	UTC -8:00 / -7:00	Tue, 03:41:12   America/Los_Angeles
MT	Mountain Time	UTC -7:00 / -6:00	Tue, 04:41:12   America/Denver
CT	Central Time	UTC -6:00 / -5:00	Tue, 05:41:12   America/Chicago
ET	Eastern Time	UTC -5:00 / -4:00	Tue,            America/New_York

// https://www.timeanddate.com/time/zone/usa
// https://simple.wikipedia.org/wiki/List_of_U.S._states_and_territories_by_time_zone
// https://www.time-zones-map.com/   这个时区很重要参考的网站
// https://state.1keydata.com/time-zone-state.php
*/

const (
	USA_TIMEZONE_PT = "PT"
	USA_TIMEZONE_MT = "MT"
	USA_TIMEZONE_CT = "CT"
	USA_TIMEZONE_ET = "ET"

	USA_TIMEZONE_AK = "Alaska Time Zone"
	USA_TIMEZONE_HI = "Hawaii-Aleutian Time Zone"
)

func GetLocationByUsaTimezone(usaTimezone string) (*time.Location, error) {
	if usaTimezone == USA_TIMEZONE_PT {
		return time.LoadLocation("America/Los_Angeles")
	} else if usaTimezone == USA_TIMEZONE_MT {
		return time.LoadLocation("America/Denver")
	} else if usaTimezone == USA_TIMEZONE_CT {
		return time.LoadLocation("America/Chicago")
	} else if usaTimezone == USA_TIMEZONE_ET {
		return time.LoadLocation("America/New_York")
	} else if usaTimezone == USA_TIMEZONE_AK {
		return time.LoadLocation("America/Anchorage")
	} else if usaTimezone == USA_TIMEZONE_HI {
		// "Etc/GMT+10"：
		// 这是一个符合 IANA 时区数据库 中的时区名称。
		// 特别注意，在 Etc/GMT 时区表示中，符号的使用与我们通常的认知是相反的：
		// "Etc/GMT+10" 实际上表示比 UTC 慢 10 小时，即UTC-10。
		return time.LoadLocation("Etc/GMT+10")
	} else {
		return nil, errors.New("Error timezone")
	}
}

func UsaTime() {
	now := time.Now()
	pt, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
	lib.DPrintln("PT:", now.In(pt).Format(time.DateTime))

	mt, err := time.LoadLocation("America/Denver")
	if err != nil {
		panic(err)
	}
	lib.DPrintln("MT:", now.In(mt).Format(time.DateTime))

	ct, err := time.LoadLocation("America/Chicago")
	if err != nil {
		panic(err)
	}
	lib.DPrintln("MT:", now.In(ct).Format(time.DateTime))

	et, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	lib.DPrintln("ET:", now.In(et).Format(time.DateTime))
}
