package utils

import "slices"

func Filter[T any](s []T, pred func(T) bool) []T {
	var result []T
	for _, v := range s {
		if pred(v) {
			result = append(result, v)
		}
	}
	return result
}

func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0)
	for _, v := range slice {
		if _, found := seen[v]; !found {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func Map[S ~[]E, E any, R any](s S, f func(E) R) []R {
	result := make([]R, 0, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

func FilterMap[S ~[]E, E any, R any](s S, f func(E) (R, error)) []R {
	result := make([]R, 0, len(s))
	for _, v := range s {
		r, err := f(v)
		if err == nil {
			result = append(result, r)
		}
	}
	return slices.Clip(result)
}
