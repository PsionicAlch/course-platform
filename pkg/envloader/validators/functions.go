package validators

import (
	"fmt"
	"strconv"
)

func Integer(name, data string, bitSize int) error {
	if _, err := strconv.ParseInt(data, 10, bitSize); err != nil {
		if bitSize == 0 {
			return fmt.Errorf("failed to convert %s to int: %s", name, err)
		} else {
			return fmt.Errorf("failed to convert %s to int%d: %s", name, bitSize, err)
		}
	}

	return nil
}

func UnsignedInteger(name, data string, bitSize int) error {
	if _, err := strconv.ParseUint(data, 10, bitSize); err != nil {
		return fmt.Errorf("failed to convert %s to uint%d: %s", name, bitSize, err)
	}

	return nil
}

func Float(name, data string, bitSize int) error {
	if _, err := strconv.ParseFloat(data, bitSize); err != nil {
		return fmt.Errorf("failed to convert %s to float%d: %s", name, bitSize, err)
	}

	return nil
}

func Complex(name, data string, bitSize int) error {
	if _, err := strconv.ParseComplex(data, bitSize); err != nil {
		return fmt.Errorf("failed to convert %s to complex%d: %s", name, bitSize, err)
	}

	return nil
}
