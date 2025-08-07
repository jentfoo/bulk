package bulk

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
