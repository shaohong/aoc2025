// https://adventofcode.com/2025/day/3

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// generate the largest two-digit number by picking digits from the input slice
// the ten's digt must come before the one's digit in the input slice
func LargestTwoDigitNumber(numberStr string) int {
	numbers := make([]int, 0)
	for _, ch := range numberStr {
		numbers = append(numbers, int(ch-'0'))
	}

	if len(numbers) < 2 {
		return -1 // or some error value indicating not enough digits
	}

	// find the largest digits from numbers[:len(numbers)-2]
	maxTens := -1
	maxTensIndex := -1
	for i := 0; i <= len(numbers)-2; i++ {
		if numbers[i] > maxTens {
			maxTens = numbers[i]
			maxTensIndex = i
		}
	}

	// find the largest digit from numbers[maxTensIndex+1:]
	maxOnes := -1
	for i := maxTensIndex + 1; i < len(numbers); i++ {
		if numbers[i] > maxOnes {
			maxOnes = numbers[i]
		}
	}

	return maxTens*10 + maxOnes

}

func LargestNDigitNumber(numberStr string, ndigits int) int {
	numbers := make([]int, 0)
	for _, ch := range numberStr {
		numbers = append(numbers, int(ch-'0'))
	}

	if len(numbers) < ndigits {
		panic(fmt.Sprintf("Input string %s has less than %d digits", numberStr, ndigits))
	}

	maxDigits := make([]int, ndigits)

	// find the largest digit and its index from numbers[:len(numbers)-n+1]
	// after that, find the next largest from the remaining slice, and so on
	startIndex := 0
	for i := 1; i <= ndigits; i++ {
		maxDigit := -1
		maxIndex := -1
		for j := startIndex; j < len(numbers)-ndigits+i; j++ {
			if numbers[j] > maxDigit {
				maxDigit = numbers[j]
				maxIndex = j
			}
		}
		maxDigits[i-1] = maxDigit
		startIndex = maxIndex + 1
	}

	// construct the largest n-digit number
	largestNumber := 0
	for i := 0; i < ndigits; i++ {
		largestNumber = largestNumber*10 + maxDigits[i]
	}

	return largestNumber
}

func ParseInput(r io.Reader) []string {
	data, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")
	return lines
}

func Part1() {

	inputs := ParseInput(os.Stdin)

	sumJoltages := 0
	for _, line := range inputs {
		joltage := LargestTwoDigitNumber(line)
		sumJoltages += joltage
	}

	fmt.Println("Sum of largest two-digit joltages:", sumJoltages)
}

func Part2() {
	inputs := ParseInput(os.Stdin)

	sumJoltages := 0
	const ndigits = 12
	for _, line := range inputs {

		joltage := LargestNDigitNumber(line, ndigits)
		sumJoltages += joltage
	}

	fmt.Println("Sum of largest four-digit joltages:", sumJoltages)
}

func main() {
	// Part1()
	Part2()
}
