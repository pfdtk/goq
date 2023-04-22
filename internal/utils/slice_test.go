package utils

import "testing"

func TestInSlice(t *testing.T) {
	r := InSlice(1, []int{2})
	println(r)
}
