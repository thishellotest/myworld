package lib

import (
	"fmt"
	"testing"
)

func Test_TypeMap_KeyExists(t *testing.T) {
	c := make(TypeMap)
	c.Set("aac.bb", 1)
	c.Set("aac.bbc", 222)
	c.Set("aac.bbc", "222335")
	//c.Set("aac.bbc.inffo", []string{"ssss", "cc"})
	c.Set("aac.bbc3.inffo", []string{"ssss", "cc"})
	//nestedMap := map[string]interface{}{
	//	"level1": map[string]interface{}{
	//		"level2": map[string]interface{}{
	//			"key1": "value1",
	//			"k2": map[string]interface{}{
	//				"aa": "cccc",
	//			},
	//		},
	//	},
	//}
	cccc := c.KeyExists("aac.bbc3.inffo1")
	DPrintln(cccc)
	DPrintln(c)
}

func Test_TypeMap_Set(t *testing.T) {
	c := make(TypeMap)
	c.Set("aac.bb", 1)
	c.Set("aac.bbc", 222)
	c.Set("aac.bbc", "222335")
	//c.Set("aac.bbc.inffo", []string{"ssss", "cc"})
	c.Set("aac.bbc3.inffo", []string{"ssss", "cc"})

	d := c.Get("aac.bbc.ss")
	fmt.Println("===:", d)
	DPrintln(c)
}

func Test_TypeMapMerge(t *testing.T) {

	c := make(TypeMap)
	c.Set("aac.bb", 1)
	c.Set("aac.bbc", 222)
	c.Set("aac.bbc", "222335")
	//c.Set("aac.bbc.inffo", []string{"ssss", "cc"})
	c.Set("aac.bbc3.inffo", []string{"ssss", "cc"})

	b := make(TypeMap)
	b.Set("bb", "11")
	b.Set("bb2", "112")

	dd := TypeMapMerge(c, b)
	DPrintln(dd)
}
