package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func reverseString(s string) string {
	// Convert the string to a slice of runes
	runes := []rune(s)

	// Reverse the slice of runes
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Convert the runes back to a string
	return string(runes)
}

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func getHorizontalLines(lines []string) []string {
	// Returns a list of strings that are constructed by concatenating the characters in each column of each line
	var horizontalLines []string
	for i := 0; i < len(lines[0]); i++ {
		var line string
		for j := 0; j < len(lines); j++ {
			line += string(lines[j][i])
		}
		horizontalLines = append(horizontalLines, line)
	}
	return horizontalLines
}

func getDiagonalLines(lines []string) []string {
	// Returns a list of strings that are constructed by concatenating the characters in each diagonal of the grid
	var diagonalLines []string
	for i := 0; i < len(lines); i++ {
		var line string
		for j := 0; j < len(lines); j++ {
			if i+j < len(lines) {
				line += string(lines[j][i+j])
			}
		}
		diagonalLines = append(diagonalLines, line)
	}

	// Get the remaining diagonals
	for i := 1; i < len(lines); i++ {
		var line string
		for j := 0; j < len(lines); j++ {
			if i+j < len(lines) {
				line += string(lines[i+j][j])
			}
		}
		diagonalLines = append(diagonalLines, line)
	}
	return diagonalLines
}

func getRightToLeftDiagonalLines(lines []string) []string {
	var diagonalLines []string
	// Get right-to-left diagonals
	for i := len(lines) - 1; i >= 0; i-- {
		var line string
		for j := 0; j < len(lines); j++ {
			if i-j >= 0 {
				line += string(lines[j][i-j])
			}
		}
		diagonalLines = append(diagonalLines, line)
	}

	for i := 1; i < len(lines); i++ {
		var line string
		for j := 0; j < len(lines); j++ {
			if i+j < len(lines) {
				line += string(lines[i+j][len(lines)-j-1])
			}
		}
		diagonalLines = append(diagonalLines, line)
	}

	return diagonalLines
}

func countOccurrences(lines []string, word string) int {
	// Returns the number of times the word appears in the lines
	count := 0
	xmasRegex := regexp.MustCompile(word)
	for _, line := range lines {
		count += len(xmasRegex.FindAllString(line, -1))
	}
	return count
}

func getXMASCount(lines []string) int {
	// Returns the number of times 2 word "MAS" appear in the lines with the shape of an X
	count := 0
	for i := 1; i < len(lines)-1; i++ {
		for j := 1; j < len(lines[i])-1; j++ {
			if lines[i][j] == 'A' {
				if ((lines[i-1][j-1] == 'M' && lines[i+1][j+1] == 'S') || (lines[i-1][j-1] == 'S' && lines[i+1][j+1] == 'M')) && ((lines[i-1][j+1] == 'M' && lines[i+1][j-1] == 'S') || (lines[i-1][j+1] == 'S' && lines[i+1][j-1] == 'M')) {
					count++
				}
			}
		}
	}
	return count
}

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	inputFilePath := "./input.txt"

	// Open the file
	file, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		// Get the current line
		lines = append(lines, scanner.Text())
	}

	horizontalLines := getHorizontalLines(lines)
	diagonalLines := getDiagonalLines(lines)
	rightToLeftDiagonalLines := getRightToLeftDiagonalLines(lines)

	var originalLines []string
	originalLines = append(originalLines, lines...)

	lines = append(lines, horizontalLines...)
	lines = append(lines, diagonalLines...)
	lines = append(lines, rightToLeftDiagonalLines...)

	reversedLines := make([]string, len(lines))
	for i, line := range lines {
		reversedLines[i] = reverseString(line)
	}
	lines = append(lines, reversedLines...)

	// Count the number of times the word "XMAS" appears in the lines
	word := "XMAS"
	count := countOccurrences(lines, word)
	fmt.Printf("The word %s appears %d times in the lines\n", word, count)

	// Count the number of times the word "MAS" appears in the lines with the shape of an X
	count = getXMASCount(originalLines)
	fmt.Printf("The word MAS appears %d times in the lines with the shape of an X\n", count)

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
