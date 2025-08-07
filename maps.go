package bulk

// MapUnion takes in a series of maps, adding the values into a single map. On conflicts last value wins.
func MapUnion[K comparable, V any](maps ...map[K]V) map[K]V {
	switch len(maps) {
	case 0:
		return map[K]V{}
	case 1:
		return maps[0]
	}

	var size int
	for _, m := range maps {
		size += len(m)
	}

	result := make(map[K]V, size)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// mapKeys collects the keys from a map, not exported as modern go offers `slices.Collect(maps.Keys(m))`.
func mapKeys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
