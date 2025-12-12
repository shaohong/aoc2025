package main

import (
	"testing"
)

func TestDatabaseHasID(t *testing.T) {

	db := DataBase{
		IdRange{start: 3, end: 5},
		IdRange{start: 10, end: 14},
		IdRange{start: 16, end: 20},
		IdRange{start: 12, end: 18},
	}

	tests := []struct {
		id       int
		expected bool
	}{
		{id: 1, expected: false},
		{id: 5, expected: true},
		{id: 8, expected: false},
		{id: 11, expected: true},
		{id: 17, expected: true},
		{id: 32, expected: false},
	}

	for _, test := range tests {
		result := db.HasID(test.id)
		if result != test.expected {
			t.Errorf("HasID(%d) = %v; want %v", test.id, result, test.expected)
		}
	}

}
