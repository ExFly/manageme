package util

import (
	"reflect"
	"time"
)

// IsEmpty https://github.com/issue9/assert/blob/master/util.go
func IsEmpty(expr interface{}) bool {
	if expr == nil {
		return true
	}
	switch v := expr.(type) {
	case bool:
		return !v
	case int:
		return 0 == v
	case int8:
		return 0 == v
	case int16:
		return 0 == v
	case int32:
		return 0 == v
	case int64:
		return 0 == v
	case uint:
		return 0 == v
	case uint8:
		return 0 == v
	case uint16:
		return 0 == v
	case uint32:
		return 0 == v
	case uint64:
		return 0 == v
	case string:
		return len(v) == 0
	case float32:
		return 0 == v
	case float64:
		return 0 == v
	case time.Time:
		return v.IsZero()
	case *time.Time:
		return v.IsZero()
	}
	// 符合 IsNil 条件的，都为 Empty
	if IsNil(expr) {
		return true
	}

	// 长度为 0 的数组也是 empty
	v := reflect.ValueOf(expr)
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		return 0 == v.Len()
	}

	return false
}

// IsNil 判断一个值是否为 nil。
// 当特定类型的变量，已经声明，但还未赋值时，也将返回 true
func IsNil(expr interface{}) bool {
	if nil == expr {
		return true
	}

	v := reflect.ValueOf(expr)
	k := v.Kind()

	return k >= reflect.Chan && k <= reflect.Slice && v.IsNil()
}
