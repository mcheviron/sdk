package utils

import "slices"

// Filter applies a predicate function to each element in the input slice and returns a new slice
// containing only the elements for which the predicate function returns true.
func Filter[T any](s []T, pred func(T) bool) []T {
	var result []T
	for _, v := range s {
		if pred(v) {
			result = append(result, v)
		}
	}
	return result
}

// Unique returns a new slice containing unique elements from the input slice.
// The returned slice will preserve the order of the unique elements.
func Unique[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	for _, v := range s {
		seen[v] = struct{}{}
	}
	return slices.Clip(Keys(seen))
}

// Map applies a transformer function to each element in the input slice and returns a new slice
// containing the results, effictively morphing the input slice.
func Map[S ~[]E, E any, R any](s S, f func(E) R) []R {
	result := make([]R, 0, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// FilterMap applies a function that can error to each element of a slice and returns a new slice
// of the successful results. Courtesy of Rust
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

// Keys returns a slice of keys from the given map.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of values from the given map.
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
