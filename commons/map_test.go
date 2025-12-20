package commons

import "testing"

func TestMap(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}

	result := Map(input, func(x int) int {
		return x * 2
	})

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d but got %d at index %d", expected[i], v, i)
		}
	}
}
