package utils

import "crypto/rand"

// InSlice checks to see if an item is in a list of items.
func InSlice[T comparable](item T, items []T) bool {
	for _, i := range items {
		if item == i {
			return true
		}
	}

	return false
}

func RandomByteSlice(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
