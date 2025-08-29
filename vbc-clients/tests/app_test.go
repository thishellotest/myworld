package tests

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	AppMain()
	m.Run()
	fmt.Println("end")
}

func TestCond(t *testing.T) {

}
