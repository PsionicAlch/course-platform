package utils_test

import (
	"testing"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

func TestInSlice(t *testing.T) {
	item := "development"
	items := []string{"development", "testing", "production"}

	if !utils.InSlice(item, items) {
		t.Fatal("\"development\" was not found in slice containing \"development\"")
	}
}
