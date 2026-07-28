package utils

// SliceFromSet converts a map[T]struct{} to a slice of T
func SliceFromSet[T comparable](m map[T]struct{}) []T {
	result := make([]T, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
