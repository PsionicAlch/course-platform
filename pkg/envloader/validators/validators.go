package validators

import (
	"fmt"
	"slices"
	"strconv"
)

type ValidationFunc func(name, data string) error

func Chain(validatorFuncs ...ValidationFunc) ValidationFunc {
	return func(name, data string) error {
		for _, validator := range validatorFuncs {
			err := validator(name, data)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func Empty(name, data string) error {
	return nil
}

func NotEmpty(name, data string) error {
	if data == "" {
		return fmt.Errorf("%s cannot be empty", name)
	}

	return nil
}

func Int(name, data string) error {
	return Integer(name, data, 0)
}

func Int8(name, data string) error {
	return Integer(name, data, 8)
}

func Int16(name, data string) error {
	return Integer(name, data, 16)
}

func Int32(name, data string) error {
	return Integer(name, data, 32)
}

func Int64(name, data string) error {
	return Integer(name, data, 64)
}

func Uint(name, data string) error {
	return UnsignedInteger(name, data, 64)
}

func Uint8(name, data string) error {
	return UnsignedInteger(name, data, 8)
}

func Uint16(name, data string) error {
	return UnsignedInteger(name, data, 16)
}

func Uint32(name, data string) error {
	return UnsignedInteger(name, data, 32)
}

func Uint64(name, data string) error {
	return UnsignedInteger(name, data, 64)
}

func Float32(name, data string) error {
	return Float(name, data, 32)
}

func Float64(name, data string) error {
	return Float(name, data, 64)
}

func Bool(name, data string) error {
	if _, err := strconv.ParseBool(data); err != nil {
		return fmt.Errorf("failed to convert %s to bool: %s", name, err)
	}

	return nil
}

func Complex64(name, data string) error {
	return Complex(name, data, 64)
}

func Complex128(name, data string) error {
	return Complex(name, data, 128)
}

func InSlice(items []string) ValidationFunc {
	return func(name, data string) error {
		if slices.Contains(items, data) {
			return nil
		}

		return fmt.Errorf("%s's %s is not in the list of viable options: %v", name, data, items)
	}
}
