package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// InSlice checks to see if an item is in a list of items.
func InSlice[T comparable](item T, items []T) bool {
	for _, i := range items {
		if item == i {
			return true
		}
	}

	return false
}

// InSliceFunc checks to see if an item is in a list of items and uses a user passed function
// to do the comparison.
func InSliceFunc[T, A any](item T, items []A, compareFunc func(itemA T, itemB A) bool) (int, bool) {
	for index, i := range items {
		if compareFunc(item, i) {
			return index, true
		}
	}

	return -1, false
}

func RandomByteSlice(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func EncodeString(src []byte) string {
	return base64.RawStdEncoding.EncodeToString(src)
}

func DecodeString(src string) ([]byte, error) {
	return base64.RawStdEncoding.Strict().DecodeString(src)
}

func Filter[S ~[]E, E any](items S, compareFunc func(E) bool) S {
	var collection S

	for _, item := range items {
		if compareFunc(item) {
			collection = append(collection, item)
		}
	}

	return collection
}
