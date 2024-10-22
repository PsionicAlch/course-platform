package envloader

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PsionicAlch/psionicalch-home/pkg/envloader/validators"
	"github.com/joho/godotenv"
)

// LoadEnvironment uses github.com/joho/godotenv to load the specified env files afterwards it ensures that
// all the required variables are set and validated. Just pass it a map containing the names of the environment
// variables that you want to have set as well as a validation function to ensure that the variable is in the
// correct format.
func LoadEnvironment(settings map[string]validators.ValidationFunc, filenames ...string) error {
	err := godotenv.Load(filenames...)
	if err != nil {
		return err
	}

	for name, validateFunc := range settings {
		variable, variableAvailable := os.LookupEnv(name)
		if !variableAvailable {
			return fmt.Errorf("%s environment variable was not set", name)
		}

		err = validateFunc(name, variable)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetVariable is a generic function that allows you to read an environment variable and automatically convert
// it to one of the following types: string, int, int8, int16, int32, int64, uint8, uint16, uint32, uint64,
// float32, float64, bool, complex64, complex128.
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

	default:
		return empty, fmt.Errorf("unsupported type for %s", name)
	}
}
