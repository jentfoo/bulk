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
				_ = SliceFilter(mixedPredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
				})
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilterInPlace(mixedPredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
				})
			}
		})
		b.Run("Into", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilterInto(nil, mixedPredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
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
				_ = SliceFilter(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilterInPlace(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
		b.Run("Into", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilterInto(nil, consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
	})
	b.Run("consecutive_middle", func(b *testing.B) {
		consecutivePredicate := func(v int) bool {
			return v >= 30 && v < 70
		}
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilter(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilterInPlace(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
		b.Run("Into", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilterInto(nil, consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
	})
	b.Run("consecutive_end", func(b *testing.B) {
		consecutivePredicate := func(v int) bool {
			return v >= 60
		}
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilter(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilterInPlace(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
		b.Run("Into", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SliceFilterInto(nil, consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
	})
}

func BenchmarkSliceSplits(b *testing.B) {
	b.Run("mixed", func(b *testing.B) {
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplit(func(v int) bool {
					return v%2 == 0
				}, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
				})
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlace(func(v int) bool {
					return v%2 == 0
				}, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
				})
			}
		})
		b.Run("InPlaceUnstable", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlaceUnstable(func(v int) bool {
					return v%2 == 0
				}, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
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
				_, _ = SliceSplit(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlace(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
		b.Run("InPlaceUnstable", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlaceUnstable(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
	})
	b.Run("consecutive_middle", func(b *testing.B) {
		consecutivePredicate := func(v int) bool {
			return v >= 30 && v < 70
		}
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplit(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlace(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
		b.Run("InPlaceUnstable", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlaceUnstable(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
	})
	b.Run("consecutive_end", func(b *testing.B) {
		consecutivePredicate := func(v int) bool {
			return v >= 60
		}
		b.Run("basic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplit(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
		b.Run("InPlace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlace(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
		b.Run("InPlaceUnstable", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = SliceSplitInPlaceUnstable(consecutivePredicate, []int{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
				})
			}
		})
	})
}

// -- TEST BELOW ARE BASED ON UNIT TEST INPUTS (VARY OVER TIME) --

func BenchmarkSliceFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceTestCases {
			_ = SliceFilter(tc.testFunc, tc.input)
		}
	}
}

func BenchmarkSliceFilterMultiple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceMultipleTestCases {
			_ = SliceFilter(tc.testFunc, tc.slices...)
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
			_ = SliceFilterInPlace(tc.testFunc, testInputs[j])
		}
	}
}

func BenchmarkSliceSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceTestCases {
			_, _ = SliceSplit(tc.testFunc, tc.input)
		}
	}
}

func BenchmarkSliceSplitMultiple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceMultipleTestCases {
			_, _ = SliceSplit(tc.testFunc, tc.slices...)
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
			_, _ = SliceSplitInPlace(tc.testFunc, testInputs[j])
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
			_, _ = SliceSplitInPlaceUnstable(tc.testFunc, testInputs[j])
		}
	}
}

func BenchmarkSliceToSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceToSetTests {
			_ = SliceToSet(tc.input...)
		}
	}
}

func BenchmarkSliceToSetBy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceToSetByTests {
			_ = SliceToSetBy(tc.conversion, tc.inputSlices...)
		}
	}
}

func BenchmarkSliceToCounts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceToCountsTests {
			_ = SliceToCounts(tc.input...)
		}
	}
}

func BenchmarkSliceToCountsBy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceToCountsByTests {
			_ = SliceToCountsBy(tc.conversion, tc.inputSlices...)
		}
	}
}

func BenchmarkSliceToIndexBy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceToIndexByTests {
			_ = SliceToIndexBy(tc.conversion, tc.inputSlices...)
		}
	}
}

func BenchmarkSliceToGroupsBy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceToGroupsByTests {
			_ = SliceToGroupsBy(tc.conversion, tc.inputSlices...)
		}
	}
}

func BenchmarkSliceIntersect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceSetOperationTests {
			_ = SliceIntersect(tc.sliceA, tc.sliceB)
		}
	}
}

func BenchmarkSliceDifference(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range sliceSetOperationTests {
			_ = SliceDifference(tc.sliceA, tc.sliceB)
		}
	}
}

func BenchmarkSliceConcat(b *testing.B) {
	b.Run("single slice", func(b *testing.B) {
		input := [][]int{sliceLargeInput}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sliceConcat(input, false)
		}
	})

	b.Run("multiple slices", func(b *testing.B) {
		input := [][]int{
			{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			{20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sliceConcat(input, false)
		}
	})

	b.Run("many small slices", func(b *testing.B) {
		input := make([][]int, 10)
		for i := range input {
			input[i] = []int{i, i + 1}
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sliceConcat(input, false)
		}
	})

	b.Run("empty slices", func(b *testing.B) {
		input := [][]int{{}, {}, {}}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sliceConcat(input, false)
		}
	})
}
