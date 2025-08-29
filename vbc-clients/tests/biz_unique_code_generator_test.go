package tests

import (
	"fmt"
	"testing"
	"vbc/internal/biz"
)

func Test_UniqueCodeGeneratorUsecase_t(t *testing.T) {
	i := 0
	for {
		i++
		a, err := UT.UniqueCodeGeneratorUsecase.GenUuid(biz.UniqueCodeGenerator_Type_ClientUniqCode, 0)
		fmt.Println(a, err)
		if err != nil {
			panic(err)
		}
		if i > 10 {
			break
		}
	}
}
