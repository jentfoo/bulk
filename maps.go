package bulk

// MapInvert swaps keys and values in a map. If duplicate values exist, the resulting key is nondeterministic.
func MapInvert[K comparable, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))
	MapInvertInto(result, m)
	return result
}

// MapInvertInto swaps keys and values from the source map into the destination map.
// If duplicate values exist, the resulting key is nondeterministic.
func MapInvertInto[K comparable, V comparable](dest map[V]K, m map[K]V) {
	for key, value := range m {
		dest[value] = key
	}
}
