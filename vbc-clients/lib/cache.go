package lib

/*
Cache[T] 是一个类型别名，不是指针类型：

在 UT.TUsecase.DataByGidWithCaches 中，如果 caches 被按值传递给函数，则会生成 Cache[T] 的副本。此时，任何修改都只影响副本，而不会反映到原始 caches 上。
结果就是 map 的地址看起来可能发生了改变，因为不同函数内部操作的是不同的 Cache[T] 副本。
map 的底层地址可能发生变化：

如果对 map 的容量没有合理的初始化，频繁向 map 中添加数据可能导致底层存储地址发生变化。这通常不是主要问题，但可能会让调试变得混乱。

改动点：
Cache[T] 的所有方法都改为指针接收器（*Cache[T]）。
CacheInit 返回指针类型。
这样可以保证在所有地方操作的都是同一个 Cache 实例。
*/

type Cache[T any] map[string]T

func (c *Cache[T]) Get(key string) (entity T, exist bool) {
	if v, ok := (*c)[key]; ok {
		return v, true
	}
	exist = false
	return
}

func (c *Cache[T]) Set(key string, entity T) {
	(*c)[key] = entity
}

func CacheInit[T any]() Cache[T] {
	return make(Cache[T])
}
