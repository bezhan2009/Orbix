package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

// ConvertToString пытается преобразовать значение в строку
func ConvertToString(value reflect.Value) (string, error) {
	switch value.Kind() {
	case reflect.String:
		return value.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), nil
	default:
		return "", fmt.Errorf("cannot convert %s to string", value.Kind())
	}
}

// ConvertToInt пытается преобразовать значение в целое число
func ConvertToInt(value reflect.Value) (int64, error) {
	switch value.Kind() {
	case reflect.String:
		return strconv.ParseInt(value.String(), 10, 64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int(), nil
	case reflect.Float32, reflect.Float64:
		return int64(value.Float()), nil
	default:
		return 0, fmt.Errorf("cannot convert %s to int", value.Kind())
	}
}

// ConvertToBool пытается преобразовать значение в булевое
func ConvertToBool(value reflect.Value) (bool, error) {
	switch value.Kind() {
	case reflect.String:
		return strconv.ParseBool(value.String())
	case reflect.Bool:
		return value.Bool(), nil
	default:
		return false, fmt.Errorf("cannot convert %s to bool", value.Kind())
	}
}
