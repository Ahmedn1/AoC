package main

import (
	"bufio"
	"fmt"
	"os"
)

func createFrequenciesMap(cityMap [][]string) map[rune][][2]int {
	frequencies := make(map[rune][][2]int)
	for i, row := range cityMap {
		for j, char := range row {
			if char != "." {
				runeChar := []rune(char)[0]
				if _, ok := frequencies[runeChar]; !ok {
					frequencies[runeChar] = [][2]int{{i, j}}
				} else {
					frequencies[runeChar] = append(frequencies[runeChar], [2]int{i, j})
				}
			}
		}
	}
	return frequencies
}

func getAntiNodes(frequencies map[rune][][2]int, cityMap [][]string) map[[2]int]bool {
	antiNodes := make(map[[2]int]bool)
	for _, freq := range frequencies {
		for i := 0; i < len(freq); i++ {
			for j := i + 1; j < len(freq); j++ {
				yDiff := freq[j][0] - freq[i][0]
				xDiff := freq[i][1] - freq[j][1]

				antinodes := make([][2]int, 2)
				antinodes[0] = [2]int{freq[i][0] - yDiff, freq[i][1] + xDiff}
				antinodes[1] = [2]int{freq[j][0] + yDiff, freq[j][1] - xDiff}

				for _, antinode := range antinodes {
					if antinode[0] >= 0 && antinode[0] < len(cityMap) && antinode[1] >= 0 && antinode[1] < len(cityMap[0]) {
						antiNodes[antinode] = true
					}
				}
			}
		}
	}
	return antiNodes
}

func getPart2AntiNodes(frequencies map[rune][][2]int, cityMap [][]string) map[[2]int]bool {
	antiNodes := make(map[[2]int]bool)
	for _, freq := range frequencies {
		for i := 0; i < len(freq); i++ {
			for j := i + 1; j < len(freq); j++ {
				yDiff := freq[j][0] - freq[i][0]
				xDiff := freq[i][1] - freq[j][1]

				//Add the antenna locations first as antinodes
				antiNodes[freq[i]] = true
				antiNodes[freq[j]] = true

				iteratorXDiff := xDiff
				iteratorYDiff := yDiff
				// Get antinodes before first antenna
				for {
					antinode := [2]int{freq[i][0] - iteratorYDiff, freq[i][1] + iteratorXDiff}
					if antinode[0] >= 0 && antinode[0] < len(cityMap) && antinode[1] >= 0 && antinode[1] < len(cityMap[0]) {
						antiNodes[antinode] = true
					} else {
						break
					}
					iteratorXDiff += xDiff
					iteratorYDiff += yDiff
				}

				// Get antinodes after second antenna
				iteratorXDiff = xDiff
				iteratorYDiff = yDiff
				for {
					antinode := [2]int{freq[j][0] + iteratorYDiff, freq[j][1] - iteratorXDiff}
					if antinode[0] >= 0 && antinode[0] < len(cityMap) && antinode[1] >= 0 && antinode[1] < len(cityMap[0]) {
						antiNodes[antinode] = true
					} else {
						break
					}
					iteratorXDiff += xDiff
					iteratorYDiff += yDiff
				}
			}
		}
	}
	return antiNodes
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
	var cityMap [][]string

	// Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var line []string
		for _, char := range scanner.Text() {
			line = append(line, string(char))
		}
		cityMap = append(cityMap, line)
	}

	frequencies := createFrequenciesMap(cityMap)
	antiNodes := getAntiNodes(frequencies, cityMap)
	fmt.Println("Anti-nodes:", len(antiNodes))

	antiNodesPart2 := getPart2AntiNodes(frequencies, cityMap)
	fmt.Println("Anti-nodes Part 2:", len(antiNodesPart2))

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
