package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func IsRepeatingSequence(productID string) bool {
	repeatingTwice := false
	n := len(productID)
	if n%2 == 0 {
		firstHalf := productID[:n/2]
		secondHalf := productID[n/2:]
		if firstHalf == secondHalf {
			repeatingTwice = true
		}
	}

	return repeatingTwice
}

func IsRepeatingSequenceInteger(productID uint) bool {
	// convert an integer to string and check if it is a repeating sequence
	productIDStr := fmt.Sprintf("%d", productID)
	return IsRepeatingSequence(productIDStr)
}

func HasLeadingZero(productID string) bool {
	return productID[0] == '0'
}

func InvalidProductIDs(lowerBound uint, upperBound uint) []uint {
	invalidProductIDs := make([]uint, 0)
	for i := lowerBound; i <= upperBound; i++ {
		if IsRepeatingSequenceInteger(i) {
			invalidProductIDs = append(invalidProductIDs, i)
		}
	}
	return invalidProductIDs
}

type ProductIDRange struct {
	lowerBound uint
	upperBound uint
}

func ParseInput(input string) []ProductIDRange {
	ranges := make([]ProductIDRange, 0)
	// split the input line by comma,

	input = strings.ReplaceAll(input, "\n", "")

	parts := strings.Split(input, ",")
	for _, part := range parts {
		bounds := strings.Split(part, "-")
		if len(bounds) != 2 {
			continue
		}
		var lower, upper uint
		fmt.Sscanf(bounds[0], "%d", &lower)
		fmt.Sscanf(bounds[1], "%d", &upper)
		ranges = append(ranges, ProductIDRange{lowerBound: lower, upperBound: upper})
	}
	return ranges
}

func Part1() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	productIDRanges := ParseInput(string(data))

	totalSum := 0
	for _, pidRange := range productIDRanges {
		invalidIDs := InvalidProductIDs(pidRange.lowerBound, pidRange.upperBound)
		fmt.Printf("Invalid Product IDs in range %d-%d: %v\n", pidRange.lowerBound, pidRange.upperBound, invalidIDs)
		for _, id := range invalidIDs {
			totalSum += int(id)
		}
	}
	fmt.Printf("Total sum of invalid product IDs: %d\n", totalSum)
}

func isRepeated(s string) bool {
	n := len(s)
	for size := 1; size <= n/2; size++ {
		if n%size != 0 {
			continue
		}
		prefix := s[:size]
		if strings.Repeat(prefix, n/size) == s {
			return true
		}
	}
	return false
}

func Part2() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	productIDRanges := ParseInput(string(data))

	totalSum := 0
	for _, pidRange := range productIDRanges {
		invalidIDs := make([]uint, 0)
		for i := pidRange.lowerBound; i <= pidRange.upperBound; i++ {
			idStr := fmt.Sprintf("%d", i)
			if isRepeated(idStr) {
				invalidIDs = append(invalidIDs, i)
			}
		}
		fmt.Printf("Invalid Product IDs in range %d-%d: %v\n", pidRange.lowerBound, pidRange.upperBound, invalidIDs)
		for _, id := range invalidIDs {
			totalSum += int(id)
		}
	}
	fmt.Printf("Total sum of invalid product IDs: %d\n", totalSum)
}

func main() {
	//Part1()
	Part2()
}
