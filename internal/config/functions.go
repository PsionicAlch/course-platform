package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetVariable[T any](name string) (T, error) {
	var empty T

	variable, variableAvailable := os.LookupEnv(name)
	if !variableAvailable {
		return empty, fmt.Errorf("%s environment variable was not set", name)
	}

	switch any(empty).(type) {
	case string:
		return any(variable).(T), nil

	case int:
		v, err := strconv.Atoi(variable)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to int: %v", name, err)
		}

		return any(v).(T), nil

	case int8:
		v, err := strconv.ParseInt(variable, 10, 8)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to int8: %v", name, err)
		}

		return any(int8(v)).(T), nil

	case int16:
		v, err := strconv.ParseInt(variable, 10, 16)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to int16: %v", name, err)
		}

		return any(int16(v)).(T), nil

	case int32:
		v, err := strconv.ParseInt(variable, 10, 32)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to int32: %v", name, err)
		}

		return any(int32(v)).(T), nil

	case int64:
		v, err := strconv.ParseInt(variable, 10, 64)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to int64: %v", name, err)
		}

		return any(v).(T), nil

	case uint:
		v, err := strconv.ParseUint(variable, 10, 64)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to uint: %v", name, err)
		}

		return any(uint(v)).(T), nil

	case uint8:
		v, err := strconv.ParseUint(variable, 10, 8)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to uint8: %v", name, err)
		}

		return any(uint8(v)).(T), nil

	case uint16:
		v, err := strconv.ParseUint(variable, 10, 16)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to uint16: %v", name, err)
		}

		return any(uint16(v)).(T), nil

	case uint32:
		v, err := strconv.ParseUint(variable, 10, 32)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to uint32: %v", name, err)
		}

		return any(uint32(v)).(T), nil

	case uint64:
		v, err := strconv.ParseUint(variable, 10, 64)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to uint64: %v", name, err)
		}

		return any(v).(T), nil

	case float32:
		v, err := strconv.ParseFloat(variable, 32)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to float32: %v", name, err)
		}

		return any(float32(v)).(T), nil

	case float64:
		v, err := strconv.ParseFloat(variable, 64)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to float64: %v", name, err)
		}

		return any(v).(T), nil

	case bool:
		v, err := strconv.ParseBool(variable)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to bool: %v", name, err)
		}

		return any(v).(T), nil

	case complex64:
		v, err := strconv.ParseComplex(variable, 64)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to complex64: %v", name, err)
		}

		return any(complex64(v)).(T), nil

	case complex128:
		v, err := strconv.ParseComplex(variable, 128)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to complex128: %v", name, err)
		}

		return any(v).(T), nil

	case time.Duration:
		v, err := strconv.ParseInt(variable, 10, 64)
		if err != nil {
			return empty, fmt.Errorf("failed to convert %s to int64: %v", name, err)
		}

		return any(time.Duration(v)).(T), nil
	default:
		return empty, fmt.Errorf("unsupported type for %s", name)
	}
}
