package utils

import (
	"iter"
)

// InSliceFunc checks to see if an item is in a list of items and uses a user passed function
// to do the comparison.
func InSliceFunc[T comparable, A any](item T, items []A, compareFunc func(itemA T, itemB A) bool) (int, bool) {
	for index, i := range items {
		if compareFunc(item, i) {
			return index, true
		}
	}

	return -1, false
}

// InSeq checks to see if an item is present in a given sequence.
func InSeq[T iter.Seq[A], A comparable](item A, items T) bool {
	for i := range items {
		if item == i {
			return true
		}
	}

	return false
}

// Find uses a user provided function to linearly search through a slice to find the index of an element.
func Find[A comparable](items []A, compareFunc func(item A) bool) (int, bool) {
	for index, item := range items {
		if compareFunc(item) {
			return index, true
		}
	}

	return -1, false
}
