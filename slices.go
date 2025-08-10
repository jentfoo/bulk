package bulk

// SliceFilter returns elements that pass the predicate function.
// May return the original slice if all elements pass (no allocation).
func SliceFilter[T any](slice []T, predicate func(v T) bool) []T {
	for falseIndex, v := range slice {
		if predicate(v) {
			continue // continue till first false is found
		}

		// build and return result
		if falseIndex == 0 {
			// Track state transitions for view optimization
			firstTrueIndex := -1

			// Find first true element
			for i := falseIndex + 1; i < len(slice); i++ {
				if predicate(slice[i]) {
					firstTrueIndex = i
					break
				}
			}
			if firstTrueIndex == -1 {
				return slice[:0] // No true elements found
			}

			// Check if all remaining elements are consecutive and true
			consecutiveEnd := firstTrueIndex
			nonConsecutiveStart := -1
			for i := firstTrueIndex + 1; i < len(slice); i++ {
				if predicate(slice[i]) {
					consecutiveEnd = i
					continue
				}
				// Found a false element, break to check if consecutive section continues
				nonConsecutiveStart = i
				break
			}
			if nonConsecutiveStart < 0 {
				return slice[firstTrueIndex:] // All elements from firstTrueIndex to end are true
			}

			// if any more trues, we have to allocate and append, otherwise return a view
			for j := nonConsecutiveStart + 1; j < len(slice); j++ {
				if predicate(slice[j]) {
					// Found another true after false, not consecutive - need to allocate
					// worst case size: (consecutiveEnd-firstTrueIndex+1) + 1 + (len(slice) - j - 1) (+1 -1 simplified out)
					result := make([]T, 0, (consecutiveEnd-firstTrueIndex+1)+capGuess(len(slice)-j))
					result = append(result, slice[firstTrueIndex:consecutiveEnd+1]...)
					result = append(result, slice[j])

					// Continue appending remaining true elements
					for k := j + 1; k < len(slice); k++ {
						if predicate(slice[k]) {
							result = append(result, slice[k])
						}
					}
					return result
				}
			}
			// No more true elements found, return consecutive view
			return slice[firstTrueIndex : consecutiveEnd+1]
		} else { // Started true, now first false found
			// Find first true element after falseIndex
			secondTrueIndex := -1
			for i := falseIndex + 1; i < len(slice); i++ {
				if predicate(slice[i]) {
					secondTrueIndex = i
					break
				}
			}
			if secondTrueIndex < 0 {
				return slice[:falseIndex] // No true elements in suffix, return prefix only
			}

			// true+ -> false+ -> true - We must allocate at this point
			result := make([]T, 0, falseIndex+capGuess(len(slice)-secondTrueIndex))
			result = append(result, slice[:falseIndex]...)
			result = append(result, slice[secondTrueIndex])
			for i := secondTrueIndex + 1; i < len(slice); i++ {
				if predicate(slice[i]) {
					result = append(result, slice[i])
				}
			}
			return result
		}
	}
	return slice // all records tested to true
}

// capGuess attempts to guess the sizing for an allocation, large allocations are reduced.
func capGuess(remaining int) int {
	if remaining > 2048 {
		return remaining / 2
	}
	return remaining
}

// SliceFilterInPlace returns elements that pass the predicate function.
// Input slice is modified and must be discarded after calling.
func SliceFilterInPlace[T any](slice []T, predicate func(v T) bool) []T {
	var n int
	for i := range slice {
		if predicate(slice[i]) {
			slice[n] = slice[i]
			n++
		}
	}
	return slice[:n]
}

// SliceFilterInto appends elements that pass the predicate function from the input slices into dest.
func SliceFilterInto[T any](dest []T, predicate func(T) bool, inputs ...[]T) []T {
	for _, input := range inputs {
		for _, v := range input {
			if predicate(v) {
				dest = append(dest, v)
			}
		}
	}
	return dest
}

// SliceSplit partitions elements based on the predicate function.
// Returns (trueElements, falseElements).
func SliceSplit[T any](slice []T, predicate func(v T) bool) ([]T, []T) {
	if len(slice) == 0 {
		return slice, nil
	}

	var splitIndex int
	first := predicate(slice[0])
	for splitIndex = 1; splitIndex < len(slice); splitIndex++ {
		if first != predicate(slice[splitIndex]) {
			break
		}
	}

	// If all are the same, return early
	if splitIndex == len(slice) {
		if first {
			return slice, nil
		} else {
			return nil, slice
		}
	}

	// Allocate slices and copy first segment
	remainingBuff := capGuess(len(slice) - splitIndex)
	var trueList, falseList []T
	if first {
		trueList = append(make([]T, 0, splitIndex+remainingBuff-1), slice[:splitIndex]...)
		falseList = append(make([]T, 0, remainingBuff), slice[splitIndex])
	} else {
		falseList = append(make([]T, 0, splitIndex+remainingBuff-1), slice[:splitIndex]...)
		trueList = append(make([]T, 0, remainingBuff), slice[splitIndex])
	}
	// Finish iterating appending remaining elements
	for i := splitIndex + 1; i < len(slice); i++ {
		if predicate(slice[i]) {
			trueList = append(trueList, slice[i])
		} else {
			falseList = append(falseList, slice[i])
		}
	}

	return trueList, falseList
}

// SliceSplitInPlace partitions elements based on the predicate function.
// Input slice is modified and must be discarded after calling. Resulting slices will remain in the original order.
// If order is not important use SliceSplitInPlaceUnstable for an even faster implementation.
// Returns (trueElements, falseElements).
func SliceSplitInPlace[T any](slice []T, predicate func(v T) bool) ([]T, []T) {
	n := len(slice)
	if n == 0 {
		return slice, nil
	}

	if predicate(slice[0]) { // first element is true
		// Reuse front of slice for TRUEs; allocate FALSE buffer lazily.
		trueList := slice[:0]
		trueList = append(trueList, slice[0]) // first element already known true

		var falseBuf []T // stays nil if we never see a false
		for i := 1; i < n; i++ {
			v := slice[i]
			isTrue := predicate(v) // one evaluation per element
			if isTrue {
				trueList = append(trueList, v) // writes to earlier indices only; safe
			} else {
				if falseBuf == nil {
					// Allocate when we discover the split
					falseBuf = make([]T, 0, capGuess(n-i))
				}
				falseBuf = append(falseBuf, v)
			}
		}
		return trueList, falseBuf
	}

	// Reuse front of slice for FALSEs; allocate TRUE buffer lazily.
	falseList := slice[:0]
	falseList = append(falseList, slice[0]) // first element already known false

	var trueBuf []T // stays nil if we never see a true
	for i := 1; i < n; i++ {
		v := slice[i]
		isTrue := predicate(v) // one evaluation per element
		if isTrue {
			if trueBuf == nil {
				trueBuf = make([]T, 0, capGuess(n-i))
			}
			trueBuf = append(trueBuf, v)
		} else {
			falseList = append(falseList, v) // safe as above
		}
	}
	return trueBuf, falseList
}

// SliceSplitInPlaceUnstable partitions elements based on the predicate function.
// Input slice is modified and must be discarded after calling.
// Resulting slices order may change from the original input.
// Returns (trueElements, falseElements).
func SliceSplitInPlaceUnstable[T any](slice []T, predicate func(v T) bool) ([]T, []T) {
	if len(slice) == 0 {
		return slice, nil
	}

	i, j := 0, len(slice)-1
	for {
		for i <= j && predicate(slice[i]) {
			i++
		}
		for i <= j && !predicate(slice[j]) {
			j--
		}
		if i >= j {
			break
		}
		slice[i], slice[j] = slice[j], slice[i]
		i++
		j--
	}

	switch {
	case i == 0:
		return nil, slice
	case i == len(slice):
		return slice, nil
	default:
		return slice[:i], slice[i:]
	}
}

// SliceTransform converts each element using the conversion function.
func SliceTransform[I any, R any](conversion func(I) R, inputs ...[]I) []R {
	result := make([]R, 0, sliceTotalSize(inputs))
	for _, input := range inputs {
		for _, v := range input {
			result = append(result, conversion(v))
		}
	}
	return result
}

// SliceToSet accepts slices of a comparable type and returns a Map with the entries as the key.
// This allows an easy de-duplicated union between slices, as well as providing a map for fast lookup if values are present.
func SliceToSet[T comparable](slices ...[]T) map[T]struct{} {
	result := make(map[T]struct{}, sliceTotalSize(slices))
	SliceIntoSet(result, slices...)
	return result
}

// SliceIntoSet accepts slices of a comparable type and a Map to set entries into.
func SliceIntoSet[T comparable](m map[T]struct{}, slices ...[]T) {
	for _, slice := range slices {
		for _, value := range slice {
			m[value] = struct{}{}
		}
	}
}

// SliceToSetBy accepts slices of any type with a function to convert those types while storing the result
// as the key to the resulting map. This allows in a single step a combination of SliceTransform with SliceToSet,
func SliceToSetBy[I any, R comparable](keyfunc func(I) R, slices ...[]I) map[R]struct{} {
	result := make(map[R]struct{}, sliceTotalSize(slices))
	SliceIntoSetBy(keyfunc, result, slices...)
	return result
}

// SliceIntoSetBy accepts slices of any type with a function to convert those types while storing the result
// as the key to the resulting map.
func SliceIntoSetBy[I any, R comparable](keyfunc func(I) R, m map[R]struct{}, slices ...[]I) {
	for _, slice := range slices {
		for _, inputVal := range slice {
			m[keyfunc(inputVal)] = struct{}{}
		}
	}
}

// SliceToCounts accepts slices of a comparable type and returns a Map with the entries as the key similar to SliceToSet.
// SliceToCounts will count how many times each entry is witnessed, provided through the returned map values.
func SliceToCounts[T comparable](slices ...[]T) map[T]int {
	result := make(map[T]int, sliceTotalSize(slices))
	SliceIntoCounts(result, slices...)
	return result
}

// SliceIntoCounts accepts slices of a comparable type and a Map with the entries to add the count to.
func SliceIntoCounts[T comparable](m map[T]int, slices ...[]T) {
	for _, slice := range slices {
		for _, value := range slice {
			m[value]++
		}
	}
}

// SliceToCountsBy accepts slices of any type and returns a Map with the generated keys and their counts.
// SliceToCountsBy will count how many times each generated key is witnessed, provided through the returned map values.
func SliceToCountsBy[T any, K comparable](keyfunc func(T) K, slices ...[]T) map[K]int {
	result := make(map[K]int, sliceTotalSize(slices))
	SliceIntoCountsBy(keyfunc, result, slices...)
	return result
}

// SliceIntoCountsBy accepts slices of any type and a Map with the generated keys and their counts which are added to.
func SliceIntoCountsBy[T any, K comparable](keyfunc func(T) K, m map[K]int, slices ...[]T) {
	for _, slice := range slices {
		for _, value := range slice {
			m[keyfunc(value)]++
		}
	}
}

// SliceToIndexBy accepts a function to convert the values to a comparable key, creating an index map.
// Expects each key to be unique. If duplicate keys exist, later values overwrite earlier ones.
func SliceToIndexBy[T any, K comparable](keyfunc func(T) K, slices ...[]T) map[K]T {
	result := make(map[K]T, sliceTotalSize(slices))
	SliceIntoIndexBy(keyfunc, result, slices...)
	return result
}

// SliceIntoIndexBy accepts a function to convert the values to a comparable key, adding to the provided index map.
// Expects each key to be unique. If duplicate keys exist, later values overwrite earlier ones.
func SliceIntoIndexBy[T any, K comparable](keyfunc func(T) K, m map[K]T, slices ...[]T) {
	for _, slice := range slices {
		for _, value := range slice {
			key := keyfunc(value)
			m[key] = value
		}
	}
}

// SliceToGroupsBy accepts a function to convert the values to a comparable key, and groups the values based on the keys.
func SliceToGroupsBy[T any, K comparable](keyfunc func(T) K, slices ...[]T) map[K][]T {
	result := make(map[K][]T, sliceTotalSize(slices))
	SliceIntoGroupsBy(keyfunc, result, slices...)
	return result
}

// SliceIntoGroupsBy accepts a function to convert the values to a comparable key, and groups the values based on the keys.
func SliceIntoGroupsBy[T any, K comparable](keyfunc func(T) K, m map[K][]T, slices ...[]T) {
	for _, slice := range slices {
		for _, value := range slice {
			key := keyfunc(value)
			m[key] = append(m[key], value)
		}
	}
}

func sliceTotalSize[T any](slices [][]T) int {
	var size int
	for _, slice := range slices {
		size += len(slice)
	}
	return capGuess(size)
}
