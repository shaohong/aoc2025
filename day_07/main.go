// https://adventofcode.com/2025/day/7
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	queue "github.com/shaohong/aoc2025/commons"
	stack "github.com/shaohong/aoc2025/commons"
)

const splitterChar byte = '^'
const startChar byte = 'S'

type Lab struct {
	Rows [][]byte
}

type Position struct {
	row int
	col int
}

func (lab *Lab) FindStart() (Position, error) {
	// just assume start is always at row 0
	for col, char := range lab.Rows[0] {
		if byte(char) == startChar {
			fmt.Printf("Found start at (0,%d)\n", col)
			return Position{row: 0, col: col}, nil
		}
	}
	return Position{}, fmt.Errorf("start position not found")
}

func ParseInput() Lab {
	// read stdin into a two dimensional array of bytes
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")
	lab := Lab{Rows: make([][]byte, 0)}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			break
		}
		row := []byte(line)
		lab.Rows = append(lab.Rows, row)
	}

	// fmt.Println("Lab layout:")
	// for _, row := range lab.Rows {
	// 	fmt.Println(string(row))
	// }
	return lab
}

func Part1() {
	lab := ParseInput()

	visitedSpliters := make(map[Position]bool)
	visitedPositions := make(map[Position]bool)

	startPos, _ := lab.FindStart()
	q := queue.Queue[Position]{}
	q.Enqueue(startPos)

	for q.Len() > 0 {
		v, _ := q.Dequeue()
		visitedPositions[v] = true

		// from position v, travel downward in the lab, until hitting a splitter or the bottom
		row := v.row + 1
		col := v.col
		for row < len(lab.Rows) {
			char := lab.Rows[row][col]
			// mark position as visited
			visitedPositions[Position{row: row, col: col}] = true

			if char == splitterChar {
				// enqueue left and right positions, if they have not been visited before
				if col > 0 {
					candidatePos := Position{row: row, col: col - 1}
					if !visitedPositions[candidatePos] && !q.Contains(candidatePos) {
						q.Enqueue(candidatePos)
					}
				}

				if col < len(lab.Rows[row])-1 {
					candidatePos := Position{row: row, col: col + 1}
					if !visitedPositions[candidatePos] && !q.Contains(candidatePos) {
						q.Enqueue(Position{row: row, col: col + 1})
					}
				}

				// mark splitter as visited
				visitedSpliters[Position{row: row, col: col}] = true
				break

			} else if char == '.' || char == startChar {
				// continue downward
				row++
			}
		}

	}

	// print number of visited splitters
	fmt.Printf("Part 1: Number of visited splitters: %d\n", len(visitedSpliters))

}

func Part2_old() {

	lab := ParseInput()
	startPos, _ := lab.FindStart()
	fmt.Printf("Part 2: Start position is at (%d,%d)\n", startPos.row, startPos.col)

	stack := stack.Stack[Position]{}
	stack.Push(startPos)

	totalPaths := 0
	totalRows := len(lab.Rows)

	// start from the start position, going downwards, if a splitter is hit, push left and right position, done
	// if bottom is hit, count one path
	for !stack.IsEmpty() {
		// if stack.Len()%10 == 0 {
		// 	fmt.Printf("Part 2: Stack size: %d, total paths so far: %d\n", stack.Len(), totalPaths)
		// }
		if totalPaths%1000 == 0 {
			fmt.Printf("Part 2: Total distinct paths so far: %d\n", totalPaths)
		}

		v, _ := stack.Pop()
		row := v.row
		col := v.col

		for row < totalRows && lab.Rows[row][col] != splitterChar {
			row++
		}

		// why do I stop, did I hit bottom or a splitter?
		if row == len(lab.Rows) {
			totalPaths++
			continue
		}

		// push right and left positions
		if col < len(lab.Rows[row])-1 {
			candidatePos := Position{row: row, col: col + 1}
			stack.Push(candidatePos)
		}
		if col > 0 {
			candidatePos := Position{row: row, col: col - 1}
			stack.Push(candidatePos)
		}
	}

	fmt.Println("Part 2: Total distinct paths to the bottom:", totalPaths)

}

var pathCountMemo map[Position]int = make(map[Position]int)

// count the number of paths from the given position to the bottom of the lab
func CountPaths(lab *Lab, pos Position) (nPaths int) {
	// check if we can go all the way to the bottom from the current position

	if _, ok := pathCountMemo[pos]; ok {
		// log.Printf("already known paths from position %+v: %d\n", pos, pathCountMemo[pos])
		return pathCountMemo[pos]
	}

	// go down until hitting a splitter or the bottom
	row := pos.row
	col := pos.col
	totalRows := len(lab.Rows)
	for row < totalRows && lab.Rows[row][col] != splitterChar {
		row++
	}

	// did we hit the bottom?
	if row == totalRows {
		// log.Printf("position %+v can hit bottom directly \n", pos)
		nPaths = 1
		pathCountMemo[pos] = 1
		return nPaths
	}

	// we hit a splitter at (row, col), explore left and right
	if col > 0 {
		leftPos := Position{row: row, col: col - 1}
		nPaths += CountPaths(lab, leftPos)
	}
	if col < len(lab.Rows[row])-1 {
		rightPos := Position{row: row, col: col + 1}
		nPaths += CountPaths(lab, rightPos)
	}

	pathCountMemo[pos] = nPaths
	log.Printf("paths from position (%d,%d): %d\n", pos.row, pos.col, nPaths)
	return nPaths
}

func Part2() {
	lab := ParseInput()
	startPos, _ := lab.FindStart()
	fmt.Printf("Part 2: Start position is at (%d,%d)\n", startPos.row, startPos.col)
	totalPaths := CountPaths(&lab, startPos)
	fmt.Printf("Part 2: Total distinct paths to the bottom: %d\n", totalPaths)
}

func main() {
	fmt.Println("--- Day 7: Laboratories ---")
	//Part1()
	Part2()
}
