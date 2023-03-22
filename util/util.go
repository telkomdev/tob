package util

import (
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
