// https://adventofcode.com/2025/day/5
package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type IdRange struct {
	start int
	end   int
}

func (r IdRange) Contains(id int) bool {
	return id >= r.start && id <= r.end
}

type DataBase []IdRange

func (db DataBase) HasID(id int) bool {
	for _, r := range db {
		if r.Contains(id) {
			return true
		}
	}
	return false
}

func ParseIDRanges(input string) []IdRange {
	idRanges := make([]IdRange, 0)
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		var start, end int
		fmt.Sscanf(line, "%d-%d", &start, &end)
		idRanges = append(idRanges, IdRange{start: start, end: end})
	}

	return idRanges
}

func ParseIntegerList(input string) []int {

	intList := make([]int, 0)
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		var value int
		fmt.Sscanf(line, "%d", &value)
		intList = append(intList, value)
	}

	return intList
}

func ParseInput(r io.Reader) ([]IdRange, []int) {
	input, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	// split the input by a blank line
	sections := strings.Split(string(input), "\n\n")

	idRanges := ParseIDRanges(sections[0])
	availablIngredients := ParseIntegerList(sections[1])

	return idRanges, availablIngredients

}

func Part1() {
	idRanges, availableIngredients := ParseInput(os.Stdin)

	dataBase := DataBase(idRanges)

	freshIngredientCount := 0
	for _, ingredientID := range availableIngredients {
		if dataBase.HasID(ingredientID) {
			freshIngredientCount++
		}
	}

	fmt.Println("Number of fresh ingredients available:", freshIngredientCount)
}

func Part2() {
	// find how many total effective ingredient IDs are available
	idRanges, _ := ParseInput(os.Stdin)

	// we need to essentially merge the idRanges to find the total coverage
	// reference: https://stackoverflow.com/questions/32585990/algorithm-merge-overlapping-segments

	sort.Slice(idRanges, func(i, j int) bool {
		return idRanges[i].start < idRanges[j].start
	})

	effectiveRanges := make([]IdRange, 0)
	for _, idRange := range idRanges {
		// if the start bigger than the end of the last range, add a new range
		if len(effectiveRanges) == 0 || idRange.start > effectiveRanges[len(effectiveRanges)-1].end {
			effectiveRanges = append(effectiveRanges, idRange)
		} else {
			// otherwise, merge the ranges
			lastRange := &effectiveRanges[len(effectiveRanges)-1]
			if idRange.end > lastRange.end {
				lastRange.end = idRange.end
			}
		}
	}

	totalEffectiveIDs := 0
	for _, r := range effectiveRanges {
		totalEffectiveIDs += (r.end - r.start + 1)
	}

	fmt.Println("Total effective ingredient IDs available:", totalEffectiveIDs)
}

func main() {
	fmt.Println("Day 5: Cafeteria")

	// Part1()
	Part2()

}
