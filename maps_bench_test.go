package bulk

import (
	"maps"
	"slices"
	"strconv"
	"testing"
)

func BenchmarkMapInvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range mapInvertTests {
			_ = MapInvert(tc.input)
		}
	}
}

func BenchmarkMapInvertInto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range mapInvertTests {
			result := make(map[string]int)
			MapInvertInto(result, tc.input)
		}
	}
}

func BenchmarkMapKeysSlice(b *testing.B) {
	smallTestMaps := make([]map[string]struct{}, 0, 20)
	medTestMaps := make([]map[string]struct{}, 0, 20)
	largeTestMaps := make([]map[string]struct{}, 0, 20)
	for _, count := range []int{1, 2, 4, 5, 6, 7, 8, 9,
		10, 20, 30, 40, 50, 60, 70, 80, 90,
	} {
		m := make(map[string]struct{}, count)
		for i := 0; i < count; i++ {
			m[strconv.Itoa(i)] = struct{}{}
		}
		smallTestMaps = append(smallTestMaps, m)
	}
	for _, count := range []int{100, 200, 300, 400, 500, 600, 700, 800, 900,
		1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000,
	} {
		m := make(map[string]struct{}, count)
		for i := 0; i < count; i++ {
			m[strconv.Itoa(i)] = struct{}{}
		}
		medTestMaps = append(medTestMaps, m)
	}
	for _, count := range []int{
		10_000, 20_000, 30_000, 40_000, 50_000, 60_000, 70_000, 80_000, 90_000,
		100_000, 200_000, 300_000, 400_000, 500_000, 600_000, 700_000, 800_000, 900_000,
	} {
		m := make(map[string]struct{}, count)
		for i := 0; i < count; i++ {
			m[strconv.Itoa(i)] = struct{}{}
		}
		largeTestMaps = append(largeTestMaps, m)
	}
	b.ResetTimer()

	b.Run("small-slices", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, m := range smallTestMaps {
				_ = slices.Collect(maps.Keys(m))
			}
		}
	})
	b.Run("small-bulk", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, m := range smallTestMaps {
				_ = MapKeysSlice(m)
			}
		}
	})
	b.Run("med-slices", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, m := range medTestMaps {
				_ = slices.Collect(maps.Keys(m))
			}
		}
	})
	b.Run("med-bulk", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, m := range medTestMaps {
				_ = MapKeysSlice(m)
			}
		}
	})
	b.Run("large-slices", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, m := range largeTestMaps {
				_ = slices.Collect(maps.Keys(m))
			}
		}
	})
	b.Run("large-bulk", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, m := range largeTestMaps {
				_ = MapKeysSlice(m)
			}
		}
	})
}
