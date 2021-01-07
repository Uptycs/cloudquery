package utilities

import (
	"encoding/json"
	"strconv"
)

func GetStringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	switch value.(type) {
	case string:
		return value.(string)
	case int:
		return strconv.FormatFloat(float64(value.(int)), 'g', -1, 64)
	case int8:
		return strconv.FormatFloat(float64(value.(int8)), 'g', -1, 64)
	case int16:
		return strconv.FormatFloat(float64(value.(int16)), 'g', -1, 64)
	case int32:
		return strconv.FormatFloat(float64(value.(int32)), 'g', -1, 64)
	case int64:
		return strconv.FormatFloat(float64(value.(int64)), 'g', -1, 64)
	case uint:
		return strconv.FormatFloat(float64(value.(uint)), 'g', -1, 64)
	case uint8:
		return strconv.FormatFloat(float64(value.(uint8)), 'g', -1, 64)
	case uint16:
		return strconv.FormatFloat(float64(value.(uint16)), 'g', -1, 64)
	case uint32:
		return strconv.FormatFloat(float64(value.(uint32)), 'g', -1, 64)
	case uint64:
		return strconv.FormatFloat(float64(value.(uint64)), 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'g', -1, 64)
	case float64:
		return strconv.FormatFloat(value.(float64), 'g', -1, 64)
	case json.Number:
		val, _ := value.(json.Number).Int64()
		return strconv.FormatInt(val, 10)
	case bool:
		return strconv.FormatBool(value.(bool))
	}

	return ""
}

func GetFloat64Value(value interface{}) float64 {
	if value == nil {
		return 0
	}
	switch value.(type) {
	case string:
		num, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return 0
		}
		return num
	case int:
		return float64(value.(int))
	case int8:
		return float64(value.(int8))
	case int16:
		return float64(value.(int16))
	case int32:
		return float64(value.(int32))
	case int64:
		return float64(value.(int64))
	case uint:
		return float64(value.(uint))
	case uint8:
		return float64(value.(uint8))
	case uint16:
		return float64(value.(uint16))
	case uint32:
		return float64(value.(uint32))
	case uint64:
		return float64(value.(uint64))
	case float32:
		return float64(value.(float32))
	case float64:
		return value.(float64)
	case json.Number:
		val, err := value.(json.Number).Float64()
		if err != nil {
			return 0
		}
		return val
	}
	// type which we can't operate on.
	return 0
}

func GetIntegerValue(value interface{}) int {
	if value == nil {
		return 0
	}
	switch value.(type) {
	case string:
		num, err := strconv.ParseInt(value.(string), 10, 0)
		if err != nil {
			return 0
		}
		return int(num)
	case int:
		return value.(int)
	case int8:
		return int(value.(int8))
	case int16:
		return int(value.(int16))
	case int32:
		return int(value.(int32))
	case int64:
		return int(value.(int64))
	case uint:
		return int(value.(uint))
	case uint8:
		return int(value.(uint8))
	case uint16:
		return int(value.(uint16))
	case uint32:
		return int(value.(uint32))
	case uint64:
		return int(value.(uint64))
	case float32:
		return int(value.(float32))
	case float64:
		return int(value.(float64))
	case json.Number:
		val, err := value.(json.Number).Int64()
		if err != nil {
			return 0
		}
		return int(val)
	}
	// type which we can't operate on.
	return 0
}
