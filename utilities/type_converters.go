package utilities

import (
	"encoding/json"
	"strconv"
	"strings"
)

func getStringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	switch value.(type) {
	case string:
		strIn := []rune(value.(string))
		if len(strIn) >= 2 && strIn[0] == '"' && strIn[len(strIn)-1] == '"' {
			// Remove quotes
			noLeadingQuotes := strIn[1:]
			return string(noLeadingQuotes[:len(noLeadingQuotes)-1])
		}
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

func getFloat64Value(value interface{}) float64 {
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

func getIntegerValue(value interface{}) int {
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

// Get boolean value of given variable
func getBooleanValue(value interface{}) bool {
	if value == nil {
		return false
	}

	switch value.(type) {
	case bool:
		return value.(bool)
	case string:
		var str = value.(string)
		if strings.EqualFold(str, "true") || strings.EqualFold(str, "yes") {
			return true
		} else {
			return false
		}
	case int:
		if value.(int) > 0 {
			return true
		} else {
			return false
		}
	case int8:
		if value.(int8) > 0 {
			return true
		} else {
			return false
		}
	case int16:
		if value.(int16) > 0 {
			return true
		} else {
			return false
		}
	case int32:
		if value.(int32) > 0 {
			return true
		} else {
			return false
		}
	case int64:
		if value.(int64) > 0 {
			return true
		} else {
			return false
		}
	case uint:
		if value.(uint) > 0 {
			return true
		} else {
			return false
		}
	case uint8:
		if value.(uint8) > 0 {
			return true
		} else {
			return false
		}
	case uint16:
		if value.(uint16) > 0 {
			return true
		} else {
			return false
		}
	case uint32:
		if value.(uint32) > 0 {
			return true
		} else {
			return false
		}
	case uint64:
		if value.(uint64) > 0 {
			return true
		} else {
			return false
		}
	}
	return true
}
