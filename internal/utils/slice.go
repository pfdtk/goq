package utils

import (
	"github.com/pfdtk/goq/internal/types"
)

type InSliceT interface {
	types.Int | types.Uint | types.Float | ~string | ~bool
}

func InSlice[T InSliceT](item T, items []T) bool {
	for i := range items {
		if items[i] == item {
			return true
		}
	}
	return false
}
