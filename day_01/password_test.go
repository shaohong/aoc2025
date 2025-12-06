package main

import (
	"testing"
)

func TestDial(t *testing.T) {
	// create a dial of size 10
	dial := Dial{size: 10, position: 0}

	// move right 3 steps
	dial.Move(Instruction{direction: "R", steps: 3})
	if dial.position != 3 {
		t.Errorf("Expected position 3, got %d", dial.position)
	}

	// move left 5 steps
	dial.Move(Instruction{direction: "L", steps: 5})
	if dial.position != 8 {
		t.Errorf("Expected position 8, got %d", dial.position)
	}
}

func TestDialMovePastZero(t *testing.T) {
	dial := Dial{size: 10, position: 0}

	count := dial.MovePastZero(Instruction{direction: "R", steps: 25})
	if dial.position != 5 || count != 2 {
		t.Errorf("Expected position 5 and count 2, got position %d and count %d", dial.position, count)
	}

	dial.position = 0

	count = dial.MovePastZero(Instruction{direction: "L", steps: 15})
	if dial.position != 5 || count != 1 {
		t.Errorf("Expected position 5 and count 1, got position %d and count %d", dial.position, count)
	}

	dial.position = 3
	count = dial.MovePastZero(Instruction{direction: "L", steps: 4})
	if dial.position != 9 || count != 1 {
		t.Errorf("Expected position 9 and count 1, got position %d and count %d", dial.position, count)
	}

}
