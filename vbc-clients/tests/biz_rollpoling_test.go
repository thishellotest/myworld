package tests

import (
	"fmt"
	"testing"
)

func Test_RollpoingUsecase_Upsert(t *testing.T) {
	er := UT.RollpoingUsecase.Upsert("ss", "11")
	fmt.Println(er)
}
