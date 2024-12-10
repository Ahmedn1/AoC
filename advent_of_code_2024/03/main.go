package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func executeInstructions(instructions string) int {
	// Initialize the accumulator and the instruction pointer
	result := 0
	pattern := `mul\(\d{1,3},\d{1,3}\)`

	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(instructions, -1)
	for _, match := range matches {
		sepIndex := strings.Index(match, ",")
		n1, _ := strconv.Atoi(match[4:sepIndex])
		n2, _ := strconv.Atoi(match[sepIndex+1 : len(match)-1])
		result += n1 * n2
	}
	return result
}

func executeConditionalInstructions(instructions string) int {
	// Initialize the accumulator and the instruction pointer
	result := 0
	pattern := `mul\(\d{1,3},\d{1,3}\)`
	doPattern := `do\(\)`
	dontPattern := `don't\(\)`

	reMul := regexp.MustCompile(pattern)
	reDo := regexp.MustCompile(doPattern)
	reDont := regexp.MustCompile(dontPattern)

	doMatches := reDo.FindAllStringIndex(instructions, -1)
	dontMatches := reDont.FindAllStringIndex(instructions, -1)

	var conditions []struct {
		start  int
		end    int
		source int // 1 for doMatches, 0 for dontMatches
	}
	for _, match := range doMatches {
		conditions = append(conditions, struct {
			start  int
			end    int
			source int
		}{start: match[0], end: match[1], source: 1})
	}

	for _, match := range dontMatches {
		conditions = append(conditions, struct {
			start  int
			end    int
			source int
		}{start: match[0], end: match[1], source: 0})
	}

	sort.Slice(conditions, func(i, j int) bool {
		return conditions[i].start < conditions[j].start
	})

	// Find MUL instructions
	var matches []string
	if conditions[0].start != 0 {
		matches = append(matches, reMul.FindAllString(instructions[:conditions[0].start], -1)...)
	}
	for i := 0; i < len(conditions); i++ {
		condition := conditions[i]
		if condition.source == 1 {
			if i < len(conditions)-1 {
				matches = append(matches, reMul.FindAllString(instructions[condition.end:conditions[i+1].start], -1)...)
			} else {
				matches = append(matches, reMul.FindAllString(instructions[condition.start:], -1)...)
			}
		}
	}

	for _, match := range matches {
		sepIndex := strings.Index(match, ",")
		n1, _ := strconv.Atoi(match[4:sepIndex])
		n2, _ := strconv.Atoi(match[sepIndex+1 : len(match)-1])
		result += n1 * n2
	}
	return result
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
	instructions, err := os.ReadFile(inputFilePath)
	// Execute the instructions
	//instructions := "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))"
	result := executeInstructions(string(instructions))
	conditionalResult := executeConditionalInstructions(string(instructions))
	fmt.Println("Result:", result)
	fmt.Println("Conditional Result:", conditionalResult)
}
