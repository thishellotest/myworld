package lib

import (
	"encoding/json"
	"strconv"
	"strings"
)

// SqlBindValueForSelect 把sql的值转换
func SqlBindValueForSelect(value interface{}) string {
	var key string
	if value == nil {
		return "\"" + key + "\""
	}
	needBackslash := false
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		needBackslash = true
		key = value.(string)
	case []byte:
		needBackslash = true
		key = string(value.([]byte))
	case error:
		needBackslash = true
		key = value.(error).Error()
	default:
		needBackslash = true
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	if needBackslash {
		key = SqlValueBackslashForSelect(key)
		key = "\"" + key + "\""
	}

	return key
}

// SqlBindValue 把sql的值转换
func SqlBindValue(value interface{}) string {
	var key string
	if value == nil {
		return "null"
	}
	needBackslash := false
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		needBackslash = true
		key = value.(string)
	case []byte:
		needBackslash = true
		key = string(value.([]byte))
	case error:
		needBackslash = true
		key = value.(error).Error()
	default:
		needBackslash = true
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	if needBackslash {
		key = SqlValueBackslash(key)
		key = "\"" + key + "\""
	}

	return key
}

// SqlValueBackslash sql的值转义
func SqlValueBackslash(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "'", "\\'")
	return value
}

// SqlValueBackslashForSelect sql的值转义
func SqlValueBackslashForSelect(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "%", "\\%")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "'", "\\'")
	return value
}
