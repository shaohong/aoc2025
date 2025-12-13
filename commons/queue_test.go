package queue

import "testing"

type Position struct {
	row int
	col int
}

func TestQueueFIFO(t *testing.T) {
	var q Queue[Position]

	positions := []Position{
		{row: 0, col: 1},
		{row: 2, col: 3},
		{row: 4, col: 5},
	}

	for _, p := range positions {
		q.Enqueue(p)
	}

	if q.Len() != len(positions) {
		t.Fatalf("expected length %d, got %d", len(positions), q.Len())
	}

	for i, expected := range positions {
		item, ok := q.Dequeue()
		if !ok {
			t.Fatalf("expected dequeue to succeed at index %d", i)
		}
		if item != expected {
			t.Fatalf("expected %+v, got %+v", expected, item)
		}
	}

	if _, ok := q.Dequeue(); ok {
		t.Fatalf("expected dequeue on empty queue to fail")
	}
	if q.Len() != 0 {
		t.Fatalf("expected length 0 after draining queue, got %d", q.Len())
	}
}
