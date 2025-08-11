package bulk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var sliceFuzzSeeds = [][]byte{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
	{88, 84, 82, 80, 64, 32, 16, 8, 4, 2, 0},
	{1, 3, 5, 7, 9, 11, 13, 15},
	{2, 4, 8, 16, 1, 3, 5, 7},
	{1, 3, 5, 7, 2, 4, 8, 16},
	{1, 3, 5, 7, 2, 4, 8, 16, 9, 11, 13, 15},
	{2, 4, 8, 16, 1, 3, 5, 7, 32, 64, 80, 128},
}

func verifyByteSliceMod(t *testing.T, bytes []byte, mod, expect byte) {
	t.Helper()

	for i := range bytes {
		assert.Equal(t, expect, bytes[i]%mod)
	}
}

func FuzzSliceFilter(f *testing.F) {
	for _, seed := range sliceFuzzSeeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input []byte) {
		result := SliceFilter(func(b byte) bool {
			return b%2 == 0
		}, input)
		verifyByteSliceMod(t, result, 2, 0)
	})
}

func FuzzSliceFilterInPlace(f *testing.F) {
	for _, seed := range sliceFuzzSeeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input []byte) {
		result := SliceFilterInPlace(func(i byte) bool {
			return i%2 == 0
		}, input)
		verifyByteSliceMod(t, result, 2, 0)
	})
}

func FuzzSliceSplit(f *testing.F) {
	for _, seed := range sliceFuzzSeeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input []byte) {
		even, odd := SliceSplit(func(i byte) bool {
			return i%2 == 0
		}, input)
		verifyByteSliceMod(t, even, 2, 0)
		verifyByteSliceMod(t, odd, 2, 1)
	})
}

func FuzzSliceSplitInPlace(f *testing.F) {
	for _, seed := range sliceFuzzSeeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input []byte) {
		even, odd := SliceSplitInPlace(func(b byte) bool {
			return b%2 == 0
		}, input)
		verifyByteSliceMod(t, even, 2, 0)
		verifyByteSliceMod(t, odd, 2, 1)
	})
}

func FuzzSliceSplitInPlaceUnstable(f *testing.F) {
	for _, seed := range sliceFuzzSeeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input []byte) {
		even, odd := SliceSplitInPlaceUnstable(func(b byte) bool {
			return b%2 == 0
		}, input)
		verifyByteSliceMod(t, even, 2, 0)
		verifyByteSliceMod(t, odd, 2, 1)
	})
}
