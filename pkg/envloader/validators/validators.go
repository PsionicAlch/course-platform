package validators

import (
	"fmt"
	"strconv"
)

type ValidationFunc func(name, data string) error

func ChainValidators(validatorFuncs ...ValidationFunc) ValidationFunc {
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

func EmptyValidator(name, data string) error {
	return nil
}

func NotEmptyValidator(name, data string) error {
	if data == "" {
		return fmt.Errorf("%s cannot be empty", name)
	}

	return nil
}

func IntValidator(name, data string) error {
	return IntegerValidator(name, data, 0)
}

func Int8Validator(name, data string) error {
	return IntegerValidator(name, data, 8)
}

func Int16Validator(name, data string) error {
	return IntegerValidator(name, data, 16)
}

func Int32Validator(name, data string) error {
	return IntegerValidator(name, data, 32)
}

func Int64Validator(name, data string) error {
	return IntegerValidator(name, data, 64)
}

func UintValidator(name, data string) error {
	return UnsignedIntegerValidator(name, data, 64)
}

func Uint8Validator(name, data string) error {
	return UnsignedIntegerValidator(name, data, 8)
}

func Uint16Validator(name, data string) error {
	return UnsignedIntegerValidator(name, data, 16)
}

func Uint32Validator(name, data string) error {
	return UnsignedIntegerValidator(name, data, 32)
}

func Uint64Validator(name, data string) error {
	return UnsignedIntegerValidator(name, data, 64)
}

func Float32Validator(name, data string) error {
	return FloatValidator(name, data, 32)
}

func Float64Validator(name, data string) error {
	return FloatValidator(name, data, 64)
}

func BoolValidator(name, data string) error {
	if _, err := strconv.ParseBool(data); err != nil {
		return fmt.Errorf("failed to convert %s to bool: %s", name, err)
	}

	return nil
}

func Complex64Validator(name, data string) error {
	return ComplexValidator(name, data, 64)
}

func Complex128Validator(name, data string) error {
	return ComplexValidator(name, data, 128)
}

func InSliceValidator(items []string) ValidationFunc {
	return func(name, data string) error {
		for _, item := range items {
			if item == data {
				return nil
			}
		}

		return fmt.Errorf("%s's %s is not in the list of viable options: %v", name, data, items)
	}
}
