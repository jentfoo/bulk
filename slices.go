package bulk

// SliceFilter returns elements that pass the predicate function.
// May return the original slice if all elements pass (no allocation).
func SliceFilter[T any](predicate func(v T) bool, slices ...[]T) []T {
	switch len(slices) {
	case 0:
		return nil
	case 1:
		result, _ := singleSliceFilter(predicate, slices[0])
		return result
	}

	results := make([][]T, 0, len(slices))
	concatInPlace := true
	for i, slice := range slices {
		partResult, view := singleSliceFilter(predicate, slice)
		if len(results) == 1 && len(results[0]) == 0 {
			// Our head slice is empty, check if we should prefer this slice instead
			if len(partResult) > 0 ||
				// if no results this will be a view, but maybe we can retain a larger view for the user
				cap(partResult) > cap(results[0]) {
				results[0] = partResult
				concatInPlace = !view
				continue
			}
		}
		// head not replaced, check if we should append
		if i == 0 || len(partResult) > 0 {
			concatInPlace = concatInPlace && !view // if view is used we have to copy in concat
			results = append(results, partResult)
		}
	}
	return sliceConcat(results, concatInPlace)
}

func singleSliceFilter[T any](predicate func(v T) bool, slice []T) ([]T, bool) {
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
				return slice[:0], true // No true elements found
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
				return slice[firstTrueIndex:], true // All elements from firstTrueIndex to end are true
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
					return SliceFilterInto(result, predicate, slice[j+1:]), false
				}
			}
			// No more true elements found, return consecutive view
			return slice[firstTrueIndex : consecutiveEnd+1], true
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
				return slice[:falseIndex], true // No true elements in suffix, return prefix only
			}

			// true+ -> false+ -> true - We must allocate at this point
			result := make([]T, 0, falseIndex+capGuess(len(slice)-secondTrueIndex))
			result = append(result, slice[:falseIndex]...)
			result = append(result, slice[secondTrueIndex])
			return SliceFilterInto(result, predicate, slice[secondTrueIndex+1:]), false
		}
	}
	return slice, true // all records tested to true
}

// sliceConcat is similar to slices.Concat (which is why API is not elevated),
// but differs in if provided a single slice it's returned without a copy.
func sliceConcat[T any](slices [][]T, inPlace bool) []T {
	switch len(slices) {
	case 0:
		return make([]T, 0)
	case 1:
		return slices[0]
	case 2:
		if len(slices[0]) == 0 { // empty first slice, easy result choice
			if len(slices[1]) > 0 || cap(slices[1]) > cap(slices[0]) {
				return slices[1]
			}
			return slices[0] // both empty, 0 is largest capacity
		}
	}

	if totalSize := sliceTotalSize(slices); inPlace && totalSize <= cap(slices[0]) {
		result := slices[0]
		for i := 1; i < len(slices); i++ { // skip the first slice (set directly)
			result = append(result, slices[i]...)
		}
		return result
	} else { // allocate and copy
		result := make([]T, 0, totalSize)
		for _, slice := range slices {
			result = append(result, slice...)
		}
		return result
	}
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

// SliceFilterInPlace returns elements that pass the predicate function.
// Input slice is modified and must be discarded after calling.
func SliceFilterInPlace[T any](predicate func(v T) bool, slices ...[]T) []T {
	switch len(slices) {
	case 0:
		return nil
	case 1:
		return singleSliceFilterInPlace(predicate, slices[0])
	}

	results := make([][]T, 0, len(slices))
	firstCapacity := len(slices[0]) // initialize with the raw initial length to consider as capacity (NOT full cap range)
	for i, slice := range slices {
		partResult := singleSliceFilterInPlace(predicate, slice)
		if len(results) == 1 && len(results[0]) == 0 {
			// Our head slice is empty, check if we should prefer this slice instead
			if i == len(slices)-1 { // last record, pick the ideal result and break
				if len(partResult) > 0 || cap(partResult) > cap(results[0]) {
					results[0] = partResult
				}
				break
			}
			// if the current result offers more capacity retain it instead
			// otherwise just fall through and both will be retained
			// if the current is the only added slice sliceConcat will optimize away the empty head slice
			// if more slices are added (3+ total), then sliceConcat may use this slice to avoid allocations
			currCapacity := len(slice) - len(partResult)
			if firstCapacity < currCapacity {
				firstCapacity = currCapacity
				results[0] = partResult
				continue
			}
		} // if not continue above, fall to check below
		if i == 0 || len(partResult) > 0 {
			firstCapacity -= len(partResult)
			results = append(results, partResult)
		}
	}
	return sliceConcat(results, firstCapacity >= 0)
}

func singleSliceFilterInPlace[T any](predicate func(v T) bool, slice []T) []T {
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
func SliceSplit[T any](predicate func(v T) bool, slices ...[]T) ([]T, []T) {
	switch len(slices) {
	case 0:
		return nil, nil
	case 1:
		tSlice, fSlice, _, _ := singleSliceSplit(predicate, slices[0])
		return tSlice, fSlice
	}

	trueResults, falseResults := make([][]T, 0, len(slices)), make([][]T, 0, len(slices))
	trueConcatInPlace, falseConcatInPlace := true, true
	for i, slice := range slices {
		tSlice, fSlice, tView, fView := singleSliceSplit(predicate, slice)
		if len(trueResults) == 1 && len(trueResults[0]) == 0 && // check for empty head to replace
			(len(tSlice) > 0 || // replace with actual results
				(!tView && !trueConcatInPlace) || // can be used to upgrade into an in place copy
				(!trueConcatInPlace && cap(tSlice) > cap(trueResults[0]))) { // wont downgrade and has more capacity
			trueResults[0] = tSlice
			trueConcatInPlace = !tView
		} else if len(tSlice) > 0 || i == 0 /* ensure at least one result */ {
			trueConcatInPlace = trueConcatInPlace && !tView
			trueResults = append(trueResults, tSlice)
		}
		if len(falseResults) == 1 && len(falseResults[0]) == 0 && // check for empty head to replace
			(len(fSlice) > 0 || // replace with actual results
				(!fView && !falseConcatInPlace) || // can be used to upgrade into an in place copy
				(!falseConcatInPlace && cap(fSlice) > cap(falseResults[0]))) { // wont downgrade and has more capacity
			falseResults[0] = fSlice
			falseConcatInPlace = !fView
		} else if len(fSlice) > 0 || i == 0 /* ensure at least one result */ {
			falseConcatInPlace = falseConcatInPlace && !fView
			falseResults = append(falseResults, fSlice)
		}
	}
	trueResult := sliceConcat(trueResults, trueConcatInPlace)
	falseResult := sliceConcat(falseResults, falseConcatInPlace)
	return trueResult, falseResult
}

// SliceSplit partitions elements based on the predicate function.
// Returns (trueElements, falseElements).
func singleSliceSplit[T any](predicate func(v T) bool, slice []T) ([]T, []T, bool, bool) {
	if len(slice) == 0 {
		return slice, nil, true, false
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
			return slice, nil, true, false
		} else {
			return nil, slice, false, true
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

	return trueList, falseList, false, false
}

// SliceSplitInPlace partitions elements based on the predicate function.
// Input slice is modified and must be discarded after calling. Resulting slices will remain in the original order.
// If order is not important use SliceSplitInPlaceUnstable for an even faster implementation.
// Returns (trueElements, falseElements).
func SliceSplitInPlace[T any](predicate func(v T) bool, slice []T) ([]T, []T) {
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
				if falseBuf == nil { // Allocate when we discover the split
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
func SliceSplitInPlaceUnstable[T any](predicate func(v T) bool, slice []T) ([]T, []T) {
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
	SliceIntoSetBy(result, keyfunc, slices...)
	return result
}

// SliceIntoSetBy accepts slices of any type with a function to convert those types while storing the result
// as the key to the resulting map.
func SliceIntoSetBy[I any, R comparable](m map[R]struct{}, keyfunc func(I) R, slices ...[]I) {
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
	SliceIntoCountsBy(result, keyfunc, slices...)
	return result
}

// SliceIntoCountsBy accepts slices of any type and a Map with the generated keys and their counts which are added to.
func SliceIntoCountsBy[T any, K comparable](m map[K]int, keyfunc func(T) K, slices ...[]T) {
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
	SliceIntoIndexBy(result, keyfunc, slices...)
	return result
}

// SliceIntoIndexBy accepts a function to convert the values to a comparable key, adding to the provided index map.
// Expects each key to be unique. If duplicate keys exist, later values overwrite earlier ones.
func SliceIntoIndexBy[T any, K comparable](m map[K]T, keyfunc func(T) K, slices ...[]T) {
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
	SliceIntoGroupsBy(result, keyfunc, slices...)
	return result
}

// SliceIntoGroupsBy accepts a function to convert the values to a comparable key, and groups the values based on the keys.
func SliceIntoGroupsBy[T any, K comparable](m map[K][]T, keyfunc func(T) K, slices ...[]T) {
	for _, slice := range slices {
		for _, value := range slice {
			key := keyfunc(value)
			m[key] = append(m[key], value)
		}
	}
}

// SliceIntersect returns elements that exist in both slices, preserving order from slice a.
func SliceIntersect[T comparable](a, b []T) []T {
	if len(a) == 0 {
		return a
	} else if len(b) == 0 {
		return b
	}

	maxCount := len(a)
	if len(b) < maxCount {
		maxCount = len(b)
	}

	// Collect intersection, preserving order from slice a
	bLookup := SliceToSet(b)
	var result []T
	var seen map[T]struct{}
	for aIdx, v := range a {
		if _, exists := bLookup[v]; exists {
			if _, duplicate := seen[v]; !duplicate {
				if result == nil { // allocate based of potential remaining
					if aMax := len(a) - aIdx; aMax < maxCount {
						maxCount = aMax // conditional because b may still have been the min
					}
					seen = make(map[T]struct{}, maxCount)
					result = make([]T, 0, capGuess(maxCount))
				}
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result
}

// SliceDifference returns elements that exist in slice a, but not in slice b, preserving order from slice a.
func SliceDifference[T comparable](a, b []T) []T {
	if len(a) == 0 {
		return a
	}

	// Collect elements from a that are not in b, with deduplication
	exclude := SliceToSet(b)
	var result []T
	var seen map[T]struct{}
	for aIdx, v := range a {
		if _, exists := exclude[v]; !exists {
			if _, duplicate := seen[v]; !duplicate {
				if result == nil { // allocate based of potential remaining
					seen = make(map[T]struct{}, len(a)-aIdx)
					result = make([]T, 0, capGuess(len(a)-aIdx))
				}
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result
}

// capGuess attempts to guess the sizing for an allocation, large allocations are reduced.
func capGuess(remaining int) int {
	if remaining > 2048 {
		return remaining / 2
	}
	return remaining
}

func sliceTotalSize[T any](slices [][]T) int {
	var size int
	for _, slice := range slices {
		size += len(slice)
	}
	return capGuess(size)
}
