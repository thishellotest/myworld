package tests

import (
	"fmt"
	"testing"
	"time"
	"vbc/configs"
	"vbc/lib"
)

func Test_PollingTransIdByTime(t *testing.T) {
	org1 := "2006-01-01 00:00:00"
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04:05", org1, configs.LoadLocation)
	a, begin := configs.PollingTransIdByTime(currentTime)
	fmt.Println("org:", org1, a, begin.Format("2006-01-02 15:04:05"))

	org1 = "2006-01-02 15:01:00"
	currentTime, _ = time.Parse("2006-01-02 15:04:05", org1)
	a, begin = configs.PollingTransIdByTime(currentTime)
	fmt.Println("org:", org1, a, begin.Format("2006-01-02 15:04:05"))

	org1 = "2006-01-02 15:55:00"
	currentTime, _ = time.Parse("2006-01-02 15:04:05", org1)
	a, begin = configs.PollingTransIdByTime(currentTime)
	fmt.Println("org:", org1, a, begin.Format("2006-01-02 15:04:05"))

	org1 = "2006-01-02 15:56:59"
	currentTime, _ = time.Parse("2006-01-02 15:04:05", "2006-01-02 15:59:00")
	a, begin = configs.PollingTransIdByTime(currentTime)
	fmt.Println("org:", org1, a, begin.Format("2006-01-02 15:04:05"))
}

func Test_time(t *testing.T) {
	a := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(a)
}

func Test_a(t *testing.T) {
	url := "https://hooks.zapier.com/hooks/catch/17404971/3wy84lw/"
	a, err := lib.HTTPPostFormData(url, map[string]interface{}{
		"id":         0,
		"first_name": "F",
		"last_name":  "L",
		"email":      "lialing@foxmail.com",
	})
	lib.DPrintln(a)
	lib.DPrintln(err)
}
