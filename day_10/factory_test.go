package main

import (
	"testing"
)

func TestPowerSet(t *testing.T) {
	input := []int{1, 2, 3}
	expected := [][]int{
		{},
		{1},
		{2},
		{1, 2},
		{3},
		{1, 3},
		{2, 3},
		{1, 2, 3},
	}

	result := PowerSet(input)
	if len(result) != len(expected) {
		t.Errorf("Expected %d subsets, but got %d", len(expected), len(result))
	}

	result = PowerSetViaBacktracking(input)
	if len(result) != len(expected) {
		t.Errorf("Expected %d subsets, but got %d", len(expected), len(result))
	}
}
