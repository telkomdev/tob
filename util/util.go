package util

import (
	"reflect"
	"strconv"
)

// InterfaceToFloat64 will convert any interface value to float64
func InterfaceToFloat64(val interface{}) float64 {

	var i float64
	switch t := val.(type) {
	case int:
		i = float64(t)
		break
	case int8:
		i = float64(t)
		break
	case int16:
		i = float64(t)
		break
	case int32:
		i = float64(t)
		break
	case int64:
		i = float64(t)
		break
	case float32:
		i = float64(t)
		break
	case float64:
		i = float64(t)
		break
	case uint8:
		i = float64(t)
		break
	case uint16:
		i = float64(t)
		break
	case uint32:
		i = float64(t)
		break
	case uint64:
		i = float64(t)
		break
	case string:
		num, err := strconv.Atoi(t)
		if err != nil {
			return 0.0
		}

		i = float64(num)
		break
	default:
		i = 0.0
	}

	return i
}

// IsNilish function is a function to check for nullability on the interface
func IsNilish(val any) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}

	return false
}
