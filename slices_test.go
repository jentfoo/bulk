package bulk

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sliceLargeInput = []int{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
	27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
}
var sliceTestCases = []struct {
	name        string
	input       []int
	testFunc    func(v int) bool
	expectTrue  []int
	expectFalse []int
}{
	{
		name:        "nil",
		input:       nil,
		testFunc:    func(v int) bool { return v > 0 },
		expectTrue:  nil,
		expectFalse: nil,
	},
	{
		name:        "empty",
		input:       []int{},
		testFunc:    func(v int) bool { return v > 0 },
		expectTrue:  nil,
		expectFalse: nil,
	},
	{
		name:        "all_true",
		input:       []int{2, 4, 6},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6},
		expectFalse: nil,
	},
	{
		name:        "all_false",
		input:       []int{1, 3, 5},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  nil,
		expectFalse: []int{1, 3, 5},
	},
	{
		name:        "all_true_large",
		input:       sliceLargeInput,
		testFunc:    func(v int) bool { return true },
		expectTrue:  sliceLargeInput,
		expectFalse: nil,
	},
	{
		name:        "all_false_large",
		input:       sliceLargeInput,
		testFunc:    func(v int) bool { return false },
		expectTrue:  nil,
		expectFalse: sliceLargeInput,
	},
	{
		name:        "single_true",
		input:       []int{2},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2},
		expectFalse: nil,
	},
	{
		name:        "single_false",
		input:       []int{1},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  nil,
		expectFalse: []int{1},
	},
	{
		name:        "one_true",
		input:       []int{1, 2, 3, 4, 6, 8},
		testFunc:    func(v int) bool { return v == 1 },
		expectTrue:  []int{1},
		expectFalse: []int{2, 3, 4, 6, 8},
	},
	{
		name:        "true_end",
		input:       []int{2, 3, 4, 6, 8, 1, 1},
		testFunc:    func(v int) bool { return v == 1 },
		expectTrue:  []int{1, 1},
		expectFalse: []int{2, 3, 4, 6, 8},
	},
	{
		name:        "mixed_split_first",
		input:       []int{1, 2, 1, 4, 6, 8},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6, 8},
		expectFalse: []int{1, 1},
	},
	{
		name:        "mixed_split_second",
		input:       []int{2, 1, 4, 6, 8},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6, 8},
		expectFalse: []int{1},
	},
	{
		name:        "mixed_split_middle",
		input:       []int{2, 4, 6, 1, 3, 5, 8},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6, 8},
		expectFalse: []int{1, 3, 5},
	},
	{
		name:        "mixed_split_last",
		input:       []int{2, 4, 6, 1},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6},
		expectFalse: []int{1},
	},
	{
		name:     "mixed_split_large",
		input:    sliceLargeInput,
		testFunc: func(v int) bool { return v%2 == 0 },
		expectTrue: []int{
			2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26,
			28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50,
		},
		expectFalse: []int{
			1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25,
			27, 29, 31, 33, 35, 37, 39, 41, 43, 45, 47, 49,
		},
	},
	// Test cases for consecutive chunk optimizations
	{
		name:        "consecutive_start_after_one_false",
		input:       []int{1, 2, 3, 4, 5},
		testFunc:    func(v int) bool { return v >= 2 },
		expectTrue:  []int{2, 3, 4, 5},
		expectFalse: []int{1},
	},
	{
		name:        "consecutive_start_after_multiple_false",
		input:       []int{1, 3, 5, 6, 7, 8},
		testFunc:    func(v int) bool { return v >= 6 },
		expectTrue:  []int{6, 7, 8},
		expectFalse: []int{1, 3, 5},
	},
	{
		name:        "consecutive_middle_chunk_only",
		input:       []int{1, 3, 6, 7, 8, 9, 5, 11},
		testFunc:    func(v int) bool { return v >= 6 && v <= 9 },
		expectTrue:  []int{6, 7, 8, 9},
		expectFalse: []int{1, 3, 5, 11},
	},
	{
		name:        "consecutive_end_chunk_only",
		input:       []int{1, 3, 5, 8, 9, 10},
		testFunc:    func(v int) bool { return v >= 8 },
		expectTrue:  []int{8, 9, 10},
		expectFalse: []int{1, 3, 5},
	},
	{
		name:        "prefix_plus_consecutive_suffix",
		input:       []int{2, 4, 1, 3, 8, 10, 12},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 8, 10, 12},
		expectFalse: []int{1, 3},
	},
	{
		name:        "prefix_plus_non_consecutive_suffix",
		input:       []int{2, 4, 1, 8, 9, 10},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 8, 10},
		expectFalse: []int{1, 9},
	},
	{
		name:        "consecutive_chunk_then_gap_then_more",
		input:       []int{1, 6, 7, 8, 3, 5, 12, 14},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{6, 8, 12, 14},
		expectFalse: []int{1, 7, 3, 5},
	},
	{
		name:        "single_true_after_false",
		input:       []int{1, 3, 5, 8},
		testFunc:    func(v int) bool { return v == 8 },
		expectTrue:  []int{8},
		expectFalse: []int{1, 3, 5},
	},
	{
		name:        "all_false_after_initial_false",
		input:       []int{1, 3, 5, 7, 9},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  nil,
		expectFalse: []int{1, 3, 5, 7, 9},
	},
	{
		name:        "consecutive_at_very_end",
		input:       []int{1, 3, 5, 7, 8, 10, 12},
		testFunc:    func(v int) bool { return v >= 8 },
		expectTrue:  []int{8, 10, 12},
		expectFalse: []int{1, 3, 5, 7},
	},
	{
		name:        "prefix_only_no_suffix_true",
		input:       []int{2, 4, 6, 1, 3, 5},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6},
		expectFalse: []int{1, 3, 5},
	},
	{
		name:        "long_consecutive_in_middle",
		input:       []int{1, 3, 10, 11, 12, 13, 14, 15, 5, 7},
		testFunc:    func(v int) bool { return v >= 10 && v <= 15 },
		expectTrue:  []int{10, 11, 12, 13, 14, 15},
		expectFalse: []int{1, 3, 5, 7},
	},
	{
		name:        "multiple_gaps_non_consecutive",
		input:       []int{1, 2, 4, 5, 6, 8, 9, 10},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6, 8, 10},
		expectFalse: []int{1, 5, 9},
	},
	{
		name:        "consecutive_with_single_gap",
		input:       []int{1, 6, 7, 8, 9, 5, 12},
		testFunc:    func(v int) bool { return v >= 6 },
		expectTrue:  []int{6, 7, 8, 9, 12},
		expectFalse: []int{1, 5},
	},
	{
		name:        "large_consecutive_chunk_optimization",
		input:       []int{1, 3, 5, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		testFunc:    func(v int) bool { return v >= 10 },
		expectTrue:  []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		expectFalse: []int{1, 3, 5},
	},
}

func sliceDup[T any](s []T) []T {
	result := make([]T, len(s))
	copy(result, s)
	return result
}

func TestSliceFilter(t *testing.T) {
	t.Parallel()

	for i, tt := range sliceTestCases {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			slice := SliceFilter(tt.input, tt.testFunc)
			if len(tt.expectTrue) == 0 {
				assert.Empty(t, slice)
			} else {
				assert.Equal(t, tt.expectTrue, slice)
			}
		})
	}
}

func TestSliceFilterInPlace(t *testing.T) {
	t.Parallel()

	for i, tt := range sliceTestCases {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			// make a copy of the test input to avoid changing it
			slice := SliceFilterInPlace(sliceDup(tt.input), tt.testFunc)
			if len(tt.expectTrue) == 0 {
				assert.Empty(t, slice)
			} else {
				assert.Equal(t, tt.expectTrue, slice)
			}
		})
	}
}

func TestSliceSplit(t *testing.T) {
	t.Parallel()

	for i, tt := range sliceTestCases {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			trueSlice, falseSlice := SliceSplit(tt.input, tt.testFunc)
			assert.Equal(t, tt.expectTrue, trueSlice)
			assert.Equal(t, tt.expectFalse, falseSlice)
		})
	}
}

func TestSliceSplitInPlace(t *testing.T) {
	t.Parallel()

	for i, tt := range sliceTestCases {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			// make a copy of the test input to avoid changing it
			trueSlice, falseSlice := SliceSplitInPlace(sliceDup(tt.input), tt.testFunc)
			assert.Equal(t, tt.expectTrue, trueSlice)
			assert.Equal(t, tt.expectFalse, falseSlice)
		})
	}
}

func TestSliceSplitInPlaceUnstable(t *testing.T) {
	t.Parallel()

	for i, tt := range sliceTestCases {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			// make a copy of the test input to avoid changing it
			trueSlice, falseSlice := SliceSplitInPlaceUnstable(sliceDup(tt.input), tt.testFunc)
			assert.ElementsMatch(t, tt.expectTrue, trueSlice)
			assert.ElementsMatch(t, tt.expectFalse, falseSlice)
		})
	}
}

func TestSliceTransform(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		var input []int
		result := SliceTransform(input, func(i int) int { return i })
		assert.Empty(t, result)
	})

	t.Run("empty", func(t *testing.T) {
		result := SliceTransform([]int{}, func(i int) int { return i })
		assert.Empty(t, result)
	})

	t.Run("int_to_string", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := SliceTransform(input, func(i int) string { return strconv.Itoa(i) })
		assert.Equal(t, []string{"1", "2", "3"}, result)
	})

	t.Run("float_to_int", func(t *testing.T) {
		input := []float64{1.1, 2.2}
		result := SliceTransform(input, func(f float64) int { return int(f) })
		assert.Equal(t, []int{1, 2}, result)
	})
}

var sliceRemoveTestCases = []struct {
	name   string
	input  []int
	index  int
	expect []int
}{
	{
		name:   "nil",
		input:  nil,
		index:  0,
		expect: nil,
	},
	{
		name:   "empty",
		input:  []int{},
		index:  0,
		expect: nil,
	},
	{
		name:   "single",
		input:  []int{42},
		index:  0,
		expect: nil,
	},
	{
		name:   "first",
		input:  []int{1, 2, 3, 4, 5},
		index:  0,
		expect: []int{2, 3, 4, 5},
	},
	{
		name:   "middle",
		input:  []int{1, 2, 3, 4, 5},
		index:  2,
		expect: []int{1, 2, 4, 5},
	},
	{
		name:   "last",
		input:  []int{1, 2, 3, 4, 5},
		index:  4,
		expect: []int{1, 2, 3, 4},
	},
	{
		name:   "negative_index",
		input:  []int{1, 2, 3},
		index:  -1,
		expect: []int{1, 2, 3},
	},
	{
		name:   "index_out_of_bounds",
		input:  []int{1, 2, 3},
		index:  5,
		expect: []int{1, 2, 3},
	},
	{
		name:  "large_first",
		input: sliceLargeInput,
		index: 0,
		expect: []int{
			2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
			27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
		},
	},
	{
		name:  "large_middle1",
		input: sliceLargeInput,
		index: 10,
		expect: []int{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
			27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
		},
	},
	{
		name:  "large_middle2",
		input: sliceLargeInput,
		index: 40,
		expect: []int{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
			27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 42, 43, 44, 45, 46, 47, 48, 49, 50,
		},
	},
	{
		name:  "large_last",
		input: sliceLargeInput,
		index: 49,
		expect: []int{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
			27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
		},
	},
}

func TestSliceRemoveAt(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceRemoveTestCases {
		t.Run(tt.name, func(t *testing.T) {
			result := SliceRemoveAt(tt.input, tt.index)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestSliceRemoveAtInPlace(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceRemoveTestCases {
		t.Run(tt.name, func(t *testing.T) {
			result := SliceRemoveAtInPlace(sliceDup(tt.input), tt.index)
			assert.Equal(t, tt.expect, result)
		})
	}
}
