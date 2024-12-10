package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func isReportSafe(levels []int) bool {
	// Check if the report is safe
	var direction string
	if len(levels) > 1 && levels[0] < levels[1] {
		direction = "increase"
	} else {
		direction = "decrease"
	}

	for i := 0; i < len(levels)-1; i++ {
		if !(AbsInt(levels[i]-levels[i+1]) >= 1 && AbsInt(levels[i]-levels[i+1]) <= 3) {
			return false
		}
		if i > 0 {
			if levels[i] > levels[i+1] && direction == "increase" {
				return false
			} else if levels[i] < levels[i+1] && direction == "decrease" {
				return false
			}
		}
	}

	return true
}

func isReportSafeWithDampener(levels []int) bool {
	// Check if the report is safe with the Problem Dampener

	if isReportSafe(levels) {
		return true
	} else {
		for i := 0; i < len(levels); i++ {
			dampenedSlice := slices.Concat(levels[:i], levels[i+1:])
			if isReportSafe(dampenedSlice) {
				return true
			}
		}
	}

	return false
}

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	inputFilePath := "input.txt"

	// Open the file
	file, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	safeReportsCount := 0
	safeReportsCountWithDampener := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Get the current line
		report := scanner.Text()
		parts := strings.Split(report, " ")
		levels := func(strings []string) []int {
			ints := make([]int, len(strings))
			for i, s := range strings {
				ints[i], _ = strconv.Atoi(s)
			}
			return ints
		}(parts)

		if isReportSafe(levels) {
			safeReportsCount++
		}
		if isReportSafeWithDampener(levels) {
			safeReportsCountWithDampener++
		}
	}

	// Print the number of safe reports
	fmt.Println("Number of safe reports:", safeReportsCount)
	fmt.Println("Number of safe reports with dampener:", safeReportsCountWithDampener)
}
