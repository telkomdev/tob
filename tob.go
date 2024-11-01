package tob

import (
	"reflect"
)

const (
	// Version number

	Version = "2.0.5"

	// OK service status
	OK = "OK"

	// NotOk service status
	NotOk = "NOT_OK"
)

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
