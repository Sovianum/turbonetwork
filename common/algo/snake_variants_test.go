package algo

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetSnakeInts(t *testing.T) {
	tc := []struct{
		start int
		end int
		len int
		expected []int
	}{
		{
			start:0,
			end:0,
			len:0,
			expected:[]int{},
		},
		{
			start:0,
			end:0,
			len:3,
			expected:[]int{0, 0, 0},
		},
		{
			start:0,
			end:2,
			len:8,
			expected:[]int{0, 1, 2, 2, 1, 0, 0, 1},
		},
	}

	for i, c := range tc {
		assert.Equal(t, c.expected, getSnakeInts(c.start, c.end, c.len), "%d", i)
	}
}

func TestGetSnakeVariants(t *testing.T) {
	tc := []struct{
		limits []int
		expected [][]int
	}{
		{
			limits:[]int{0, 0, 0},
			expected:[][]int{},
		},
		{
			limits:[]int{1, 1, 1},
			expected:[][]int{
				{0, 0, 0},
			},
		},
		{
			limits:[]int{3, 3},
			expected:[][]int{
				{0, 0},
				{1, 0},
				{2, 0},
				{2, 1},
				{1, 1},
				{0, 1},
				{0, 2},
				{1, 2},
				{2, 2},
			},
		},
		{
			limits:[]int{2, 2, 2},
			expected:[][]int{
				{0, 0, 0},
				{1, 0, 0},
				{1, 1, 0},
				{0, 1, 0},
				{0, 1, 1},
				{1, 1, 1},
				{1, 0, 1},
				{0, 0, 1},
			},
		},
	}

	for i, c := range tc {
		assert.Equal(t, c.expected, GetSnakeVariants(c.limits), "%d", i)
	}
}
