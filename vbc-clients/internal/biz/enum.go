package biz

type EnumMap[T int | int32 | int64 | string] map[T]string

func (c EnumMap[T]) Name(k T) string {
	if v, ok := c[k]; ok {
		return v
	}
	return ""
}
