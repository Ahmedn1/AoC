package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseIntList(s string) []int {
	parts := strings.Split(s, ",")
	numbers := make([]int, len(parts))
	for i, p := range parts {
		num, err := strconv.Atoi(p)
		if err != nil {
			fmt.Println("Error converting to int:", err)
			continue
		}
		numbers[i] = num
	}
	return numbers
}

func buildRules(ruleList []string) map[int][]int {
	rules := make(map[int][]int)
	for _, rule := range ruleList {
		parts := strings.Split(rule, "|")
		firstPage, _ := strconv.Atoi(parts[0])
		secondPage, _ := strconv.Atoi(parts[1])
		rules[firstPage] = append(rules[firstPage], secondPage)
	}
	return rules
}

func checkValidUpdate(update []int, ruleMap map[int][]int) bool {
	for i := 0; i < len(update); i++ {
		currentPage := update[i]
		for j := i + 1; j < len(update); j++ {
			testPage := update[j]
			_, exists := ruleMap[testPage]
			if exists {
				for _, nextPage := range ruleMap[testPage] {
					if nextPage == currentPage {
						return false
					}
				}
			}
		}
	}
	return true
}

func getMiddleNumber(numbers []int) int {
	if len(numbers)%2 == 0 {
		return numbers[len(numbers)/2]
	}
	return numbers[len(numbers)/2]
}

func getPart1Result(rules []string, updates [][]int) int {
	ruleMap := buildRules(rules)
	totalMiddleNumbers := 0
	for _, update := range updates {
		if checkValidUpdate(update, ruleMap) {
			totalMiddleNumbers += getMiddleNumber(update)
		}
	}
	return totalMiddleNumbers
}

func fixUpdate(update []int, ruleMap map[int][]int) []int {
	for i := 0; i < len(update); i++ {
		swapIndex := i
		currentPage := update[i]
		for j := i + 1; j < len(update); j++ {
			testPage := update[j]
			_, exists := ruleMap[testPage]
			if exists {
				for _, nextPage := range ruleMap[testPage] {
					if nextPage == currentPage {
						update[swapIndex] = testPage
						swapIndex = j
					}
				}
				update[swapIndex] = currentPage
			}
		}
	}
	return update
}

func getPart2Result(rules []string, updates [][]int) int {
	ruleMap := buildRules(rules)
	totalMiddleNumbers := 0
	for _, update := range updates {
		if !checkValidUpdate(update, ruleMap) {
			fixedUpdate := fixUpdate(update, ruleMap)
			for !checkValidUpdate(fixedUpdate, ruleMap) {
				fixedUpdate = fixUpdate(fixedUpdate, ruleMap)
			}
			totalMiddleNumbers += getMiddleNumber(fixedUpdate)
		}
	}
	return totalMiddleNumbers
}

func main() {
	// read data from test_input.txt
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Initialize variables
	var rules []string
	var updates [][]int
	isSecondPart := false

	// Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for empty line to switch between parts
		if line == "" {
			isSecondPart = true
			continue
		}

		if isSecondPart {
			// Parse line into a slice of integers
			numbers := parseIntList(line)
			updates = append(updates, numbers)
		} else {
			// Store line as is
			rules = append(rules, line)
		}
	}

	// Get result
	result := getPart1Result(rules, updates)
	fmt.Println("Result:", result)

	// Get result for part 2
	result2 := getPart2Result(rules, updates)
	fmt.Println("Result for part 2:", result2)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
