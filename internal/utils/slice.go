package utils

import "github.com/pfdtk/goq/common"

func InSlice[T common.Int | common.Uint | common.Float | string | bool](item T, items []T) bool {
	for i := range items {
		if items[i] == item {
			return true
		}
	}
	return false
}
