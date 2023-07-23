package utils

func Clone[T any](slice []T) []T {
	n := make([]T, 0, len(slice))
	n = append(n, slice...)
	return n
}
