package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/draffensperger/golp"
	commons "github.com/shaohong/aoc2025/commons"
)

// generate all subsets (the power set) of a given set of integers
func PowerSet(set []int) [][]int {
	var result [][]int
	for i := 0; i < (1 << len(set)); i++ {
		var subset []int
		for j := 0; j < len(set); j++ {
			if (i & (1 << j)) != 0 {
				subset = append(subset, set[j])
			}
		}
		result = append(result, subset)
	}

	// sort result subsets by length
	sort.Slice(result, func(i, j int) bool {
		return len(result[i]) < len(result[j])
	})
	return result
}

func PowerSetViaBacktracking(set []int) [][]int {
	var result [][]int
	var backtrack func(start int, current []int)
	backtrack = func(start int, current []int) {
		// append a copy of current subset to result
		subset := make([]int, len(current))
		copy(subset, current)
		result = append(result, subset)

		for i := start; i < len(set); i++ {
			current = append(current, set[i])
			backtrack(i+1, current)
			current = current[:len(current)-1] // backtrack
		}
	}
	backtrack(0, []int{})

	sort.Slice(result, func(i, j int) bool {
		return len(result[i]) < len(result[j])
	})
	return result
}

type Machine struct {
	desired_light_status []int
	button_wires         [][]int
	jotage               []int
}

func (m Machine) String() string {
	return fmt.Sprintf("Desired Light Status: %v, Button Wires: %v, Jotage: %v",
		m.desired_light_status, m.button_wires, m.jotage)
}

// calculate how many button presses are needed to achieve the desired light status
func (m Machine) ToggleLights() int {
	numOfButtons := len(m.button_wires)

	// generate the power set of button indices
	buttonIndices := make([]int, numOfButtons)
	for i := 0; i < numOfButtons; i++ {
		buttonIndices[i] = i
	}
	subsets := PowerSet(buttonIndices)

	// try each subset of button presses
	for _, subset := range subsets {
		if len(subset) == 0 {
			continue
		}

		// simulate the button presses on the buttons in the subset
		lightStatus := make([]int, len(m.desired_light_status))
		for _, buttonIndex := range subset {
			for i := 0; i < len(lightStatus); i++ {
				lightStatus[i] = (lightStatus[i] + m.button_wires[buttonIndex][i]) % 2
			}
		}
		// check if the light status matches the desired status
		if slices.Equal(lightStatus, m.desired_light_status) {
			return len(subset)
		}
	}

	return -1

}

func ButtonJoltageTotal(counters []int) int {
	sum := 0
	for _, c := range counters {
		sum += c
	}
	return sum
}

func isButtonUsable(button []int, targetJoltage []int) bool {
	// if any corresponding value in the button Slice is greater than the targetJoltage, the button is not usable
	for i := 0; i < len(targetJoltage); i++ {
		if button[i] > targetJoltage[i] {
			return false
		}
	}

	return true
}

func isJoltageAchieved(currentJoltage []int) bool {
	for _, val := range currentJoltage {
		if val != 0 {
			return false
		}
	}
	return true
}

// substract button values from currentJoltage, return the difference
func SubstractButtonEffect(currentJoltage []int, button []int) []int {
	result := make([]int, len(currentJoltage))
	for i := 0; i < len(currentJoltage); i++ {
		result[i] = currentJoltage[i] - button[i]
	}
	return result
}

type SolutionCandidate struct {
	buttonPressed         int
	targetJoltageAsString string // comma separated int, as a string
}

// search for minimum buttons to press to achieve the target Joltage.
// using BFS
func SolveForJoltage(buttons [][]int, targetJoltage []int) int {

	solutionQueue := commons.Queue[SolutionCandidate]{}

	// initialize the queue with some solution candidates
	for _, button := range buttons {
		if isButtonUsable(button, targetJoltage) {
			newJoltage := SubstractButtonEffect(targetJoltage, button)
			if isJoltageAchieved(newJoltage) {
				return 1
			}
			solutionCandidate := SolutionCandidate{
				buttonPressed:         1,
				targetJoltageAsString: JoltageToString(newJoltage),
			}
			solutionQueue.Enqueue(solutionCandidate)
		}
	}

	for solutionQueue.Len() > 0 {
		currentCandidate, _ := solutionQueue.Dequeue()
		targetJoltage := ParseJoltage(currentCandidate.targetJoltageAsString)

		for _, button := range buttons {
			if isButtonUsable(button, targetJoltage) {
				newJoltage := SubstractButtonEffect(targetJoltage, button)
				if isJoltageAchieved(newJoltage) {
					return currentCandidate.buttonPressed + 1
				}
				solutionCandidate := SolutionCandidate{
					buttonPressed:         currentCandidate.buttonPressed + 1,
					targetJoltageAsString: JoltageToString(newJoltage),
				}
				solutionQueue.Enqueue(solutionCandidate)
			}
		}
	}

	return -1
}

// find the minimum buttons to press to achieve the jotage
func (m Machine) ConstructJoltage() int {
	return SolveForJoltage(m.button_wires, m.jotage)
}

// parsse light status string into slice of 0 and 1s.
// e.g "[.##.]" to [0,1,1,0]
func ParseLightStatus(input string) []int {
	status := make([]int, 0)
	for _, ch := range input {
		if ch == '#' {
			status = append(status, 1)
		} else if ch == '.' {
			status = append(status, 0)
		}
	}
	return status
}

// parse joltage string into slice of integers.
// e.g "{3,0,4,1}" to [3,0,4,1]
func ParseJoltage(input string) []int {
	joltageList := make([]int, 0)
	content := strings.Trim(input, "{}")

	numbers := strings.Split(content, ",")
	for _, numStr := range numbers {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			panic(err)
		}
		joltageList = append(joltageList, num)
	}
	return joltageList
}

// serialize joltage slice to string
func JoltageToString(joltage []int) string {
	strs := make([]string, len(joltage))
	for i, val := range joltage {
		strs[i] = strconv.Itoa(val)
	}
	return "{" + strings.Join(strs, ",") + "}"
}

func ParseButtonWiring(input string, numOfWires int) []int {
	buttonWires := make([]int, numOfWires)
	content := strings.Trim(input, "()")
	numbers := strings.Split(content, ",")
	for _, numStr := range numbers {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			panic(err)
		}
		buttonWires[num] = 1
	}
	return buttonWires
}

func ParseInput() []Machine {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(input), "\n")
	machines := make([]Machine, len(lines))
	for i, line := range lines {
		machines[i] = Machine{}
		parts := strings.Fields(strings.TrimSpace(line))

		// first part is desired light status
		machines[i].desired_light_status = ParseLightStatus(parts[0])
		// last part is jotage
		machines[i].jotage = ParseJoltage(parts[len(parts)-1])

		// middle parts are button wire schematics
		numLights := len(machines[i].desired_light_status)
		buttonWires := make([][]int, 0)
		for _, part := range parts[1 : len(parts)-1] {
			buttonWires = append(buttonWires, ParseButtonWiring(part, numLights))
		}
		machines[i].button_wires = buttonWires
	}
	return machines
}

func Part1() {
	machines := ParseInput()
	totalPresses := 0
	for i, machine := range machines {
		numPresses := machine.ToggleLights()
		fmt.Printf("Machine %d: Minimum Button Presses = %d\n", i, numPresses)
		totalPresses += numPresses
	}

	fmt.Printf("Total Minimum Button Presses: %d\n", totalPresses)
}

func vectorsToA(vectors [][]int) (A [][]int, dim int, n int) {
	n = len(vectors)
	if n == 0 {
		return nil, 0, 0
	}
	dim = len(vectors[0])

	// Validate
	for i := 0; i < n; i++ {
		if len(vectors[i]) != dim {
			panic("all vectors must have same dimension")
		}
	}

	// A[d][i] = vectors[i][d]
	A = make([][]int, dim)
	for d := 0; d < dim; d++ {
		A[d] = make([]int, n)
		for i := 0; i < n; i++ {
			A[d][i] = vectors[i][d]
		}
	}
	return A, dim, n
}

// solve the linear equation with integer coefficients to achieve the target joltage, using golp
func (machine Machine) SolveJoltageLP() int {

	fmt.Println("Solving Joltage LP for machine:", machine)
	// copy machine.jotage to b
	b := make([]int, len(machine.jotage))
	copy(b, machine.jotage)

	A, dim, n := vectorsToA(machine.button_wires)
	if len(b) != dim {
		log.Fatalf("b dimension %d != vector dimension %d", len(b), dim)
	}
	// Create LP with n variables
	lp := golp.NewLP(0, n)
	lp.SetVerboseLevel(golp.IMPORTANT)

	// Objective: minimize sum(x) (optional but useful)
	obj := make([]float64, n)
	for i := 0; i < n; i++ {
		obj[i] = 1
	}
	lp.SetObjFn(obj)

	// Constraints: for each dimension d, sum_i A[d][i]*x[i] == b[d]
	for d := 0; d < dim; d++ {
		row := make([]float64, n)
		for i := 0; i < n; i++ {
			row[i] = float64(A[d][i])
		}
		lp.AddConstraint(row, golp.EQ, float64(b[d]))
	}

	// set some constraints: x[i] >= 0 and integer
	for i := 0; i < n; i++ {
		lp.SetInt(i, true)
		lp.SetBounds(i, 0, 1000)
	}

	status := lp.Solve()
	if status != golp.OPTIMAL && status != golp.SUBOPTIMAL {
		fmt.Println("No solution. Solve status:", status)
		return -1
	}

	x := lp.Variables()
	fmt.Println("Solution x:", x)

	// calculate total button presses
	totalPresses := 0
	for i := 0; i < n; i++ {
		totalPresses += int(x[i])
	}
	fmt.Println("Total Button Presses:", totalPresses)
	return totalPresses
}

func Part2() {
	machines := ParseInput()

	totalPresses := 0
	for i, machine := range machines {
		numPresses := machine.SolveJoltageLP()
		fmt.Printf("Machine %d: Minimum Button Presses for Joltage = %d\n", i, numPresses)
		totalPresses += numPresses
	}

	fmt.Printf("Total Minimum Button Presses for Joltage: %d\n", totalPresses)
}

func main() {
	fmt.Println("--- Day 10: Factory ---")
	//Part1()
	Part2()
}
