package main

import "testing"

func TestIsRepeatingSequence(t *testing.T) {
	tests := []struct {
		productID string
		expected  bool
	}{
		{"abab", true},
		{"abcabc", true},
		{"abcd", false},
		{"aaaa", true},
	}

	for _, test := range tests {
		result := IsRepeatingSequence(test.productID)
		if result != test.expected {
			t.Errorf("For productID %s, expected %v but got %v", test.productID, test.expected, result)
		}
	}
}
