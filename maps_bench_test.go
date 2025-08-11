package bulk

import (
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
