package bulk

// mapKeys collects the keys from a map, not exported as modern go offers `slices.Collect(maps.Keys(m))`.
func mapKeys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
