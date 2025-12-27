// The polyomino packing problem
package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type Polyomino struct {
	id    int     // unique identifier
	cells [][]int // 0/1 grid representation of the polyomino
}

func (p Polyomino) solidCellsCount() int {
	count := 0
	for r := 0; r < len(p.cells); r++ {
		for c := 0; c < len(p.cells[0]); c++ {
			if p.cells[r][c] == 1 {
				count++
			}
		}
	}
	return count
}

type GridPacking struct {
	width, height   int
	polyominoCounts []int // how many of each polyomino type to place
}

type ProblemSpace struct {
	polyominos   []Polyomino
	gridPackings []GridPacking
}

// needed for parsing input
func IsPolyominoSection(firstLine string) bool {
	polyominoIdPattern := regexp.MustCompile(`^\d+:$`)
	if polyominoIdPattern.MatchString(strings.TrimSpace(firstLine)) {
		return true
	}
	return false
}

func ParsePolyomino(section string) Polyomino {
	lines := strings.Split(strings.TrimSpace(section), "\n")
	// first line is ID
	var polyId int
	fmt.Sscanf(lines[0], "%d:", &polyId)
	cells := make([][]int, len(lines)-1)
	for i, line := range lines[1:] {
		line = strings.TrimSpace(line)
		cells[i] = make([]int, len(line))
		for j, char := range line {
			if char == '#' {
				cells[i][j] = 1
			} else {
				cells[i][j] = 0
			}
		}
	}

	return Polyomino{id: polyId, cells: cells}
}

func ParseGridPacking(line string) GridPacking {

	parts := strings.Split(line, ":")
	var width, height int
	fmt.Sscanf(strings.TrimSpace(parts[0]), "%dx%d", &width, &height)
	polyCounts := make([]int, 0)
	countParts := strings.Split(strings.TrimSpace(parts[1]), " ")
	for _, countStr := range countParts {
		var count int
		fmt.Sscanf(countStr, "%d", &count)
		polyCounts = append(polyCounts, count)
	}

	return GridPacking{width: width, height: height, polyominoCounts: polyCounts}
}

func ParseInput() ProblemSpace {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	problem := ProblemSpace{polyominos: make([]Polyomino, 0), gridPackings: make([]GridPacking, 0)}

	sections := strings.Split(string(data), "\n\n")

	for _, section := range sections {
		lines := strings.Split(strings.TrimSpace(section), "\n")
		firstLine := lines[0]
		if IsPolyominoSection(firstLine) {
			// Parse polyomino
			poly := ParsePolyomino(section)
			problem.polyominos = append(problem.polyominos, poly)
			continue
		} else {
			// Parse grid packing from each line
			for _, line := range lines {
				pack := ParseGridPacking(line)
				problem.gridPackings = append(problem.gridPackings, pack)
			}
		}
	}
	// TODO: parse data into []Polyomino
	return problem
}

func Part1_serious(problem ProblemSpace) {

	successCount := 0
	for i, gridPack := range problem.gridPackings {
		fmt.Printf("Problem %d: grid %dx%d, poly counts: %v\n", i, gridPack.width, gridPack.height, gridPack.polyominoCounts)

		// build the grid
		grid := make([][]int, gridPack.height)
		for r := 0; r < gridPack.height; r++ {
			grid[r] = make([]int, gridPack.width)
		}

		ok, _, _ := CanPackWithCounts(grid, problem.polyominos, gridPack.polyominoCounts)
		fmt.Printf("problem %d  Can pack: %v\n", i, ok)
		if ok {
			successCount++
		}
	}
	fmt.Printf("Total successful packings: %d\n", successCount)

}

func Part1_guestimation(problem ProblemSpace) {
	possiblePacking := 0
	magicFactor := 1.23
	for k, gridPack := range problem.gridPackings {

		gridSize := gridPack.width * gridPack.height
		fmt.Printf("grid %d, total gridSize: %d\n", k, gridSize)
		totalPolyCells := 0
		for i, count := range problem.gridPackings[k].polyominoCounts {
			totalPolyCells += problem.polyominos[i].solidCellsCount() * count
		}
		fmt.Printf("total solid poly cells %d\n", totalPolyCells)
		if int(float64(totalPolyCells)*magicFactor) <= gridSize {
			possiblePacking++
		}
	}

	fmt.Printf("Total possible packings (guestimate): %d\n", possiblePacking)
}

func main() {
	fmt.Println("--- Day 12: Christmas Tree Farm ---")
	problem := ParseInput()
	fmt.Printf("Parsed Problem Space:\n%+v\n", problem)
	Part1_guestimation(problem)
}
