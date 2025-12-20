// https://adventofcode.com/2025/day/9
package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
	"strings"
)

// a tile is indicated by its x,y coordinate
type Tile struct {
	x int
	y int
}

type TilePair struct {
	tileA Tile
	tileB Tile
	area  int // area of the rectangle formed by the two tiles
}

type Interval struct {
	start int
	end   int
}

// merge intervales into non-overlapping ones
func MergeIntervals(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return intervals
	}

	// sort intervals by start
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].start < intervals[j].start
	})

	nonOverlappingIntervals := make([]Interval, 0)
	for _, interval := range intervals {
		if len(nonOverlappingIntervals) == 0 {
			nonOverlappingIntervals = append(nonOverlappingIntervals, interval)
			continue
		}
		lastInterval := nonOverlappingIntervals[len(nonOverlappingIntervals)-1]
		if interval.start <= lastInterval.end {
			// merge intervals
			nonOverlappingIntervals[len(nonOverlappingIntervals)-1].end = slices.Max([]int{lastInterval.end, interval.end})
		} else {
			nonOverlappingIntervals = append(nonOverlappingIntervals, interval)
		}
	}

	return nonOverlappingIntervals
}

func ParseInput() []Tile {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(input), "\n")
	tiles := make([]Tile, 0)
	for _, line := range lines {
		if line == "" {
			continue
		}
		newTile := Tile{}
		fmt.Sscanf(line, "%d,%d", &newTile.x, &newTile.y)
		tiles = append(tiles, newTile)
	}
	return tiles
}

// the generalized distance function between two tiles
func Distance(a, b Tile) int {
	var dx, dy int
	if a.x > b.x {
		dx = a.x - b.x + 1
	} else {
		dx = b.x - a.x + 1
	}
	if a.y > b.y {
		dy = a.y - b.y + 1
	} else {
		dy = b.y - a.y + 1
	}

	return dx * dy
}

func Part1() {
	tiles := ParseInput()

	fmt.Println("Parsed Tiles:", tiles)

	// compute all pairwise distances
	largestDistance := 0
	for i := 0; i < len(tiles); i++ {
		for j := i + 1; j < len(tiles); j++ {
			area := Distance(tiles[i], tiles[j])
			if area > largestDistance {
				largestDistance = area
			}
		}
	}

	fmt.Println("Largest Distance Area:", largestDistance)
}

func PairWithinBoundaries(pair TilePair, verticalSegments [][2]Tile, horizontalSegments [][2]Tile) bool {

	minX := slices.Min([]int{pair.tileA.x, pair.tileB.x})
	maxX := slices.Max([]int{pair.tileA.x, pair.tileB.x})
	minY := slices.Min([]int{pair.tileA.y, pair.tileB.y})
	maxY := slices.Max([]int{pair.tileA.y, pair.tileB.y})

	// helper checks whether filtered vertical segments fully span the Y interval
	covered := func(filter func([2]Tile) bool) bool {
		intervals := make([]Interval, 0)
		for _, segment := range verticalSegments {
			if !filter(segment) {
				continue
			}
			startY := slices.Min([]int{segment[0].y, segment[1].y})
			endY := slices.Max([]int{segment[0].y, segment[1].y})
			intervals = append(intervals, Interval{start: startY, end: endY})
		}
		for _, interval := range MergeIntervals(intervals) {
			if interval.start <= minY && interval.end >= maxY {
				return true
			}
		}
		return false
	}

	leftCovered := covered(func(segment [2]Tile) bool {
		return segment[0].x <= minX
	})
	if !leftCovered {
		return false
	}

	fmt.Printf("the left side of the pair rectangle(%v -> %v) is covered within the boundary\n", Tile{minX, minY}, Tile{minX, maxY})

	rightCovered := covered(func(segment [2]Tile) bool {
		return segment[0].x >= maxX
	})
	if !rightCovered {
		return false
	}

	fmt.Printf("the right side of the pair rectangle(%v -> %v) is covered within the boundary\n", Tile{maxX, minY}, Tile{maxX, maxY})

	return true
}

func Part2() {
	tiles := ParseInput()

	// find all the line segments between two consecutive tiles
	// they serve as the boundaries of the polygon/area formed by the tiles

	verticalSegments := make([][2]Tile, 0)   // each segment is represented by two tiles having the same x,
	horizontalSegments := make([][2]Tile, 0) // each segment is represented by two tiles having the same y,
	for i := 0; i < len(tiles); i++ {
		tile_1 := tiles[i]
		tile_2 := tiles[(i+1)%len(tiles)]

		if tile_1.x == tile_2.x {
			verticalSegments = append(verticalSegments, [2]Tile{tile_1, tile_2})
		}
		if tile_1.y == tile_2.y {
			horizontalSegments = append(horizontalSegments, [2]Tile{tile_1, tile_2})
		}
	}

	tilePairs := make([]TilePair, 0)
	for i := 0; i < len(tiles); i++ {
		for j := i + 1; j < len(tiles); j++ {
			area := Distance(tiles[i], tiles[j])
			tilePairs = append(tilePairs, TilePair{tileA: tiles[i], tileB: tiles[j], area: area})
		}
	}

	// sort TilePairs by area descending
	sort.Slice(tilePairs, func(i, j int) bool {
		return tilePairs[i].area > tilePairs[j].area
	})

	largestInternalArea := 0
	for _, pair := range tilePairs {
		// check if the rectangle formed by pair.tileA and pair.tileB is within the boundaries (defined by the segments)
		if PairWithinBoundaries(pair, verticalSegments, horizontalSegments) {
			largestInternalArea = pair.area
			break
		}
	}

	fmt.Println("Largest Internal Area:", largestInternalArea)
}

func main() {
	fmt.Println("--- Day 9: Movie Theater ---")
	// Part1()
	Part2()

}
