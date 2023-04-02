package internal

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/spf13/cast"
)

// Deref --
func Deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}

// GetBaseType --
func GetBaseType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = Deref(t)

	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}

	return t, nil
}

// CastValueTo -- convert value from string to specific reflect.Type
func CastValueTo(v interface{}, vType reflect.Type, ptr bool) interface{} {
	strValue := cast.ToString(v)

	switch vType.Kind() {
	case reflect.String:
		if ptr {
			return StringPtr(strValue)
		}

		return strValue
	case reflect.Bool:
		defaultValue, _ := strconv.ParseBool(strValue)

		if ptr {
			return BoolPtr(defaultValue)
		}

		return defaultValue
	case reflect.Float64:
		defaultValue, _ := strconv.ParseFloat(strValue, 64)

		if ptr {
			return Float64Ptr(defaultValue)
		}

		return defaultValue
	case reflect.Float32:
		defaultValue, _ := strconv.ParseFloat(strValue, 32)

		if ptr {
			return Float64Ptr(defaultValue)
		}

		return defaultValue
	case reflect.Int:
		defaultValue, _ := strconv.ParseInt(strValue, 10, 0)

		if ptr {
			return Int64Ptr(defaultValue)
		}

		return defaultValue
	case reflect.Int8:
		defaultValue, _ := strconv.ParseInt(strValue, 10, 8)

		if ptr {
			return Int64Ptr(defaultValue)
		}

		return defaultValue
	case reflect.Int16:
		defaultValue, _ := strconv.ParseInt(strValue, 10, 16)

		if ptr {
			return Int64Ptr(defaultValue)
		}

		return defaultValue
	case reflect.Int32:
		defaultValue, _ := strconv.ParseInt(strValue, 10, 32)

		if ptr {
			return Int64Ptr(defaultValue)
		}

		return defaultValue
	case reflect.Int64:
		defaultValue, _ := strconv.ParseInt(strValue, 10, 64)

		if ptr {
			return Int64Ptr(defaultValue)
		}

		return defaultValue
	default:
		return nil
	}
}

// IsZeroOfUnderlyingType --
func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
