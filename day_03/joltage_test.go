package main

import (
	"testing"
)

func TestLargestTwoDigitNumber(t *testing.T) {
	tests := []struct {
		joltages string
		expected int
	}{
		{"987654321111111", 98},
		{"811111111111119", 89},
		{"234234234234278", 78},
		{"818181911112111", 92},
	}

	for _, test := range tests {
		result := LargestTwoDigitNumber(test.joltages)
		if result != test.expected {
			t.Errorf("For joltages %s, expected %d but got %d", test.joltages, test.expected, result)
		}
	}
}

func TestLargestNDigitNumber(t *testing.T) {
	tests := []struct {
		joltages string
		ndigits  int
		expected int
	}{
		{"987654321111111", 3, 987},
		{"811111111111119", 4, 8119},
		{"234234234234278", 2, 78},
		{"234234234234278", 12, 434234234278},
	}

	for _, test := range tests {
		result := LargestNDigitNumber(test.joltages, test.ndigits)
		if result != test.expected {
			t.Errorf("For joltages %s and ndigits %d, expected %d but got %d", test.joltages, test.ndigits, test.expected, result)
		}
	}
}
