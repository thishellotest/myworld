package lib

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ToTypeMapByString(val string) TypeMap {
	r := StringToT[TypeMap](val)
	return r.UnwrapOr(nil)
}

func ToTypeListByString(val string) (typeList TypeList) {
	r := StringToT[TypeList](val)
	return r.UnwrapOr(nil)
}
func ToTypeList(params interface{}) (typeList TypeList) {
	defer func() {
		if err := recover(); err != nil {
			typeList = nil
		}
	}()
	records := params.([]interface{})

	for k, _ := range records {
		tt := records[k].(map[string]interface{})
		typeList = append(typeList, tt)
	}
	return
}

func ToTypeMap(params interface{}) (typeMap TypeMap) {
	defer func() {
		if err := recover(); err != nil {
			typeMap = nil
		}
	}()
	return params.(map[string]interface{})
}

func ToTypeListInterface(params interface{}) (typeListInterface TypeListInterface) {
	defer func() {
		if err := recover(); err != nil {
			typeListInterface = nil
		}
	}()
	records := params.([]interface{})
	return records
}

type TypeListInterface []interface{}

func ConvertTypeListInterface[T string | int32](input TypeListInterface) ([]T, error) {
	result := make([]T, 0, len(input))
	for _, v := range input {
		switch val := v.(type) {
		case string:
			if _, ok := any((*T)(nil)).(*string); ok {
				result = append(result, any(val).(T))
			} else {
				return nil, fmt.Errorf("type mismatch: expected %T, got string", *new(T))
			}
		case int32:
			if _, ok := any((*T)(nil)).(*int32); ok {
				result = append(result, any(val).(T))
			} else {
				return nil, fmt.Errorf("type mismatch: expected %T, got int32", *new(T))
			}
		default:
			return nil, fmt.Errorf("unsupported type: %T", v)
		}
	}
	return result, nil
}

type TypeList []TypeMap

func (c *TypeList) AppendList(list TypeList) {
	for k, _ := range list {
		*c = append(*c, list[k])
	}
}

type TypeMap map[string]interface{}

func (c TypeMap) ToBytes() []byte {
	a, _ := json.Marshal(&c)
	return a
}

func (c TypeMap) ToOrigin() map[string]interface{} {
	return c
}

func (c TypeMap) ToString() string {
	return string(c.ToBytes())
}

func (c TypeMap) GetInt64(key string) int64 {
	r := c.Get(key)
	if r == nil {
		return 0
	}
	a, _ := strconv.ParseInt(InterfaceToString(r), 10, 64)
	return a
}

func (c TypeMap) GetInt(key string) int32 {
	r := c.Get(key)
	if r == nil {
		return 0
	}
	a, _ := strconv.ParseInt(InterfaceToString(r), 10, 32)
	return int32(a)
}

func (c TypeMap) GetTypeMap(key string) TypeMap {
	r := c.Get(key)
	if r == nil {
		return nil
	}
	return ToTypeMap(r)
}

func (c TypeMap) GetString(key string) string {
	r := c.Get(key)
	if r == nil {
		return ""
	}
	return InterfaceToString(r)
}

func (c TypeMap) GetTypeList(key string) TypeList {
	r := c.Get(key)
	if r == nil {
		return nil
	}
	return ToTypeList(r)
}

func (c TypeMap) GetTypeListInterface(key string) TypeListInterface {
	r := c.Get(key)
	if r == nil {
		return nil
	}
	return ToTypeListInterface(r)
}

func (c TypeMap) KeyExists(linkKey string) bool {
	var m interface{} = c
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		return false
	}
	keys := strings.Split(linkKey, ".")
	for i, key := range keys {
		v = v.MapIndex(reflect.ValueOf(key))
		if v.IsValid() {
			if v.Kind() == reflect.Map {
				// 如果值是map类型，继续在这个map上查找
				m = v.Interface()
				if i+1 == len(keys) {
					return true
				}
				continue
			} else if v.Kind() == reflect.Interface {
				// 如果是接口类型，则进一步解析接口的值
				v = v.Elem()
				if v.IsValid() && v.Kind() == reflect.Map {
					m = v.Interface()
					if i+1 == len(keys) {
						return true
					}
					continue
				}
			}
			if i+1 == len(keys) {
				return true
			} else {
				return false
			}
		}
		return false
	}
	return false
}

// Get key data.email
func (c TypeMap) Get(key string) (r interface{}) {

	defer func() {
		if err := recover(); err != nil {
			r = nil
			fmt.Println(err)
		}
	}()

	keys := strings.Split(key, ".")
	le := len(keys)
	if le == 0 {
		return nil
	} else if le == 1 {
		return c[keys[0]]
	} else {
		newKeys := keys[1:le]
		str := strings.Join(newKeys, ".")
		var destTypeMap TypeMap
		switch c[keys[0]].(type) {
		case TypeMap:
			destTypeMap = c[keys[0]].(TypeMap)
		case map[string]interface{}:
			dd := c[keys[0]].(map[string]interface{})
			destTypeMap = dd
		}

		return destTypeMap.Get(str)
	}
}

func (c TypeMap) Set(key string, val interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("TypeMap Set", err)
		}
	}()

	keys := strings.Split(key, ".")
	le := len(keys)
	if le == 1 {
		if c == nil {
			c = make(map[string]interface{})
		}
		c[keys[0]] = val
	} else {
		newKeys := keys[1:le]
		newStr := strings.Join(newKeys, ".")
		//fmt.Println("_1", newStr)
		var cMap map[string]interface{}
		tt := c.Get(keys[0])
		if tt == nil {
			//fmt.Println("_2")
			cMap = make(map[string]interface{})
			c.Set(keys[0], cMap)
		} else {
			cMap = tt.(map[string]interface{})
		}
		TypeMap(cMap).Set(newStr, val)
	}
}

// TypeMapMerge 后面的map覆盖前面的map
func TypeMapMerge(val ...TypeMap) (res TypeMap) {
	if len(val) > 0 {
		for k, _ := range val {
			if val[k] != nil {
				if res == nil {
					res = make(TypeMap)
				}
				for k1, _ := range val[k] {
					res[k1] = val[k][k1]
				}
			}
		}
	}
	return
}
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
