package commons

import "testing"

func TestStack(t *testing.T) {
	var s Stack[Position]

	positions := []Position{
		{row: 0, col: 1},
		{row: 2, col: 3},
		{row: 4, col: 5},
	}

	if !s.IsEmpty() {
		t.Fatalf("expected stack to be empty")
	}

	for _, p := range positions {
		s.Push(p)
	}

	if s.Len() != len(positions) {
		t.Fatalf("expected stack length to be %d, got %d", len(positions), s.Len())
	}

	if s.IsEmpty() {
		t.Fatalf("expected stack to not be empty")
	}

	for i := len(positions) - 1; i >= 0; i-- {
		expected := positions[i]
		item, ok := s.Pop()
		if !ok {
			t.Fatalf("expected pop to succeed at index %d", i)
		}
		if item != expected {
			t.Fatalf("expected %+v, got %+v", expected, item)
		}
	}

	if !s.IsEmpty() {
		t.Fatalf("expected stack to be empty after popping all items")
	}

	if _, ok := s.Pop(); ok {
		t.Fatalf("expected pop on empty stack to fail")
	}
}
