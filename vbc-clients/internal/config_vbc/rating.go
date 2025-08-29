package config_vbc

type TypeFeeScheduleConfig map[int]int

var FeeScheduleConfig = TypeFeeScheduleConfig{
	0:   0,
	10:  171,
	20:  339,
	30:  524,
	40:  755,
	50:  1075,
	60:  1361,
	70:  1716,
	80:  1995,
	90:  2240,
	100: 3737,
}

func (c TypeFeeScheduleConfig) Pay(rating int) int {
	if _, ok := c[rating]; ok {
		return c[rating]
	}
	return 0
}
