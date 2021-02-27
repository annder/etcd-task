package utils

import (
	"strconv"
)

func Int(value interface{}) int {
	switch value.(type) {
	case string:
		v, _ := value.(string)
		rs, _ := strconv.Atoi(v)
		return rs
	case int:
		v, _ := value.(int)
		return v
	case int32:
		v, _ := value.(int32)
		return int(v)
	case int64:
		v, _ := value.(int64)
		return int(v)
	case float32:
		v, _ := value.(float32)
		return int(v)
	case float64:
		v, _ := value.(float64)
		return int(v)
	default:
		return int(0)
	}
}



