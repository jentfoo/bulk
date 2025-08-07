package bulk

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mapMergeTestCases = []struct {
	name   string
	inputs []map[string]int
	expect map[string]int
}{
	{
		name:   "nil",
		inputs: nil,
		expect: map[string]int{},
	},
	{
		name:   "empty",
		inputs: []map[string]int{},
		expect: map[string]int{},
	},
	{
		name:   "single_map",
		inputs: []map[string]int{{"a": 1, "b": 2}},
		expect: map[string]int{"a": 1, "b": 2},
	},
	{
		name:   "single_empty_map",
		inputs: []map[string]int{{}},
		expect: map[string]int{},
	},
	{
		name:   "single_nil_map",
		inputs: []map[string]int{nil},
		expect: nil,
	},
	{
		name:   "two_maps_no_overlap",
		inputs: []map[string]int{{"a": 1, "b": 2}, {"c": 3, "d": 4}},
		expect: map[string]int{"a": 1, "b": 2, "c": 3, "d": 4},
	},
	{
		name:   "two_maps_with_overlap_last_wins",
		inputs: []map[string]int{{"a": 1, "b": 2}, {"b": 20, "c": 3}},
		expect: map[string]int{"a": 1, "b": 20, "c": 3},
	},
	{
		name:   "three_maps_with_overlap",
		inputs: []map[string]int{{"a": 1, "b": 2}, {"b": 20, "c": 3}, {"c": 30, "d": 4}},
		expect: map[string]int{"a": 1, "b": 20, "c": 30, "d": 4},
	},
	{
		name:   "multiple_maps_same_key_last_wins",
		inputs: []map[string]int{{"x": 1}, {"x": 2}, {"x": 3}, {"x": 4}},
		expect: map[string]int{"x": 4},
	},
	{
		name:   "mixed_with_empty_maps",
		inputs: []map[string]int{{"a": 1}, {}, {"b": 2}, {}, {"c": 3}},
		expect: map[string]int{"a": 1, "b": 2, "c": 3},
	},
	{
		name:   "mixed_with_nil_maps",
		inputs: []map[string]int{{"a": 1}, nil, {"b": 2}, nil, {"c": 3}},
		expect: map[string]int{"a": 1, "b": 2, "c": 3},
	},
	{
		name: "large_merge",
		inputs: []map[string]int{
			{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5},
			{"f": 6, "g": 7, "h": 8, "i": 9, "j": 10},
			{"k": 11, "l": 12, "m": 13, "n": 14, "o": 15},
			{"p": 16, "q": 17, "r": 18, "s": 19, "t": 20},
		},
		expect: map[string]int{
			"a": 1, "b": 2, "c": 3, "d": 4, "e": 5,
			"f": 6, "g": 7, "h": 8, "i": 9, "j": 10,
			"k": 11, "l": 12, "m": 13, "n": 14, "o": 15,
			"p": 16, "q": 17, "r": 18, "s": 19, "t": 20,
		},
	},
	{
		name: "complex_overlap_pattern",
		inputs: []map[string]int{
			{"a": 1, "b": 2, "c": 3},
			{"b": 20, "c": 30, "d": 4},
			{"a": 100, "d": 40, "e": 5},
		},
		expect: map[string]int{"a": 100, "b": 20, "c": 30, "d": 40, "e": 5},
	},
}

func TestMapMerge(t *testing.T) {
	t.Parallel()

	for i, tt := range mapMergeTestCases {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			got := MapUnion(tt.inputs...)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestMapKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  map[string]int
		expect []string
	}{
		{
			name:   "basic",
			input:  map[string]int{"a": 1, "b": 2, "c": 3},
			expect: []string{"a", "b", "c"},
		},
		{
			name:   "empty_map",
			input:  map[string]int{},
			expect: []string{},
		},
		{
			name:   "nil",
			input:  nil,
			expect: []string{},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			got := mapKeys(tt.input)

			assert.Len(t, got, len(tt.expect))
			assert.ElementsMatch(t, tt.expect, got)
		})
	}
}
