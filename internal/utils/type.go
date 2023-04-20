package utils

func InterfaceSlice[T any, R any](slice []T) []R {
	res := make([]R, len(slice))
	for i, v := range slice {
		res[i] = R(v)
	}
	return res
}
