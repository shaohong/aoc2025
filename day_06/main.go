package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ParseInput() (numberRows [][]int, operatorRow []string) {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")

	numberRows = make([][]int, 0)
	operatorRow = make([]string, 0)

	for _, line := range lines {
		tokens := strings.Fields(line)
		_, err := strconv.Atoi(tokens[0])
		if err != nil {
			// this is an operator row
			operatorRow = tokens
		} else {
			// this is a number row
			numberRow := make([]int, 0)
			for _, token := range tokens {
				num, err := strconv.Atoi(token)
				if err != nil {
					panic(err)
				}
				numberRow = append(numberRow, num)
			}
			numberRows = append(numberRows, numberRow)
		}
	}

	// fmt.Println("Number Rows:", numberRows)
	// fmt.Println("Operator Row:", operatorRow)

	return numberRows, operatorRow
}

type Operation struct {
	operator string
	operands []int
}

func (op Operation) Apply() int {
	switch op.operator {
	case "*":
		result := 1
		for _, operand := range op.operands {
			result *= operand
		}
		return result

	case "+":
		result := 0
		for _, operand := range op.operands {
			result += operand
		}
		return result
	}
	return -1
}

func Part1() {
	numberRows, operatorRow := ParseInput()

	operations := make([]Operation, len(operatorRow))

	// loop through columns
	for col := 0; col < len(operatorRow); col++ {
		operations[col] = Operation{
			operator: operatorRow[col],
			operands: make([]int, len(numberRows))}

		for row := 0; row < len(numberRows); row++ {
			operations[col].operands[row] = numberRows[row][col]
		}
	}

	totalSum := 0
	for _, op := range operations {
		// fmt.Println("Operation:", op.operator, "Operands:", op.operands, "Result:", op.Apply())
		totalSum += op.Apply()
	}

	fmt.Println("Total Sum of all operations:", totalSum)
}

func Part2() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")

	// conver to columns of bytes

	// the number of columns is determined by the length of the first line
	numCols := len(lines[0])

	columnBytes := make([][]byte, numCols)
	for i := 0; i < numCols; i++ {
		columnBytes[i] = make([]byte, 0)
		for _, line := range lines {
			if len(line) > 0 {
				columnBytes[i] = append(columnBytes[i], line[i])
			}
		}
	}

	// scan from right to left, parse the operand and operators
	operations := make([]Operation, 0)
	var currentOp Operation = Operation{operands: make([]int, 0)}
	for col := numCols - 1; col >= 0; col-- {
		colStr := string(columnBytes[col])
		colStr = strings.TrimSpace(colStr)
		// fmt.Println("checking colStr", colStr)
		if colStr == "" {
			continue
		}

		// see if the last character is an operator
		if colStr[len(colStr)-1] == '+' || colStr[len(colStr)-1] == '*' {
			// this is the last column of the current operation
			currentOp.operator = colStr[len(colStr)-1:]

			restStr := strings.TrimSpace(colStr[:len(colStr)-1])
			value, _ := strconv.Atoi(restStr)
			currentOp.operands = append(currentOp.operands, value)
			operations = append(operations, currentOp)
			currentOp = Operation{operands: make([]int, 0)}
		} else {
			value, _ := strconv.Atoi(colStr)
			currentOp.operands = append(currentOp.operands, value)
		}
	}

	totalSum := 0
	for _, op := range operations {
		// fmt.Println("Operation:", op.operator, "Operands:", op.operands, "Result:", op.Apply())
		totalSum += op.Apply()
	}

	fmt.Println("Total Sum of all operations:", totalSum)
}

func main() {
	fmt.Println("--- Day 6: Trash Compactor ---")
	//Part1()
	Part2()
}
