package internal

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func Getenv[T any](key string, fallback T) (T, error) {
	var result T
	valueType := reflect.TypeOf(result)
	valuePtr := reflect.New(valueType)
	returnValue := valuePtr.Elem()

	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}
	switch returnValue.Kind() {
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return result, err
		}
		returnValue.SetInt(int64(intValue))
	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return result, err
		}
		returnValue.SetFloat(floatValue)
	case reflect.String:
		returnValue.SetString(value)
	case reflect.Struct:
		return result, fmt.Errorf("struct casting not supported")
	default:
		return result, fmt.Errorf("unsupported kind: %s", returnValue.Kind())
	}

	return returnValue.Interface().(T), nil
}
