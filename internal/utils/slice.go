package utils

import "github.com/pfdtk/goq/common"

type InSliceT interface {
	common.Int | common.Uint | common.Float | ~string | ~bool
}

func InSlice[T InSliceT](item T, items []T) bool {
	for i := range items {
		if items[i] == item {
			return true
		}
	}
	return false
}
