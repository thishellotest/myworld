package lib

func ArrayReverse[T any](params []T) (res []T) {
	if params != nil {
		res = make([]T, 0)
	}
	for i := len(params) - 1; i >= 0; i-- {
		res = append(res, params[i])
	}
	return res
}

func InArray[T comparable](val T, array []T) bool {
	for _, v := range array {
		if v == val {
			return true
		}
	}
	return false
}

// RemoveDuplicates 使用泛型去除数组中的重复元素
func RemoveDuplicates[T comparable](arr []T) []T {
	// 使用map来去除重复元素
	resultMap := make(map[T]bool)

	// 遍历数组，将每个元素加入map，重复的会被自动忽略
	for _, v := range arr {
		resultMap[v] = true
	}

	// 将map中的key放入结果切片
	var result []T
	for key := range resultMap {
		result = append(result, key)
	}

	return result
}

func ConvertToInterfaceSlice[T any](slice []T) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
