package bulk

import (
	"testing"
)

// -- STATIC INPUT TESTS --

func BenchmarkSliceFilters(b *testing.B) {
	b.Run("mixed", func(b *testing.B) {
		mixedPredicate := func(v int) bool {
			return v%2 == 0
		}
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilter([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
				}, mixedPredicate)
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilterInPlace([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
				}, mixedPredicate)
			}
		})
	})
	b.Run("consecutive_start", func(b *testing.B) {
		consecutivePredicate := func(v int) bool {
			return v < 40
		}
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilter([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilterInPlace([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
	})
	b.Run("consecutive_middle", func(b *testing.B) {
		consecutivePredicate := func(v int) bool {
			return v >= 30 && v < 70
		}
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilter([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilterInPlace([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
	})
	b.Run("consecutive_end", func(b *testing.B) {
		consecutivePredicate := func(v int) bool {
			return v >= 60
		}
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilter([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilterInPlace([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
	})
}

func BenchmarkSliceSplits(b *testing.B) {
	b.Run("mixed", func(b *testing.B) {
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplit([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
				}, func(v int) bool {
					return v%2 == 0
				})
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlace([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
				}, func(v int) bool {
					return v%2 == 0
				})
			}
		})
		b.Run("InPlaceUnstable", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlaceUnstable([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
				}, func(v int) bool {
					return v%2 == 0
				})
			}
		})
	})
	b.Run("consecutive_start", func(b *testing.B) {
		consecutivePredicate := func(v int) bool {
			return v < 40
		}
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplit([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlace([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
		b.Run("InPlaceUnstable", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlaceUnstable([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
	})
	b.Run("consecutive_middle", func(b *testing.B) {
		consecutivePredicate := func(v int) bool {
			return v >= 30 && v < 70
		}
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplit([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlace([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
		b.Run("InPlaceUnstable", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlaceUnstable([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
	})
	b.Run("consecutive_end", func(b *testing.B) {
		consecutivePredicate := func(v int) bool {
			return v >= 60
		}
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplit([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlace([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
		b.Run("InPlaceUnstable", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlaceUnstable([]int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				}, consecutivePredicate)
			}
		})
	})
}

// -- TEST BELOW ARE BASED ON UNIT TEST INPUTS (VARY OVER TIME) --

func BenchmarkSliceFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceTestCases {
			_ = SliceFilter(tc.input, tc.testFunc)
		}
	}
}

func BenchmarkSliceFilterInPlace(b *testing.B) {
	// Pre-copy test case inputs to prevent modifying the actual tc.input field
	testInputs := make([][]int, len(sliceTestCases))
	for i, tc := range sliceTestCases {
		testInputs[i] = sliceDup(tc.input)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j, tc := range sliceTestCases {
			_ = SliceFilterInPlace(testInputs[j], tc.testFunc)
		}
	}
}

func BenchmarkSliceSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceTestCases {
			_, _ = SliceSplit(tc.input, tc.testFunc)
		}
	}
}

func BenchmarkSliceSplitInPlace(b *testing.B) {
	// Pre-copy test case inputs to prevent modifying the actual tc.input field
	testInputs := make([][]int, len(sliceTestCases))
	for i, tc := range sliceTestCases {
		testInputs[i] = sliceDup(tc.input)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j, tc := range sliceTestCases {
			_, _ = SliceSplitInPlace(testInputs[j], tc.testFunc)
		}
	}
}

func BenchmarkSliceSplitInPlaceUnstable(b *testing.B) {
	// Pre-copy test case inputs to prevent modifying the actual tc.input field
	testInputs := make([][]int, len(sliceTestCases))
	for i, tc := range sliceTestCases {
		testInputs[i] = sliceDup(tc.input)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j, tc := range sliceTestCases {
			_, _ = SliceSplitInPlaceUnstable(testInputs[j], tc.testFunc)
		}
	}
}

func BenchmarkSliceRemoveAt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceRemoveTestCases {
			_ = SliceRemoveAt(tc.input, tc.index)
		}
	}
}

func BenchmarkSliceRemoveAtInPlace(b *testing.B) {
	// Pre-copy test case inputs to prevent modifying the actual tc.input field
	testInputs := make([][]int, len(sliceRemoveTestCases))
	for i, tc := range sliceRemoveTestCases {
		testInputs[i] = sliceDup(tc.input)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j, tc := range sliceRemoveTestCases {
			_ = SliceRemoveAtInPlace(testInputs[j], tc.index)
		}
	}
}

func BenchmarkSliceToMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceToMapTests {
			_ = SliceToMap(tc.input...)
		}
	}
}
