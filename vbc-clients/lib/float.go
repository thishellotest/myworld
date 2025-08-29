package lib

import "github.com/shopspring/decimal"

func Float32Sum(params ...float32) float64 {
	sum := decimal.New(0, 0)
	for _, v := range params {
		sum = sum.Add(decimal.NewFromFloat32(v))
	}
	r, _ := sum.Float64()
	return r
}

func FloatSum(params ...float64) float64 {
	sum := decimal.NewFromFloat(0)

	for _, v := range params {
		sum = sum.Add(decimal.NewFromFloat(v))
	}
	r, _ := sum.Float64()
	return r
}
