package a

import (
	"fmt"
	"regexp"
	"testing"
	"time"
	"vbc/lib"
)

func Test_aaaaa(t *testing.T) {
	a := "2021-12-02"
	b := "2022-01-01"
	if a > b {
		lib.DPrintln("a > b")
	} else {
		lib.DPrintln("no")
	}
}

func Test_Time_Format(t *testing.T) {
	a := time.Now().Format("Jan 02, 2006")
	lib.DPrintln("sss_c 123456", a)
}

type Person interface {
	SetName(name string)
}

type Employee struct {
	Name string
}

func (e *Employee) SetName(name string) {
	e.Name = name
}

func CreatePerson[T Person]() *Employee {
	obj := new(Employee) // 强制使用指针类型
	lib.DPrintln(obj, "ssss")
	return obj
}

func TestB(t *testing.T) {
	emp := CreatePerson[*Employee]()
	(*emp).SetName("Alice")
	fmt.Println((*emp).Name) // 输出：Alice

	ssss := new(*Employee)
	(*ssss).Name = "ss"
	lib.DPrintln(ssss)
}

func TestC(t *testing.T) {

	text := `
#abc
#    Name of Disability/Condition: Low back pain secondary to left knee meniscal tear s/p arthroscopic surgery
Current Treatment Facility: 
Current Medication: Tylenol, Ibuprofen

I am respectfully requesting Veteran Affairs benefits for my condition of low back pain secondary to left knee meniscal tear s/p arthroscopic surgery which began during my service. I served in the United States Navy from 2004 to 2024 as a Gunnersmate. During this service period, I developed low back pain secondary to left knee meniscal tear s/p arthroscopic surgery that continues to affect my daily life and ability to work.

##     Onset and Service Connection:`

	// 正则：行首，0个或多个空格，1或2个#，至少1个空格
	re := regexp.MustCompile(`(?m)^\s*#{1,2}\s*`)

	cleanText := re.ReplaceAllString(text, "")

	fmt.Println(cleanText)

}
