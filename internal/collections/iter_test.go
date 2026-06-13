package collections

import (
	"slices"
	"testing"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:  "less than 3",
			input: []int{1, 2, 3, 4, 5},
			predicate: func(i int) bool {
				return i < 3
			},
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq := SeqValues(slices.All(tt.input))
			ft := slices.Collect(SeqFilter(seq, tt.predicate))
			if slices.Compare(ft, tt.expected) != 0 {
				t.Errorf("filter: expected %v, got %v", tt.expected, ft)
			}
		})
	}
}
