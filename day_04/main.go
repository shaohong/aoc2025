package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Grid struct {
	nRows int
	nCols int
	data  [][]byte
}

type Position struct {
	row int
	col int
}

const rollMarker byte = '@'

func (g *Grid) Get(row, col int) byte {
	return g.data[row][col]
}

// Remove the roll of paper at (row, col)
func (g *Grid) RemovePaperRoll(row, col int) {
	g.data[row][col] = '.'
}

// Find Liftable Rolls
func (g *Grid) FindLiftableRolls() []Position {
	liftable := make([]Position, 0)
	for row := 0; row < g.nRows; row++ {
		for col := 0; col < g.nCols; col++ {
			if g.CanBeForkLifted(row, col) {
				liftable = append(liftable, Position{row: row, col: col})
			}
		}
	}
	return liftable
}

func (g *Grid) CanBeForkLifted(row, col int) bool {
	// The forklifts can only access a roll of paper if there are fewer than four rolls of paper in the eight adjacent positions
	// go through the 8 adjacent positions and count how many values equal to '@'

	// don't bother if the current position is not a roll of paper
	if g.Get(row, col) != rollMarker {
		return false
	}

	neighboringRolls := 0
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			r := row + dr
			c := col + dc
			if r >= 0 && r < g.nRows && c >= 0 && c < g.nCols {
				if g.Get(r, c) == rollMarker {
					neighboringRolls++
				}
			}
		}
	}
	return neighboringRolls < 4
}

func ParseInput(r io.Reader) Grid {
	data, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")
	grid := make([][]byte, len(lines))
	for i, line := range lines {
		grid[i] = []byte(line)
	}
	return Grid{
		nRows: len(grid),
		nCols: len(grid[0]),
		data:  grid,
	}
}

func Part1() {

	grid := ParseInput(os.Stdin)

	fmt.Println(grid)
	movableRolls := 0
	for row := 0; row < grid.nRows; row++ {
		for col := 0; col < grid.nCols; col++ {
			if grid.CanBeForkLifted(row, col) {
				movableRolls++
			}
		}
	}

	fmt.Println("Number of rolls that can be forklifted:", movableRolls)
}

func Part2() {

	grid := ParseInput(os.Stdin)
	totalMovableRolls := 0
	for {
		movableRolls := grid.FindLiftableRolls()
		if len(movableRolls) == 0 {
			break
		}
		totalMovableRolls += len(movableRolls)
		for _, pos := range movableRolls {
			grid.RemovePaperRoll(pos.row, pos.col)
		}
	}
	fmt.Println("Total number of rolls that can be forklifted:", totalMovableRolls)
}

func main() {
	//Part1()
	Part2()
}
