package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func AbsInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func getTotalDistance(locationIDs_1, locationIDs_2 []int) int {
	totalDistance := 0
	sort.Ints(locationIDs_1)
	sort.Ints(locationIDs_2)

	for i := 0; i < len(locationIDs_1); i++ {
		totalDistance += AbsInt(locationIDs_1[i] - locationIDs_2[i])
	}

	return totalDistance
}

func countElement(slice []int, element int) int {
	count := 0
	for _, value := range slice {
		if value == element {
			count++
		}
	}
	return count
}

func calculateSimilarityScore(locationIDs_1, locationIDs_2 []int) int {
	similarity := 0
	for i := 0; i < len(locationIDs_1); i++ {
		element := locationIDs_1[i]
		elementCount := countElement(locationIDs_2, element)
		similarity += element * elementCount
	}
	return similarity
}

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	inputFilePath11 := "./1_1_input.txt"

	// Open the file
	file, err := os.Open(inputFilePath11)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Initialize slices to hold integers
	var list1, list2 []int

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Get the current line
		line := scanner.Text()

		// Split the line by tab
		parts := strings.Split(line, "   ")
		if len(parts) != 2 {
			fmt.Println("Invalid line format:", line)
			continue
		}

		// Convert strings to integers
		num1, err1 := strconv.Atoi(parts[0])
		num2, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			fmt.Println("Error converting to integer:", line)
			continue
		}

		// Append to respective lists
		list1 = append(list1, num1)
		list2 = append(list2, num2)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Calculate the total distance
	totalDistance := getTotalDistance(list1, list2)
	fmt.Println("Total distance:", totalDistance)

	// Calculate the similarity score
	similarityScore := calculateSimilarityScore(list1, list2)
	fmt.Println("Similarity score:", similarityScore)
}
