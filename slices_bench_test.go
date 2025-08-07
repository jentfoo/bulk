package bulk

import (
	"testing"
)

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
	testInputs := make([][]int, len(sliceTestCases))
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
