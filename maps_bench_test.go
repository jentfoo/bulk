package bulk

import (
	"testing"
)

func BenchmarkMapMerge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range mapMergeTestCases {
			_ = MapUnion(tc.inputs...)
		}
	}
}
