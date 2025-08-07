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
				return nil // No true elements found
			}

			// Check if all remaining elements are consecutive and true
			consecutiveEnd := firstTrueIndex
			for i := firstTrueIndex + 1; i < len(slice); i++ {
				if predicate(slice[i]) {
					consecutiveEnd = i
					continue
				}

				// Found a false element, check if consecutive section continues
				nonConsecutiveStart := i + 1
				for j := i + 1; j < len(slice); j++ {
					if predicate(slice[j]) {
						// Found another true after false, not consecutive - need to allocate
						remaining := len(slice) - nonConsecutiveStart
						result := make([]T, 0, (consecutiveEnd-firstTrueIndex+1)+capGuess(remaining))
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
			}

			// All elements from firstTrueIndex to end are true
			return slice[firstTrueIndex:]
		} else {
			// Handle prefix case - find consecutive suffix section
			firstTrueIndex := -1

			// Find first true element after falseIndex
			for i := falseIndex + 1; i < len(slice); i++ {
				if predicate(slice[i]) {
					firstTrueIndex = i
					break
				}
			}

			if firstTrueIndex == -1 {
				// No true elements in suffix, return prefix only
				return slice[:falseIndex]
			}

			// Check if suffix elements are consecutive and true
			consecutiveEnd := firstTrueIndex
			for i := firstTrueIndex + 1; i < len(slice); i++ {
				if predicate(slice[i]) {
					consecutiveEnd = i
					continue
				}

				// Found a false element, check if consecutive section continues
				for j := i + 1; j < len(slice); j++ {
					if predicate(slice[j]) {
						// Found another true after false, not consecutive - need to allocate
						result := make([]T, 0, falseIndex+(consecutiveEnd-firstTrueIndex+1)+capGuess(len(slice)-j))
						result = append(result, slice[:falseIndex]...)
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
				// No more true elements found, combine prefix with consecutive suffix view
				result := make([]T, 0, falseIndex+(consecutiveEnd-firstTrueIndex+1))
				result = append(result, slice[:falseIndex]...)
				result = append(result, slice[firstTrueIndex:consecutiveEnd+1]...)
				return result
			}

			// All suffix elements from firstTrueIndex to end are true
			result := make([]T, 0, falseIndex+(len(slice)-firstTrueIndex))
			result = append(result, slice[:falseIndex]...)
			result = append(result, slice[firstTrueIndex:]...)
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

// SliceTransform converts each element using the conversion function.
func SliceTransform[I any, R any](input []I, conversion func(I) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = conversion(v)
	}
	return result
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
