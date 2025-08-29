package decimal

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

func Test_IntPart(t *testing.T) {

	de, _ := decimal.NewFromString("13312.233")
	r := de.Floor().IntPart()
	fmt.Println(r)
}

func Test_decimal(t *testing.T) {
	aaa := float64(1.2)
	cc := decimal.NewFromFloat(aaa)
	cc = cc.Add(decimal.NewFromFloat(2.44))
	ss := cc.String()
	fmt.Println(ss)
}
