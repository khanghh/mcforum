package arrays

import (
	"strings"
)

// func Contains(obj interface{}, arr interface{}) bool {
// 	targetValue := reflect.ValueOf(arr)
// 	switch reflect.TypeOf(arr).Kind() {
// 	case reflect.Slice, reflect.Array:
// 		for i := 0; i < targetValue.Len(); i++ {
// 			if targetValue.Index(i).Interface() == obj {
// 				return true
// 			}
// 		}
// 	case reflect.Map:
// 		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
// 			return true
// 		}
// 	}
// 	return false
// }

func Contains[T comparable](arr []T, target T) bool {
	for _, val := range arr {
		if val == target {
			return true
		}
	}
	return false
}

func ContainsIgnoreCase(arr []string, str string) bool {
	if len(str) == 0 {
		return false
	}
	if len(arr) == 0 {
		return false
	}
	str = strings.ToLower(str)
	for i := 0; i < len(arr); i++ {
		if strings.ToLower(arr[i]) == str {
			return true
		}
	}
	return false
}

func Distinct[T any](input []T, getKey func(T) any) (output []T) {
	tempMap := map[any]byte{}
	for _, item := range input {
		l := len(tempMap)
		tempMap[getKey(item)] = 0
		if len(tempMap) != l { // 数量发生变化，说明不存在
			output = append(output, item)
		}
	}
	return
}
