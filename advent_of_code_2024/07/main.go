package main

import (
	"bufio"
	"fmt"
	"github.com/mowshon/iterium"
	"os"
	"strconv"
	"strings"
)

func parseEquation(line string) []int {
	splitLine := strings.Split(line, ":")
	testValue, err := strconv.Atoi(splitLine[0])
	if err != nil {
		return nil
	}
	// Split the string list of numbers into a cast slice of integers
	operands := make([]int, 0)
	stringNumbers := strings.Split(splitLine[1], " ")
	for _, stringNumber := range stringNumbers {
		number, err := strconv.Atoi(stringNumber)
		if err == nil {
			operands = append(operands, number)
		}
	}
	equation := []int{testValue}
	equation = append(equation, operands...)
	return equation
}

func checkValidEquation(testValue int, operands []int, operatorsCombinations [][]string) bool {
	// Check if the equation is valid

	for _, operators := range operatorsCombinations {
		// Check if the equation is valid
		result := 0
		for i, operand := range operands {
			if i == 0 {
				result = operand
			} else {
				if operators[i-1] == "+" {
					result += operand
				} else if operators[i-1] == "*" {
					result *= operand
				}
			}
		}
		if result == testValue {
			return true
		}
	}

	return false
}

func checkEquationWithConcatenation(testValue int, operands []int, operatorsCombinations [][]string) bool {
	// The concatenation operator (||) combines the digits from its left and right inputs into a single number.
	// For example, 12 || 345 would become 12345. All operators are still evaluated left-to-right.

	// Check if the equation is valid
	for _, operators := range operatorsCombinations {
		// Check if the equation is valid
		result := 0
		for i, operand := range operands {
			if i == 0 {
				result = operand
			} else {
				if operators[i-1] == "+" {
					result += operand
				} else if operators[i-1] == "*" {
					result *= operand
				} else if operators[i-1] == "||" {
					// concatenate result and operand as strings then cast them back to int
					result, _ = strconv.Atoi(strconv.Itoa(result) + strconv.Itoa(operand))
				}
			}
		}
		if result == testValue {
			return true
		}
	}
	return false
}

func getEquationsTotal(equations map[int][]int, part1Only bool) int {
	total := 0
	for testValue, operands := range equations {
		combs := iterium.Product([]string{"+", "*"}, len(operands)-1)
		combinations, _ := combs.Slice()
		if checkValidEquation(testValue, operands, combinations) {
			total += testValue
		} else if !part1Only {
			combs = iterium.Product([]string{"+", "*", "||"}, len(operands)-1)
			combinations, _ = combs.Slice()
			if checkEquationWithConcatenation(testValue, operands, combinations) {
				total += testValue
			}
		}
	}
	return total
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
	equations := make(map[int][]int)

	// Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		equation := parseEquation(scanner.Text())
		if equation != nil {
			equations[equation[0]] = equation[1:]
		}
	}

	fmt.Println("Part 1:", getEquationsTotal(equations, true))
	fmt.Println("Part 2:", getEquationsTotal(equations, false))

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
