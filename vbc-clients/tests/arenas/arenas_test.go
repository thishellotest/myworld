package arenas

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_aaac(t *testing.T) {
	// GOEXPERIMENT=arenas
	fmt.Println("sssc")
}

func processRequest(req *http.Request) {
	//// 开始创建公共arena内存池
	//mem := arena.NewArena()
	//// 最后统一释放内存池
	//defer mem.Free()
	//
	//// 分配一系列单对象
	//for i := 0; i < 10; i++ {
	//	obj := arena.New[T](mem)
	//	obj.Foo = "Hello"
	//
	//	fmt.Printf("%v\n", obj)
	//}
	//
	//// 或者分配slice 暂时不支持map
	//// 参数 mem, length, capacity
	//slice := arena.MakeSlice[T](mem, 100, 200)
	//slice[0].Foo = "hello"
	//fmt.Printf("%v\n", slice)
	//
	//// 不能直接分配string，可借助bytes转换
	//src := "source string"
	//
	//bs := arena.MakeSlice[byte](mem, len(src), len(src))
	//copy(bs, src)
	//str := unsafe.String(&bs[0], len(bs))
	//
	//fmt.Printf("%v\n", str)
}
