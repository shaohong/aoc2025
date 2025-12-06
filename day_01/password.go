package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Instruction struct {
	direction string
	steps     int
}

type Dial struct {
	size     int
	position int
}

func (d *Dial) Move(moveInstr Instruction) {
	if moveInstr.direction == "R" {
		d.position += (moveInstr.steps % d.size)
	} else if moveInstr.direction == "L" {
		d.position -= (moveInstr.steps % d.size)
	}

	if d.position >= d.size {
		d.position -= d.size
	}
	if d.position < 0 {
		d.position += d.size
	}
}

// Move the dial according to the instruction, and return how many times it passed position 0, inccluding stopping at 0.
func (d *Dial) MovePastZero(moveInstr Instruction) int {
	countZero := 0
	if moveInstr.direction == "R" {
		newPosition := d.position + moveInstr.steps
		countZero += newPosition / d.size
		d.position = newPosition % d.size
	} else if moveInstr.direction == "L" {
		newPosition := d.position - moveInstr.steps

		for newPosition < 0 {
			// if we start from 0, then we don't count passing 0
			if d.position != 0 {
				countZero++
			}
			newPosition += d.size
			d.position = newPosition
		}

		d.position = newPosition

		// if we stop at 0, count it
		if newPosition == 0 {
			countZero++
		}
	}

	fmt.Println("Move", moveInstr, "newPosition", d.position, "countZero", countZero)

	return countZero
}

func ParseInstructions(input string) []Instruction {
	// parse lines into instructions
	instructions := []Instruction{}

	for _, line := range strings.Split(input, "\n") {
		if strings.TrimSpace(line) == "" {
			break
		}
		// the instruction string is like "L2" or "R50"
		// the first character is the direction, the rest is the number of steps
		dir := string(line[0])
		var steps int
		fmt.Sscanf(line[1:], "%d", &steps)
		instructions = append(instructions, Instruction{direction: dir, steps: steps})
	}

	return instructions
}

func Part1() {
	// read stdin to a buffer
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	instructions := ParseInstructions(string(data))

	// create a dial of size 100
	dial := Dial{size: 100, position: 50}

	countZero := 0
	for _, instr := range instructions {
		dial.Move(instr)
		if dial.position == 0 {
			countZero++
		}
	}
	fmt.Println("count of 0 position:", countZero)
}

func Part2() {
	// read stdin to a data buffer
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	instructions := ParseInstructions(string(data))

	// create a dial of size 100
	dial := Dial{size: 100, position: 50}

	countZero := 0
	for _, instr := range instructions {
		countZero += dial.MovePastZero(instr)
	}

	fmt.Println("count of passing 0 position:", countZero)
}

func main() {
	//Part1()

	Part2()
}
