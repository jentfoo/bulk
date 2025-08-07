package bulk

// SliceFilter returns elements that pass the predicate function.
// May return the original slice if all elements pass (no allocation).
func SliceFilter[T any](slice []T, predicate func(v T) bool) []T {
	for falseIndex, v := range slice {
		if !predicate(v) {
			if falseIndex == 0 {
				// iterate until a true result is found, then start appending at that point
				var result []T
				for i := falseIndex + 1; i < len(slice); i++ {
					if predicate(slice[i]) {
						if result == nil {
							result = make([]T, 0, capGuess(len(slice)-i))
						}
						result = append(result, slice[i])
					}
				}
				return result
			} else {
				// copy all records that already passed, and then finish iteration to produce result
				result := append(make([]T, 0, falseIndex+capGuess(len(slice)-falseIndex-1)), slice[:falseIndex]...)
				for i := falseIndex + 1; i < len(slice); i++ {
					if predicate(slice[i]) {
						result = append(result, slice[i])
					}
				}
				return result
			}
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

// SliceSplit partitions elements based on the predicate function.
// Returns (trueElements, falseElements).
func SliceSplit[T any](slice []T, predicate func(v T) bool) ([]T, []T) {
	if len(slice) == 0 {
		return nil, nil
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
		return nil, nil
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
					// Allocate when we discover the split. Use a conservative guess.
					rem := n - i
					falseBuf = make([]T, 0, capGuess(rem))
				}
				falseBuf = append(falseBuf, v)
			}
		}
		if len(falseBuf) == 0 {
			return trueList, nil
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
				rem := n - i
				trueBuf = make([]T, 0, capGuess(rem))
			}
			trueBuf = append(trueBuf, v)
		} else {
			falseList = append(falseList, v) // safe as above
		}
	}
	if len(trueBuf) == 0 {
		return nil, falseList
	}
	return trueBuf, falseList
}

// SliceSplitInPlaceUnstable partitions elements based on the predicate function.
// Input slice is modified and must be discarded after calling.
// Resulting slices order may change from the original input.
// Returns (trueElements, falseElements).
func SliceSplitInPlaceUnstable[T any](slice []T, predicate func(v T) bool) ([]T, []T) {
	if len(slice) == 0 {
		return nil, nil
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

// SliceReverseInPlace reverses elements in-place.
func SliceReverseInPlace[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// SliceConversion transforms each element using the conversion function.
func SliceConversion[I any, R any](input []I, conversion func(I) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = conversion(v)
	}
	return result
}

// SliceUnion returns all values combined into a single slice.
func SliceUnion[T comparable](slices ...[]T) []T {
	switch len(slices) {
	case 0:
		return nil
	case 1:
		return slices[0]
	}

	var size int
	for _, slice := range slices {
		size += len(slice)
	}
	result := make([]T, 0, size)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// SliceUnionUnique returns all unique values from the provided slices.
func SliceUnionUnique[T comparable](slices ...[]T) []T {
	var uniqMap map[T]bool
	switch len(slices) {
	case 0:
		return nil
	case 1:
		uniqMap = make(map[T]bool, len(slices[0]))
	default:
		uniqMap = make(map[T]bool)
	}

	for _, s := range slices {
		for _, v := range s {
			uniqMap[v] = true
		}
	}
	return mapKeys(uniqMap)
}

// SliceRemoveAt removes the element at the specified index.
// Returns original slice if index is out of bounds.
func SliceRemoveAt[T any](slice []T, index int) []T {
	switch len(slice) {
	case 0:
		return nil
	case 1:
		if index == 0 {
			return nil
		} else {
			return slice
		}
	}
	if index >= len(slice) || index < 0 {
		return slice // index out of range
	} else if index == 0 {
		return slice[1:]
	} else if index == len(slice)-1 {
		return slice[:index]
	}

	result := make([]T, index, len(slice)-1)
	// Copy the first half of the slice up to index
	copy(result, slice[:index])
	// Copy the second half of the slice after index
	return append(result, slice[index+1:]...)
}

// SliceRemoveAtInPlace removes the element at the specified index.
// Input slice is modified and must be discarded after calling.
func SliceRemoveAtInPlace[T any](slice []T, index int) []T {
	switch len(slice) {
	case 0:
		return nil
	case 1:
		if index == 0 {
			return nil
		} else {
			return slice
		}
	}
	if index >= len(slice) || index < 0 {
		return slice // index out of range
	} else if index == 0 {
		return slice[1:]
	} else if index == len(slice)-1 {
		return slice[:index]
	}

	// Shift the smaller portion to minimize data movement
	if index < len(slice)/2 { // Shift up
		copy(slice[1:index+1], slice[0:index])
		return slice[1:] // Remove first element
	} else { // Shift down
		copy(slice[index:], slice[index+1:])
		return slice[:len(slice)-1] // Remove last element
	}
}
