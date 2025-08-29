package lib

import (
	"fmt"
	"testing"
	"time"
)

func Test_TimeEpoch1899(t *testing.T) {
	now := time.Now()
	now = time.Date(1899, 12, 30, 23, 2, 0, 0, time.UTC)
	a := TimeEpoch1899()

	r := now.Sub(a)
	fmt.Println(r.Hours() / 24)
}

func Test_a(t *testing.T) {

	str := `2024-03-10T18:36:14+08:00`
	c, _ := TimeParse(str)
	a := c.Format(time.RFC3339Nano)
	DPrintln(a)
}
