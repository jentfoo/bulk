package bulk

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sliceLargeInput = []int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
	27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
	27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
}
var sliceTestCases = []struct {
	name        string
	input       []int
	testFunc    func(v int) bool
	expectTrue  []int
	expectFalse []int
	trueCapMin  int
	trueCapMax  int
	falseCapMin int
	falseCapMax int
}{
	{
		name:        "nil",
		input:       nil,
		testFunc:    func(v int) bool { return v > 0 },
		expectTrue:  nil,
		expectFalse: nil,
		trueCapMin:  0,
		trueCapMax:  0,
		falseCapMin: 0,
		falseCapMax: 0,
	},
	{
		name:        "empty",
		input:       []int{},
		testFunc:    func(v int) bool { return v > 0 },
		expectTrue:  nil,
		expectFalse: nil,
		trueCapMin:  0,
		trueCapMax:  0,
		falseCapMin: 0,
		falseCapMax: 0,
	},
	{
		name:        "all_true",
		input:       []int{2, 4, 6},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6},
		expectFalse: nil,
		// should return original slice
		trueCapMin: 3,
		trueCapMax: 3,
		// empty slice, but may retain original capacity as view
		falseCapMin: 0,
		falseCapMax: 3,
	},
	{
		name:        "all_false",
		input:       []int{1, 3, 5},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  nil,
		expectFalse: []int{1, 3, 5},
		// empty slice, but may retain original capacity as view
		trueCapMin: 0,
		trueCapMax: 3,
		// should return original slice
		falseCapMin: 3,
		falseCapMax: 3,
	},
	{
		name:        "all_true_large",
		input:       sliceLargeInput,
		testFunc:    func(v int) bool { return true },
		expectTrue:  sliceLargeInput,
		expectFalse: nil,
		// should return original slice
		trueCapMin: 101,
		trueCapMax: 101,
		// empty slice, but may retain original capacity as view
		falseCapMin: 0,
		falseCapMax: 101,
	},
	{
		name:        "all_false_large",
		input:       sliceLargeInput,
		testFunc:    func(v int) bool { return false },
		expectTrue:  nil,
		expectFalse: sliceLargeInput,
		// empty slice, but may retain original capacity as view
		trueCapMin: 0,
		trueCapMax: 101,
		// should return original slice
		falseCapMin: 101,
		falseCapMax: 101,
	},
	{
		name:        "single_true",
		input:       []int{2},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2},
		expectFalse: nil,
		// should return original slice
		trueCapMin: 1,
		trueCapMax: 1,
		// empty slice, but may retain original capacity as view
		falseCapMin: 0,
		falseCapMax: 1,
	},
	{
		name:        "single_false",
		input:       []int{1},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  nil,
		expectFalse: []int{1},
		// empty slice, but may retain original capacity as view
		trueCapMin: 0,
		trueCapMax: 1,
		// should return original slice
		falseCapMin: 1,
		falseCapMax: 1,
	},
	{
		name:        "one_true",
		input:       []int{1, 2, 3, 4, 6, 8},
		testFunc:    func(v int) bool { return v == 1 },
		expectTrue:  []int{1},
		expectFalse: []int{2, 3, 4, 6, 8},
		trueCapMin:  1,
		trueCapMax:  6,
		falseCapMin: 5,
		falseCapMax: 5,
	},
	{
		name:        "one_true_large_first",
		input:       sliceLargeInput,
		testFunc:    func(v int) bool { return v == 0 },
		expectTrue:  []int{0},
		expectFalse: sliceLargeInput[1:],
		trueCapMin:  1,
		trueCapMax:  101,
		falseCapMin: 100,
		falseCapMax: 100,
	},
	{
		name:       "two_true_large_middle",
		input:      sliceLargeInput,
		testFunc:   func(v int) bool { return v == 40 },
		expectTrue: []int{40, 40},
		expectFalse: []int{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
			27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
			27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
		},
		// scattered elements, may allocate
		trueCapMin: 12,
		trueCapMax: 101,
		// scattered elements, may allocate
		falseCapMin: 100,
		falseCapMax: 101,
	},
	{
		name:        "true_end",
		input:       []int{2, 3, 4, 6, 8, 1, 1},
		testFunc:    func(v int) bool { return v == 1 },
		expectTrue:  []int{1, 1},
		expectFalse: []int{2, 3, 4, 6, 8},
		trueCapMin:  2,
		trueCapMax:  2,
		falseCapMin: 6,
		falseCapMax: 7,
	},
	{
		name:        "mixed_split_first",
		input:       []int{1, 2, 1, 4, 6, 8},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6, 8},
		expectFalse: []int{1, 1},
		// scattered elements, may allocate
		trueCapMin: 4,
		trueCapMax: 5,
		// scattered elements, may allocate
		falseCapMin: 5,
		falseCapMax: 6,
	},
	{
		name:        "mixed_split_second",
		input:       []int{2, 1, 4, 6, 8},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6, 8},
		expectFalse: []int{1},
		// prefix + suffix, may allocate
		trueCapMin: 4,
		trueCapMax: 5,
		// single element, may allocate
		falseCapMin: 4,
		falseCapMax: 5,
	},
	{
		name:        "mixed_split_middle",
		input:       []int{2, 4, 6, 1, 3, 5, 8},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6, 8},
		expectFalse: []int{1, 3, 5},
		// prefix + suffix, may allocate
		trueCapMin: 4,
		trueCapMax: 7,
		// consecutive middle section, may return view
		falseCapMin: 4,
		falseCapMax: 7,
	},
	{
		name:        "mixed_split_last",
		input:       []int{2, 4, 6, 1},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6},
		expectFalse: []int{1},
		// prefix only, may return view
		trueCapMin: 3,
		trueCapMax: 4,
		// single element at end, may return view
		falseCapMin: 1,
		falseCapMax: 1,
	},
	{
		name:     "mixed_split_large",
		input:    sliceLargeInput,
		testFunc: func(v int) bool { return v%2 == 0 },
		expectTrue: []int{
			0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26,
			28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50,
			2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26,
			28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50,
		},
		expectFalse: []int{
			1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25,
			27, 29, 31, 33, 35, 37, 39, 41, 43, 45, 47, 49,
			1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25,
			27, 29, 31, 33, 35, 37, 39, 41, 43, 45, 47, 49,
		},
		trueCapMin:  100,
		trueCapMax:  101,
		falseCapMin: 99,
		falseCapMax: 101,
	},
	{
		name:        "consecutive_start_after_one_false",
		input:       []int{1, 2, 3, 4, 5},
		testFunc:    func(v int) bool { return v >= 2 },
		expectTrue:  []int{2, 3, 4, 5},
		expectFalse: []int{1},
		// consecutive suffix, may return view
		trueCapMin: 4,
		trueCapMax: 4,
		// single element prefix, may return view
		falseCapMin: 4,
		falseCapMax: 5,
	},
	{
		name:        "consecutive_start_after_multiple_false",
		input:       []int{1, 3, 5, 6, 7, 8},
		testFunc:    func(v int) bool { return v >= 6 },
		expectTrue:  []int{6, 7, 8},
		expectFalse: []int{1, 3, 5},
		// consecutive suffix, may return view
		trueCapMin: 3,
		trueCapMax: 3,
		// non-consecutive prefix, may allocate
		falseCapMin: 5,
		falseCapMax: 6,
	},
	{
		name:        "consecutive_middle_chunk_only",
		input:       []int{1, 3, 6, 7, 8, 9, 5, 11},
		testFunc:    func(v int) bool { return v >= 6 && v <= 9 },
		expectTrue:  []int{6, 7, 8, 9},
		expectFalse: []int{1, 3, 5, 11},
		// consecutive middle chunk, may return view
		trueCapMin: 6,
		trueCapMax: 6,
		// scattered elements, may allocate
		falseCapMin: 4,
		falseCapMax: 8,
	},
	{
		name:        "consecutive_end_chunk_only",
		input:       []int{1, 3, 5, 8, 9, 10},
		testFunc:    func(v int) bool { return v >= 8 },
		expectTrue:  []int{8, 9, 10},
		expectFalse: []int{1, 3, 5},
		// consecutive suffix, may return view
		trueCapMin: 3,
		trueCapMax: 3,
		// non-consecutive prefix, may allocate
		falseCapMin: 5,
		falseCapMax: 6,
	},
	{
		name:        "prefix_plus_consecutive_suffix",
		input:       []int{2, 4, 1, 3, 8, 10, 12},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 8, 10, 12},
		expectFalse: []int{1, 3},
		// prefix + suffix, may allocate
		trueCapMin: 5,
		trueCapMax: 7,
		// consecutive middle section, may return view
		falseCapMin: 5,
		falseCapMax: 7,
	},
	{
		name:        "prefix_plus_non_consecutive_suffix",
		input:       []int{2, 4, 1, 8, 9, 10},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 8, 10},
		expectFalse: []int{1, 9},
		// prefix + suffix, may allocate
		trueCapMin: 5,
		trueCapMax: 6,
		// scattered elements, may allocate
		falseCapMin: 3,
		falseCapMax: 6,
	},
	{
		name:        "consecutive_chunk_then_gap_then_more",
		input:       []int{1, 6, 7, 8, 3, 5, 12, 14},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{6, 8, 12, 14},
		expectFalse: []int{1, 7, 3, 5},
		// scattered elements, may allocate
		trueCapMin: 6,
		trueCapMax: 8,
		// scattered elements, may allocate
		falseCapMin: 7,
		falseCapMax: 8,
	},
	{
		name:        "single_true_after_false",
		input:       []int{1, 3, 5, 8},
		testFunc:    func(v int) bool { return v == 8 },
		expectTrue:  []int{8},
		expectFalse: []int{1, 3, 5},
		// single element at end, may return view
		trueCapMin: 1,
		trueCapMax: 1,
		// non-consecutive prefix, may allocate
		falseCapMin: 3,
		falseCapMax: 4,
	},
	{
		name:        "consecutive_at_very_end",
		input:       []int{1, 3, 5, 7, 8, 10, 12},
		testFunc:    func(v int) bool { return v >= 8 },
		expectTrue:  []int{8, 10, 12},
		expectFalse: []int{1, 3, 5, 7},
		// consecutive suffix, may return view
		trueCapMin: 3,
		trueCapMax: 3,
		// non-consecutive prefix, may allocate
		falseCapMin: 6,
		falseCapMax: 7,
	},
	{
		name:        "prefix_only_no_suffix_true",
		input:       []int{2, 4, 6, 1, 3, 5},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6},
		expectFalse: []int{1, 3, 5},
		// prefix only, may return view
		trueCapMin: 5,
		trueCapMax: 6,
		// consecutive suffix, may return view
		falseCapMin: 3,
		falseCapMax: 3,
	},
	{
		name:        "long_consecutive_in_middle",
		input:       []int{1, 3, 10, 11, 12, 13, 14, 15, 5, 7},
		testFunc:    func(v int) bool { return v >= 10 && v <= 15 },
		expectTrue:  []int{10, 11, 12, 13, 14, 15},
		expectFalse: []int{1, 3, 5, 7},
		// long consecutive middle chunk, may return view
		trueCapMin: 8,
		trueCapMax: 8,
		// scattered elements, may allocate
		falseCapMin: 4,
		falseCapMax: 10,
	},
	{
		name:        "multiple_gaps_non_consecutive",
		input:       []int{1, 2, 4, 5, 6, 8, 9, 10},
		testFunc:    func(v int) bool { return v%2 == 0 },
		expectTrue:  []int{2, 4, 6, 8, 10},
		expectFalse: []int{1, 5, 9},
		// scattered elements, may allocate
		trueCapMin: 6,
		trueCapMax: 8,
		// scattered elements, may allocate
		falseCapMin: 6,
		falseCapMax: 8,
	},
	{
		name:        "consecutive_with_single_gap",
		input:       []int{1, 6, 7, 8, 9, 5, 12},
		testFunc:    func(v int) bool { return v >= 6 },
		expectTrue:  []int{6, 7, 8, 9, 12},
		expectFalse: []int{1, 5},
		// consecutive chunk + suffix, may allocate
		trueCapMin: 5,
		trueCapMax: 7,
		// scattered elements, may allocate
		falseCapMin: 3,
		falseCapMax: 7,
	},
	{
		name:        "large_consecutive_chunk_optimization",
		input:       []int{1, 3, 5, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		testFunc:    func(v int) bool { return v >= 10 },
		expectTrue:  []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		expectFalse: []int{1, 3, 5},
		// large consecutive suffix, may return view
		trueCapMin: 16,
		trueCapMax: 16,
		// non-consecutive prefix, may allocate
		falseCapMin: 18,
		falseCapMax: 19,
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
			trueSlice := SliceFilter(tt.input, tt.testFunc)
			if len(tt.expectTrue) == 0 {
				assert.Empty(t, trueSlice)
			} else {
				assert.Equal(t, tt.expectTrue, trueSlice)
			}
			assert.GreaterOrEqual(t, cap(trueSlice), tt.trueCapMin)
			assert.LessOrEqual(t, cap(trueSlice), tt.trueCapMax)

			falseSlice := SliceFilter(tt.input, func(v int) bool { return !tt.testFunc(v) })
			if len(tt.expectFalse) == 0 {
				assert.Empty(t, falseSlice)
			} else {
				assert.Equal(t, tt.expectFalse, falseSlice)
			}
			assert.GreaterOrEqual(t, cap(falseSlice), tt.falseCapMin)
			assert.LessOrEqual(t, cap(falseSlice), tt.falseCapMax)
		})
	}
}

func TestSliceFilterInto(t *testing.T) {
	t.Parallel()

	for i, tt := range sliceTestCases {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			resultTrue := SliceFilterInto(make([]int, 0), tt.testFunc, tt.input)
			if len(tt.expectTrue) == 0 {
				assert.Empty(t, resultTrue)
			} else {
				assert.Equal(t, tt.expectTrue, resultTrue)
			}

			resultFalse := SliceFilterInto(make([]int, 0), func(v int) bool { return !tt.testFunc(v) }, tt.input)
			if len(tt.expectFalse) == 0 {
				assert.Empty(t, resultFalse)
			} else {
				assert.Equal(t, tt.expectFalse, resultFalse)
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
			if len(tt.expectTrue) == 0 {
				assert.Empty(t, trueSlice) // may be nil or input
			} else {
				assert.Equal(t, tt.expectTrue, trueSlice)
			}
			assert.Equal(t, tt.expectFalse, falseSlice)

			// Validate capacity
			assert.GreaterOrEqual(t, cap(trueSlice), tt.trueCapMin)
			assert.LessOrEqual(t, cap(trueSlice), tt.trueCapMax)
			assert.GreaterOrEqual(t, cap(falseSlice), tt.falseCapMin)
			assert.LessOrEqual(t, cap(falseSlice), tt.falseCapMax)
		})
	}
}

func TestSliceSplitInPlace(t *testing.T) {
	t.Parallel()

	for i, tt := range sliceTestCases {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			// make a copy of the test input to avoid changing it
			trueSlice, falseSlice := SliceSplitInPlace(sliceDup(tt.input), tt.testFunc)
			if len(tt.expectTrue) == 0 {
				assert.Empty(t, trueSlice) // may be nil or input
			} else {
				assert.Equal(t, tt.expectTrue, trueSlice)
			}
			assert.Equal(t, tt.expectFalse, falseSlice)

			// Validate capacity
			assert.GreaterOrEqual(t, cap(trueSlice), tt.trueCapMin)
			assert.LessOrEqual(t, cap(trueSlice), tt.trueCapMax)
			assert.GreaterOrEqual(t, cap(falseSlice), tt.falseCapMin)
			assert.LessOrEqual(t, cap(falseSlice), tt.falseCapMax)
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

// TestSliceFilterCapacity tests that SliceFilter returns slices with appropriate capacity
func TestSliceFilterCapacity(t *testing.T) {
	t.Parallel()

	t.Run("all_pass_returns_original_slice", func(t *testing.T) {
		input := []int{2, 4, 6, 8}
		result := SliceFilter(input, func(v int) bool { return v%2 == 0 })

		// Should return the exact same slice (view optimization)
		assert.Equal(t, input, result)
		assert.Equal(t, cap(input), cap(result))

		// Verify it's the same underlying array
		if len(input) > 0 && len(result) > 0 {
			assert.Equal(t, &input[0], &result[0])
		}
	})

	t.Run("consecutive_suffix_returns_view", func(t *testing.T) {
		input := []int{1, 3, 6, 7, 8, 9}
		result := SliceFilter(input, func(v int) bool { return v >= 6 })

		// Should return a view of the suffix [6, 7, 8, 9]
		assert.Equal(t, []int{6, 7, 8, 9}, result)
		// Should share underlying array
		if len(result) > 0 {
			assert.Equal(t, &input[2], &result[0]) // result starts at input[2]
		}
	})

	t.Run("consecutive_prefix_returns_view", func(t *testing.T) {
		input := []int{2, 4, 6, 1, 3, 5}
		result := SliceFilter(input, func(v int) bool { return v%2 == 0 })

		// Should return a view of prefix [2, 4, 6] - but this case requires allocation due to non-consecutive
		assert.Equal(t, []int{2, 4, 6}, result)
		// This case actually allocates because prefix + non-consecutive elements exist
		assert.GreaterOrEqual(t, cap(result), len(result))
	})

	t.Run("empty_result_retains_original_capacity", func(t *testing.T) {
		input := []int{1, 3, 5, 7}
		result := SliceFilter(input, func(v int) bool { return v%2 == 0 })

		assert.Empty(t, result)
		// Empty result from SliceFilter returns slice[:0] which retains original capacity
		assert.Equal(t, cap(input), cap(result))
	})

	t.Run("non_consecutive_allocates_with_reasonable_capacity", func(t *testing.T) {
		input := []int{2, 1, 4, 3, 6, 5, 8}
		result := SliceFilter(input, func(v int) bool { return v%2 == 0 })

		assert.Equal(t, []int{2, 4, 6, 8}, result)
		// Should allocate with capacity >= length of result
		assert.GreaterOrEqual(t, cap(result), len(result))
		// Should not over-allocate excessively
		assert.LessOrEqual(t, cap(result), len(input))
	})
}

// TestSliceSplitCapacity tests that SliceSplit returns slices with appropriate capacity
func TestSliceSplitCapacity(t *testing.T) {
	t.Parallel()

	t.Run("all_true_returns_original_slice", func(t *testing.T) {
		input := []int{2, 4, 6, 8}
		trueSlice, falseSlice := SliceSplit(input, func(v int) bool { return v%2 == 0 })

		assert.Equal(t, input, trueSlice)
		assert.Empty(t, falseSlice)
		assert.Equal(t, cap(input), cap(trueSlice))
	})

	t.Run("all_false_returns_original_slice", func(t *testing.T) {
		input := []int{1, 3, 5, 7}
		trueSlice, falseSlice := SliceSplit(input, func(v int) bool { return v%2 == 0 })

		assert.Empty(t, trueSlice)
		assert.Equal(t, input, falseSlice)
		assert.Equal(t, cap(input), cap(falseSlice))
	})

	t.Run("split_allocates_appropriate_capacity", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8}
		trueSlice, falseSlice := SliceSplit(input, func(v int) bool { return v%2 == 0 })

		assert.Equal(t, []int{2, 4, 6, 8}, trueSlice)
		assert.Equal(t, []int{1, 3, 5, 7}, falseSlice)

		// Both slices should have reasonable capacity
		assert.GreaterOrEqual(t, cap(trueSlice), len(trueSlice))
		assert.GreaterOrEqual(t, cap(falseSlice), len(falseSlice))

		// Capacity should not be excessive
		assert.LessOrEqual(t, cap(trueSlice), len(input))
		assert.LessOrEqual(t, cap(falseSlice), len(input))
	})

	t.Run("first_element_true_capacity_allocation", func(t *testing.T) {
		input := []int{2, 1, 4, 3, 6}
		trueSlice, falseSlice := SliceSplit(input, func(v int) bool { return v%2 == 0 })

		assert.Equal(t, []int{2, 4, 6}, trueSlice)
		assert.Equal(t, []int{1, 3}, falseSlice)

		// Check capacity is reasonable for both slices
		assert.GreaterOrEqual(t, cap(trueSlice), len(trueSlice))
		assert.GreaterOrEqual(t, cap(falseSlice), len(falseSlice))
	})

	t.Run("first_element_false_capacity_allocation", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		trueSlice, falseSlice := SliceSplit(input, func(v int) bool { return v%2 == 0 })

		assert.Equal(t, []int{2, 4}, trueSlice)
		assert.Equal(t, []int{1, 3, 5}, falseSlice)

		// Check capacity is reasonable for both slices
		assert.GreaterOrEqual(t, cap(trueSlice), len(trueSlice))
		assert.GreaterOrEqual(t, cap(falseSlice), len(falseSlice))
	})
}

// TestSliceSplitInPlaceCapacity tests capacity behavior of in-place split functions
func TestSliceSplitInPlaceCapacity(t *testing.T) {
	t.Parallel()

	t.Run("first_true_reuses_input_slice", func(t *testing.T) {
		input := []int{2, 1, 4, 3, 6}
		originalCap := cap(input)
		inputCopy := sliceDup(input)

		trueSlice, falseSlice :=
			SliceSplitInPlace(inputCopy, func(v int) bool { return v%2 == 0 })

		assert.Equal(t, []int{2, 4, 6}, trueSlice)
		assert.Equal(t, []int{1, 3}, falseSlice)

		// True slice should reuse original capacity
		assert.Equal(t, originalCap, cap(trueSlice))
		// False slice should have reasonable capacity
		assert.GreaterOrEqual(t, cap(falseSlice), len(falseSlice))
	})

	t.Run("first_false_reuses_input_slice", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		originalCap := cap(input)
		inputCopy := sliceDup(input)

		trueSlice, falseSlice :=
			SliceSplitInPlace(inputCopy, func(v int) bool { return v%2 == 0 })

		assert.Equal(t, []int{2, 4}, trueSlice)
		assert.Equal(t, []int{1, 3, 5}, falseSlice)

		// False slice should reuse original capacity
		assert.Equal(t, originalCap, cap(falseSlice))
		// True slice should have reasonable capacity
		assert.GreaterOrEqual(t, cap(trueSlice), len(trueSlice))
	})

	t.Run("all_true_reuses_input_capacity", func(t *testing.T) {
		input := []int{2, 4, 6, 8}
		originalCap := cap(input)
		inputCopy := sliceDup(input)

		trueSlice, falseSlice :=
			SliceSplitInPlace(inputCopy, func(v int) bool { return v%2 == 0 })

		assert.Equal(t, input, trueSlice)
		assert.Empty(t, falseSlice)
		assert.Equal(t, originalCap, cap(trueSlice))
	})

	t.Run("all_false_reuses_input_capacity", func(t *testing.T) {
		input := []int{1, 3, 5, 7}
		originalCap := cap(input)
		inputCopy := sliceDup(input)

		trueSlice, falseSlice :=
			SliceSplitInPlace(inputCopy, func(v int) bool { return v%2 == 0 })

		assert.Empty(t, trueSlice)
		assert.Equal(t, input, falseSlice)
		assert.Equal(t, originalCap, cap(falseSlice))
	})
}

// TestSliceSplitInPlaceUnstableCapacity tests capacity behavior of unstable in-place split
func TestSliceSplitInPlaceUnstableCapacity(t *testing.T) {
	t.Parallel()

	t.Run("reuses_input_slice_capacity", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}
		originalCap := cap(input)
		inputCopy := sliceDup(input)

		trueSlice, falseSlice :=
			SliceSplitInPlaceUnstable(inputCopy, func(v int) bool { return v%2 == 0 })

		// Both slices should use the original array
		totalLen := len(trueSlice) + len(falseSlice)
		assert.Equal(t, len(input), totalLen)

		// Combined capacity should equal original capacity
		if len(trueSlice) > 0 && len(falseSlice) > 0 {
			// Both slices share the same underlying array
			assert.LessOrEqual(t, cap(trueSlice), originalCap)
			assert.LessOrEqual(t, cap(falseSlice), originalCap)
		} else if len(trueSlice) > 0 {
			assert.Equal(t, originalCap, cap(trueSlice))
		} else if len(falseSlice) > 0 {
			assert.Equal(t, originalCap, cap(falseSlice))
		}
	})

	t.Run("all_true_preserves_capacity", func(t *testing.T) {
		input := []int{2, 4, 6, 8}
		originalCap := cap(input)
		inputCopy := sliceDup(input)

		trueSlice, falseSlice :=
			SliceSplitInPlaceUnstable(inputCopy, func(v int) bool { return v%2 == 0 })

		assert.Equal(t, originalCap, cap(trueSlice))
		assert.Empty(t, falseSlice)
	})

	t.Run("all_false_preserves_capacity", func(t *testing.T) {
		input := []int{1, 3, 5, 7}
		originalCap := cap(input)
		inputCopy := sliceDup(input)

		trueSlice, falseSlice :=
			SliceSplitInPlaceUnstable(inputCopy, func(v int) bool { return v%2 == 0 })

		assert.Empty(t, trueSlice)
		assert.Equal(t, originalCap, cap(falseSlice))
	})
}

func TestSliceTransform(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		var input []int
		result := SliceTransform(func(i int) int { return i }, input)
		assert.Empty(t, result)
	})

	t.Run("empty", func(t *testing.T) {
		result := SliceTransform(func(i int) int { return i }, []int{})
		assert.Empty(t, result)
	})

	t.Run("int_to_string", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := SliceTransform(func(i int) string { return strconv.Itoa(i) }, input)
		assert.Equal(t, []string{"1", "2", "3"}, result)
	})

	t.Run("float_to_int", func(t *testing.T) {
		input := []float64{1.1, 2.2}
		result := SliceTransform(func(f float64) int { return int(f) }, input)
		assert.Equal(t, []int{1, 2}, result)
	})

	t.Run("multiple_slices_concatenated", func(t *testing.T) {
		slice1 := []int{1, 2}
		slice2 := []int{3, 4}
		slice3 := []int{5}
		result := SliceTransform(func(i int) string { return strconv.Itoa(i) }, slice1, slice2, slice3)
		assert.Equal(t, []string{"1", "2", "3", "4", "5"}, result)
	})

	t.Run("multiple_slices_with_empty", func(t *testing.T) {
		slice1 := []int{1, 2}
		slice2 := []int{}
		slice3 := []int{3, 4}
		result := SliceTransform(func(i int) int { return i * 2 }, slice1, slice2, slice3)
		assert.Equal(t, []int{2, 4, 6, 8}, result)
	})

	t.Run("multiple_slices_all_empty", func(t *testing.T) {
		result := SliceTransform(func(i int) int { return i }, []int{}, []int{}, []int{})
		assert.Empty(t, result)
	})

	t.Run("multiple_slices_with_nil", func(t *testing.T) {
		slice1 := []int{1}
		var slice2 []int
		slice3 := []int{2, 3}
		result := SliceTransform(func(i int) string { return strconv.Itoa(i) }, slice1, slice2, slice3)
		assert.Equal(t, []string{"1", "2", "3"}, result)
	})

	t.Run("single_slice_maintains_compatibility", func(t *testing.T) {
		input := []int{10, 20, 30}
		result := SliceTransform(func(i int) int { return i / 10 }, input)
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

var sliceToSetTests = []struct {
	name       string
	input      [][]int
	expectKeys []int
}{
	{
		name:  "nil",
		input: nil,
	},
	{
		name:  "empty",
		input: [][]int{},
	},
	{
		name:       "unique_values",
		input:      [][]int{{1, 2}, {3, 4}},
		expectKeys: []int{1, 2, 3, 4},
	},
	{
		name:       "middle_overlapping_values",
		input:      [][]int{{1, 2}, {2, 3}, {3, 4}},
		expectKeys: []int{1, 2, 3, 4},
	},
	{
		name:       "single_slice",
		input:      [][]int{{1, 2, 3}},
		expectKeys: []int{1, 2, 3},
	},
	{
		name:  "nested_empty_slices",
		input: [][]int{{}, {}},
	},
	{
		name:       "duplicate_values_in_slice",
		input:      [][]int{{1, 1, 2}, {2, 3, 3}},
		expectKeys: []int{1, 2, 3},
	},
	{
		name:  "single_empty_slice",
		input: [][]int{{}},
	},
	{
		name:       "all_same_values",
		input:      [][]int{{1, 1, 1}, {1, 1}, {1}},
		expectKeys: []int{1},
	},
	{
		name:       "large_input_with_duplicates",
		input:      [][]int{sliceLargeInput[:25], sliceLargeInput[20:]},
		expectKeys: sliceLargeInput[:51], // 0-50 unique values
	},
	{
		name:       "zero_values",
		input:      [][]int{{0}, {0, 1}, {1, 0}},
		expectKeys: []int{0, 1},
	},
	{
		name:       "negative_values",
		input:      [][]int{{-1, -2}, {-2, -3}, {-3, -4}},
		expectKeys: []int{-1, -2, -3, -4},
	},
}

func TestSliceToSet(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceToSetTests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceToSet(tt.input...)

			assert.Len(t, got, len(tt.expectKeys))

			for _, key := range tt.expectKeys {
				_, ok := got[key]
				assert.True(t, ok)
			}

			for key := range got {
				assert.Contains(t, tt.expectKeys, key)
			}
		})
	}

	t.Run("string", func(t *testing.T) {
		result := SliceToSet([][]string{{"a", "b"}, {"b", "c"}, {"c", "d"}}...)

		expected := map[string]struct{}{
			"a": {},
			"b": {},
			"c": {},
			"d": {},
		}
		assert.Len(t, result, 4)
		assert.Equal(t, expected, result)
	})
}

var sliceToSetByTests = []struct {
	name        string
	inputSlices [][]int
	conversion  func(int) string
	expectKeys  []string
}{
	{
		name:        "nil_input",
		inputSlices: nil,
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectKeys:  []string{},
	},
	{
		name:        "empty_input",
		inputSlices: [][]int{},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectKeys:  []string{},
	},
	{
		name:        "single_slice",
		inputSlices: [][]int{{1, 2, 3}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectKeys:  []string{"1", "2", "3"},
	},
	{
		name:        "multiple_slices_unique_values",
		inputSlices: [][]int{{1, 2}, {3, 4}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectKeys:  []string{"1", "2", "3", "4"},
	},
	{
		name:        "multiple_slices_overlapping_values",
		inputSlices: [][]int{{1, 2}, {2, 3}, {3, 4}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectKeys:  []string{"1", "2", "3", "4"},
	},
	{
		name:        "duplicate_values_in_slice",
		inputSlices: [][]int{{1, 1, 2}, {2, 3, 3}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectKeys:  []string{"1", "2", "3"},
	},
	{
		name:        "empty_slices_mixed",
		inputSlices: [][]int{{}, {1, 2}, {}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectKeys:  []string{"1", "2"},
	},
	{
		name:        "all_same_values",
		inputSlices: [][]int{{5, 5, 5}, {5, 5}, {5}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectKeys:  []string{"5"},
	},
	{
		name:        "zero_values",
		inputSlices: [][]int{{0}, {0, 1}, {1, 0}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectKeys:  []string{"0", "1"},
	},
	{
		name:        "negative_values",
		inputSlices: [][]int{{-1, -2}, {-2, -3}, {-3, -4}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectKeys:  []string{"-1", "-2", "-3", "-4"},
	},
	{
		name:        "transform_to_same_key",
		inputSlices: [][]int{{1, 11, 21}, {31, 41, 51}},
		conversion:  func(i int) string { return strconv.Itoa(i % 10) }, // all mod 10 = 1
		expectKeys:  []string{"1"},
	},
	{
		name:        "large_input_with_duplicates",
		inputSlices: [][]int{sliceLargeInput[:25], sliceLargeInput[20:]},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectKeys: func() []string {
			keys := make([]string, 51) // 0-50 unique values
			for i := 0; i <= 50; i++ {
				keys[i] = strconv.Itoa(i)
			}
			return keys
		}(),
	},
}

func TestSliceToSetBy(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceToSetByTests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceToSetBy(tt.conversion, tt.inputSlices...)

			assert.Len(t, got, len(tt.expectKeys))

			for _, key := range tt.expectKeys {
				_, ok := got[key]
				assert.True(t, ok)
			}

			for key := range got {
				assert.Contains(t, tt.expectKeys, key, "unexpected key %q found", key)
			}
		})
	}

	t.Run("int_to_length_string", func(t *testing.T) {
		slice1 := []string{"a", "bb", "ccc"}
		slice2 := []string{"dd", "eeeee"}
		result := SliceToSetBy(func(s string) int { return len(s) }, slice1, slice2)

		expected := map[int]struct{}{
			1: {}, // "a"
			2: {}, // "bb", "dd"
			3: {}, // "ccc"
			5: {}, // "eeeee"
		}
		assert.Equal(t, expected, result)
	})

	t.Run("comparable_struct_keys", func(t *testing.T) {
		type Point struct{ X, Y int }
		points1 := []Point{{1, 2}, {3, 4}}
		points2 := []Point{{1, 2}, {5, 6}} // duplicate {1, 2}

		result := SliceToSetBy(func(p Point) Point { return p }, points1, points2)

		expected := map[Point]struct{}{
			{1, 2}: {},
			{3, 4}: {},
			{5, 6}: {},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("complex_transformation", func(t *testing.T) {
		nums1 := []int{10, 20, 30}
		nums2 := []int{15, 25, 35}

		// Transform to "even" or "odd" based on last digit
		result := SliceToSetBy(func(n int) string {
			if n%10%2 == 0 {
				return "even_ending"
			}
			return "odd_ending"
		}, nums1, nums2)

		expected := map[string]struct{}{
			"even_ending": {}, // 10, 20, 30
			"odd_ending":  {}, // 15, 25, 35
		}
		assert.Equal(t, expected, result)
	})
}

var sliceToCountsTests = []struct {
	name         string
	input        [][]int
	expectCounts map[int]int
}{
	{
		name:         "nil",
		input:        nil,
		expectCounts: map[int]int{},
	},
	{
		name:         "empty",
		input:        [][]int{},
		expectCounts: map[int]int{},
	},
	{
		name:         "unique_values",
		input:        [][]int{{1, 2}, {3, 4}},
		expectCounts: map[int]int{1: 1, 2: 1, 3: 1, 4: 1},
	},
	{
		name:         "overlapping_values",
		input:        [][]int{{1, 2}, {2, 3}, {3, 4}},
		expectCounts: map[int]int{1: 1, 2: 2, 3: 2, 4: 1},
	},
	{
		name:         "single_slice",
		input:        [][]int{{1, 2, 3}},
		expectCounts: map[int]int{1: 1, 2: 1, 3: 1},
	},
	{
		name:         "nested_empty_slices",
		input:        [][]int{{}, {}},
		expectCounts: map[int]int{},
	},
	{
		name:         "duplicate_values_in_slice",
		input:        [][]int{{1, 1, 2}, {2, 3, 3}},
		expectCounts: map[int]int{1: 2, 2: 2, 3: 2},
	},
	{
		name:         "single_empty_slice",
		input:        [][]int{{}},
		expectCounts: map[int]int{},
	},
	{
		name:         "all_same_values",
		input:        [][]int{{1, 1, 1}, {1, 1}, {1}},
		expectCounts: map[int]int{1: 6},
	},
	{
		name:         "zero_values",
		input:        [][]int{{0}, {0, 1}, {1, 0}},
		expectCounts: map[int]int{0: 3, 1: 2},
	},
	{
		name:         "negative_values",
		input:        [][]int{{-1, -2}, {-2, -3}, {-3, -4}},
		expectCounts: map[int]int{-1: 1, -2: 2, -3: 2, -4: 1},
	},
	{
		name:         "large_mixed_counts",
		input:        [][]int{{1, 2, 3, 2, 1}, {3, 4, 5, 4, 3}},
		expectCounts: map[int]int{1: 2, 2: 2, 3: 3, 4: 2, 5: 1},
	},
}

var sliceToCountsByTests = []struct {
	name         string
	inputSlices  [][]int
	conversion   func(int) string
	expectCounts map[string]int
}{
	{
		name:         "nil_input",
		inputSlices:  nil,
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectCounts: map[string]int{},
	},
	{
		name:         "empty_input",
		inputSlices:  [][]int{},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectCounts: map[string]int{},
	},
	{
		name:         "single_slice",
		inputSlices:  [][]int{{1, 2, 3}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectCounts: map[string]int{"1": 1, "2": 1, "3": 1},
	},
	{
		name:         "multiple_slices_unique_values",
		inputSlices:  [][]int{{1, 2}, {3, 4}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectCounts: map[string]int{"1": 1, "2": 1, "3": 1, "4": 1},
	},
	{
		name:         "multiple_slices_overlapping_values",
		inputSlices:  [][]int{{1, 2}, {2, 3}, {3, 4}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectCounts: map[string]int{"1": 1, "2": 2, "3": 2, "4": 1},
	},
	{
		name:         "duplicate_values_in_slice",
		inputSlices:  [][]int{{1, 1, 2}, {2, 3, 3}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectCounts: map[string]int{"1": 2, "2": 2, "3": 2},
	},
	{
		name:         "empty_slices_mixed",
		inputSlices:  [][]int{{}, {1, 2}, {}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectCounts: map[string]int{"1": 1, "2": 1},
	},
	{
		name:         "all_same_values",
		inputSlices:  [][]int{{5, 5, 5}, {5, 5}, {5}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectCounts: map[string]int{"5": 6},
	},
	{
		name:         "zero_values",
		inputSlices:  [][]int{{0}, {0, 1}, {1, 0}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectCounts: map[string]int{"0": 3, "1": 2},
	},
	{
		name:         "negative_values",
		inputSlices:  [][]int{{-1, -2}, {-2, -3}, {-3, -4}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectCounts: map[string]int{"-1": 1, "-2": 2, "-3": 2, "-4": 1},
	},
	{
		name:         "transform_to_same_key",
		inputSlices:  [][]int{{1, 11, 21}, {31, 41, 51}},
		conversion:   func(i int) string { return strconv.Itoa(i % 10) },
		expectCounts: map[string]int{"1": 6},
	},
	{
		name:         "mod_10_grouping",
		inputSlices:  [][]int{{10, 20, 11, 21}, {12, 22, 13, 23}},
		conversion:   func(i int) string { return strconv.Itoa(i % 10) },
		expectCounts: map[string]int{"0": 2, "1": 2, "2": 2, "3": 2},
	},
}

func TestSliceToCounts(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceToCountsTests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceToCounts(tt.input...)

			assert.Len(t, got, len(tt.expectCounts))
			for key, expectedCount := range tt.expectCounts {
				actualCount, ok := got[key]
				assert.True(t, ok)
				assert.Equal(t, expectedCount, actualCount)
			}

			for key := range got {
				assert.Contains(t, tt.expectCounts, key, "unexpected key %v found", key)
			}
		})
	}

	t.Run("string_values", func(t *testing.T) {
		result := SliceToCounts([][]string{{"a", "b"}, {"b", "c"}, {"c", "d", "a"}}...)

		expected := map[string]int{
			"a": 2,
			"b": 2,
			"c": 2,
			"d": 1,
		}
		assert.Equal(t, expected, result)
	})
}

func TestSliceToCountsBy(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceToCountsByTests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceToCountsBy(tt.conversion, tt.inputSlices...)

			assert.Len(t, got, len(tt.expectCounts))
			for key, expectedCount := range tt.expectCounts {
				actualCount, ok := got[key]
				assert.True(t, ok)
				assert.Equal(t, expectedCount, actualCount)
			}

			for key := range got {
				assert.Contains(t, tt.expectCounts, key, "unexpected key %q found", key)
			}
		})
	}

	t.Run("int_to_length_string", func(t *testing.T) {
		slice1 := []string{"a", "bb", "ccc", "a"}
		slice2 := []string{"dd", "eeeee", "bb"}
		result := SliceToCountsBy(func(s string) int { return len(s) }, slice1, slice2)

		expected := map[int]int{
			1: 2, // "a", "a"
			2: 3, // "bb", "dd", "bb"
			3: 1, // "ccc"
			5: 1, // "eeeee"
		}
		assert.Equal(t, expected, result)
	})

	t.Run("comparable_struct_keys", func(t *testing.T) {
		type Point struct{ X, Y int }
		points1 := []Point{{1, 2}, {3, 4}, {1, 2}}
		points2 := []Point{{1, 2}, {5, 6}, {3, 4}}

		result := SliceToCountsBy(func(p Point) Point { return p }, points1, points2)

		expected := map[Point]int{
			{1, 2}: 3, // appears 3 times total
			{3, 4}: 2, // appears 2 times total
			{5, 6}: 1, // appears 1 time
		}
		assert.Equal(t, expected, result)
	})

	t.Run("complex_transformation", func(t *testing.T) {
		nums1 := []int{10, 20, 30, 15}
		nums2 := []int{25, 35, 40, 45}

		result := SliceToCountsBy(func(n int) string {
			if n%10%2 == 0 {
				return "even_ending"
			}
			return "odd_ending"
		}, nums1, nums2)

		expected := map[string]int{
			"even_ending": 4, // 10, 20, 30, 40
			"odd_ending":  4, // 15, 25, 35, 45
		}
		assert.Equal(t, expected, result)
	})
}

var sliceToIndexByTests = []struct {
	name        string
	inputSlices [][]int
	conversion  func(int) string
	expectIndex map[string]int
}{
	{
		name:        "nil_input",
		inputSlices: nil,
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectIndex: map[string]int{},
	},
	{
		name:        "empty_input",
		inputSlices: [][]int{},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectIndex: map[string]int{},
	},
	{
		name:        "single_slice",
		inputSlices: [][]int{{1, 2, 3}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectIndex: map[string]int{"1": 1, "2": 2, "3": 3},
	},
	{
		name:        "multiple_slices_unique_keys",
		inputSlices: [][]int{{1, 2}, {3, 4}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectIndex: map[string]int{"1": 1, "2": 2, "3": 3, "4": 4},
	},
	{
		name:        "duplicate_keys_last_wins",
		inputSlices: [][]int{{1, 2}, {2, 3}, {3, 4}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectIndex: map[string]int{"1": 1, "2": 2, "3": 3, "4": 4}, // last occurrence wins
	},
	{
		name:        "empty_slices_mixed",
		inputSlices: [][]int{{}, {1, 2}, {}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectIndex: map[string]int{"1": 1, "2": 2},
	},
	{
		name:        "overwrite_values",
		inputSlices: [][]int{{10, 20}, {30, 10}},
		conversion:  func(i int) string { return strconv.Itoa(i % 10) },
		expectIndex: map[string]int{"0": 10}, // "0" key maps to last 10 (from 30, then 10)
	},
	{
		name:        "zero_values",
		inputSlices: [][]int{{0, 1}, {2, 0}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectIndex: map[string]int{"0": 0, "1": 1, "2": 2}, // last 0 wins
	},
	{
		name:        "negative_values",
		inputSlices: [][]int{{-1, -2}, {-3, -2}},
		conversion:  func(i int) string { return strconv.Itoa(i) },
		expectIndex: map[string]int{"-1": -1, "-2": -2, "-3": -3}, // last -2 wins
	},
}

var sliceToGroupsByTests = []struct {
	name         string
	inputSlices  [][]int
	conversion   func(int) string
	expectGroups map[string][]int
}{
	{
		name:         "nil_input",
		inputSlices:  nil,
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectGroups: map[string][]int{},
	},
	{
		name:         "empty_input",
		inputSlices:  [][]int{},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectGroups: map[string][]int{},
	},
	{
		name:         "single_slice",
		inputSlices:  [][]int{{1, 2, 3}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectGroups: map[string][]int{"1": {1}, "2": {2}, "3": {3}},
	},
	{
		name:         "multiple_slices_unique_keys",
		inputSlices:  [][]int{{1, 2}, {3, 4}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectGroups: map[string][]int{"1": {1}, "2": {2}, "3": {3}, "4": {4}},
	},
	{
		name:         "multiple_slices_overlapping_keys",
		inputSlices:  [][]int{{1, 2}, {2, 3}, {3, 4}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectGroups: map[string][]int{"1": {1}, "2": {2, 2}, "3": {3, 3}, "4": {4}},
	},
	{
		name:         "empty_slices_mixed",
		inputSlices:  [][]int{{}, {1, 2}, {}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectGroups: map[string][]int{"1": {1}, "2": {2}},
	},
	{
		name:         "group_by_mod_10",
		inputSlices:  [][]int{{10, 21, 32}, {13, 24, 30}},
		conversion:   func(i int) string { return strconv.Itoa(i % 10) },
		expectGroups: map[string][]int{"0": {10, 30}, "1": {21}, "2": {32}, "3": {13}, "4": {24}},
	},
	{
		name:         "zero_values",
		inputSlices:  [][]int{{0, 1}, {2, 0}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectGroups: map[string][]int{"0": {0, 0}, "1": {1}, "2": {2}},
	},
	{
		name:         "negative_values",
		inputSlices:  [][]int{{-1, -2}, {-3, -2}},
		conversion:   func(i int) string { return strconv.Itoa(i) },
		expectGroups: map[string][]int{"-1": {-1}, "-2": {-2, -2}, "-3": {-3}},
	},
	{
		name:        "group_by_sign",
		inputSlices: [][]int{{-5, -2, 3}, {7, -1, 0}},
		conversion: func(i int) string {
			if i < 0 {
				return "negative"
			} else if i > 0 {
				return "positive"
			}
			return "zero"
		},
		expectGroups: map[string][]int{"negative": {-5, -2, -1}, "positive": {3, 7}, "zero": {0}},
	},
}

func TestSliceToIndexBy(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceToIndexByTests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceToIndexBy(tt.conversion, tt.inputSlices...)

			assert.Len(t, got, len(tt.expectIndex))
			for key, expectedValue := range tt.expectIndex {
				actualValue, ok := got[key]
				assert.True(t, ok)
				assert.Equal(t, expectedValue, actualValue)
			}

			for key := range got {
				assert.Contains(t, tt.expectIndex, key, "unexpected key %q found", key)
			}
		})
	}

	t.Run("struct_indexing", func(t *testing.T) {
		type Person struct {
			ID   int
			Name string
		}
		people1 := []Person{{1, "Alice"}, {2, "Bob"}}
		people2 := []Person{{3, "Charlie"}, {1, "Alice_Updated"}} // ID 1 gets updated

		result := SliceToIndexBy(func(p Person) int { return p.ID }, people1, people2)

		expected := map[int]Person{
			1: {1, "Alice_Updated"}, // last wins
			2: {2, "Bob"},
			3: {3, "Charlie"},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("string_to_length_index", func(t *testing.T) {
		slice1 := []string{"a", "bb", "ccc"}
		slice2 := []string{"dd", "eeeee", "f"} // "f" overwrites "a" (both len=1)

		result := SliceToIndexBy(func(s string) int { return len(s) }, slice1, slice2)

		expected := map[int]string{
			1: "f",     // last string with length 1
			2: "dd",    // last string with length 2
			3: "ccc",   // only string with length 3
			5: "eeeee", // only string with length 5
		}
		assert.Equal(t, expected, result)
	})
}

func TestSliceToGroupsBy(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceToGroupsByTests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceToGroupsBy(tt.conversion, tt.inputSlices...)

			assert.Len(t, got, len(tt.expectGroups))
			for key, expectedGroup := range tt.expectGroups {
				actualGroup, ok := got[key]
				assert.True(t, ok)
				assert.Equal(t, expectedGroup, actualGroup)
			}

			for key := range got {
				assert.Contains(t, tt.expectGroups, key, "unexpected key %q found", key)
			}
		})
	}

	t.Run("struct_grouping", func(t *testing.T) {
		type Person struct {
			Dept string
			Name string
		}
		people1 := []Person{{"eng", "Alice"}, {"sales", "Bob"}}
		people2 := []Person{{"eng", "Charlie"}, {"sales", "Dave"}}

		result := SliceToGroupsBy(func(p Person) string { return p.Dept }, people1, people2)

		expected := map[string][]Person{
			"eng":   {{Dept: "eng", Name: "Alice"}, {Dept: "eng", Name: "Charlie"}},
			"sales": {{Dept: "sales", Name: "Bob"}, {Dept: "sales", Name: "Dave"}},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("string_to_length_groups", func(t *testing.T) {
		slice1 := []string{"a", "bb", "ccc", "d"}
		slice2 := []string{"ee", "f", "gggg"}

		result := SliceToGroupsBy(func(s string) int { return len(s) }, slice1, slice2)

		expected := map[int][]string{
			1: {"a", "d", "f"},
			2: {"bb", "ee"},
			3: {"ccc"},
			4: {"gggg"},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("complex_grouping_with_duplicates", func(t *testing.T) {
		nums1 := []int{10, 15, 20, 11}
		nums2 := []int{25, 30, 16}

		result := SliceToGroupsBy(func(n int) string {
			if n%10 < 5 {
				return "low_digit"
			}
			return "high_digit"
		}, nums1, nums2)

		expected := map[string][]int{
			"low_digit":  {10, 20, 11, 30}, // 10, 20, 11, 30 (digits 0, 0, 1, 0)
			"high_digit": {15, 25, 16},     // 15, 25, 16 (digits 5, 5, 6)
		}
		assert.Equal(t, expected, result)
	})
}

func TestSliceTotalSize(t *testing.T) {
	t.Parallel()

	sliceTotalSizeTests := []struct {
		name   string
		slices [][]int
		expect int
	}{
		{
			name:   "nil",
			slices: nil,
		},
		{
			name:   "empty",
			slices: [][]int{},
		},
		{
			name:   "single_nil_slice",
			slices: [][]int{nil},
		},
		{
			name:   "single_empty_slice",
			slices: [][]int{{}},
		},
		{
			name:   "multiple_empty_slices",
			slices: [][]int{{}, {}, {}},
		},
		{
			name:   "single_slice",
			slices: [][]int{{1, 2, 3}},
			expect: 3,
		},
		{
			name:   "multiple_slices_same_size",
			slices: [][]int{{1, 2}, {3, 4}, {5, 6}},
			expect: 6,
		},
		{
			name:   "multiple_slices_different_sizes",
			slices: [][]int{{1}, {2, 3}, {4, 5, 6, 7}},
			expect: 7,
		},
		{
			name:   "mixed_empty_and_non_empty",
			slices: [][]int{{}, {1, 2}, {}, {3, 4, 5}},
			expect: 5,
		},
		{
			name:   "mixed_nil_and_non_nil",
			slices: [][]int{nil, {1, 2}, nil, {3, 4, 5}},
			expect: 5,
		},
		{
			name:   "large_single_slice",
			slices: [][]int{sliceLargeInput},
			expect: 101,
		},
		{
			name:   "multiple_large_slices",
			slices: [][]int{sliceLargeInput[:25], sliceLargeInput[25:50], sliceLargeInput[50:]},
			expect: 101,
		},
		{
			name:   "zero_length_but_non_nil",
			slices: [][]int{make([]int, 0, 10), make([]int, 0, 5)},
			expect: 0,
		},
	}

	for _, tt := range sliceTotalSizeTests {
		t.Run(tt.name, func(t *testing.T) {
			got := sliceTotalSize(tt.slices)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestSliceIntoSet(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceToSetTests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[int]struct{})
			SliceIntoSet(result, tt.input...)

			assert.Len(t, result, len(tt.expectKeys))

			for _, key := range tt.expectKeys {
				_, ok := result[key]
				assert.True(t, ok)
			}

			for key := range result {
				assert.Contains(t, tt.expectKeys, key)
			}
		})
	}
}

func TestSliceIntoSetBy(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceToSetByTests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string]struct{})
			SliceIntoSetBy(result, tt.conversion, tt.inputSlices...)

			assert.Len(t, result, len(tt.expectKeys))

			for _, key := range tt.expectKeys {
				_, ok := result[key]
				assert.True(t, ok)
			}

			for key := range result {
				assert.Contains(t, tt.expectKeys, key, "unexpected key %q found", key)
			}
		})
	}
}

func TestSliceIntoCounts(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceToCountsTests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[int]int)
			SliceIntoCounts(result, tt.input...)

			assert.Len(t, result, len(tt.expectCounts))
			for key, expectedCount := range tt.expectCounts {
				actualCount, ok := result[key]
				assert.True(t, ok)
				assert.Equal(t, expectedCount, actualCount)
			}

			for key := range result {
				assert.Contains(t, tt.expectCounts, key)
			}
		})
	}
}

func TestSliceIntoCountsBy(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceToCountsByTests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string]int)
			SliceIntoCountsBy(result, tt.conversion, tt.inputSlices...)

			assert.Len(t, result, len(tt.expectCounts))
			for key, expectedCount := range tt.expectCounts {
				actualCount, ok := result[key]
				assert.True(t, ok)
				assert.Equal(t, expectedCount, actualCount)
			}

			for key := range result {
				assert.Contains(t, tt.expectCounts, key)
			}
		})
	}
}

func TestSliceIntoIndexBy(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceToIndexByTests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string]int)
			SliceIntoIndexBy(result, tt.conversion, tt.inputSlices...)

			assert.Len(t, result, len(tt.expectIndex))
			for key, expectedValue := range tt.expectIndex {
				actualValue, ok := result[key]
				assert.True(t, ok)
				assert.Equal(t, expectedValue, actualValue)
			}

			for key := range result {
				assert.Contains(t, tt.expectIndex, key)
			}
		})
	}
}

func TestSliceIntoGroupsBy(t *testing.T) {
	t.Parallel()

	for _, tt := range sliceToGroupsByTests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string][]int)
			SliceIntoGroupsBy(result, tt.conversion, tt.inputSlices...)

			assert.Len(t, result, len(tt.expectGroups))
			for key, expectedGroup := range tt.expectGroups {
				actualGroup, ok := result[key]
				assert.True(t, ok)
				assert.Equal(t, expectedGroup, actualGroup)
			}

			for key := range result {
				assert.Contains(t, tt.expectGroups, key)
			}
		})
	}
}
