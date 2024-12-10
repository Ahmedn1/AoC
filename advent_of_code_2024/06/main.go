package main

import (
	"bufio"
	"fmt"
	"github.com/sbwhitecap/tqdm"
	. "github.com/sbwhitecap/tqdm/iterators"
	"os"
)

func getGuardIndex(guardMap [][]string) []int {
	var guardIndex []int
	for i, row := range guardMap {
		for j, cell := range row {
			if cell == "^" {
				guardIndex = append(guardIndex, i, j)
				return guardIndex
			}
		}
	}
	return guardIndex
}

func buildMapVisitedLocations(visitedLocations [][3]int) map[[3]int]bool {
	uniqueVisitedLocations := make(map[[3]int]bool)
	for _, location := range visitedLocations {
		uniqueVisitedLocations[[3]int{location[0], location[1], location[2]}] = true
	}
	return uniqueVisitedLocations
}

func getGuardSteps(guardMap [][]string) [][3]int {
	// Create a visitedLocations slice to store visited indices
	// visitedLocations is a slice of slices, where each slice contains 3 integers
	// The first integer is the row index, the second integer is the column index, and the third integer represents the direction
	// ^ = 0, > = 1, v = 2, < = 3
	visitedLocations := make([][3]int, 0, len(guardMap))

	guardIndex := getGuardIndex(guardMap)
	noIterationsVisitedLocationsNotUpdated := 0
	guardDirection := guardMap[guardIndex[0]][guardIndex[1]]
	visitedLocations = append(visitedLocations, [3]int{guardIndex[0], guardIndex[1], 0})
	lastVisitedMapLocationsSize := len(buildMapVisitedLocations(visitedLocations))
	for {
		if noIterationsVisitedLocationsNotUpdated > 3 {
			return make([][3]int, 0)
		}

		if guardDirection == "^" {
			if guardIndex[0] == 0 {
				break
			} else if guardMap[guardIndex[0]-1][guardIndex[1]] == "#" {
				guardDirection = ">"
				//visitedLocations = append(visitedLocations, [3]int{guardIndex[0], guardIndex[1], 1})
			} else {
				guardIndex[0]--
				visitedLocations = append(visitedLocations, [3]int{guardIndex[0], guardIndex[1], 0})
			}
		} else if guardDirection == ">" {
			if guardIndex[1] == len(guardMap[0])-1 {
				break
			} else if guardMap[guardIndex[0]][guardIndex[1]+1] == "#" {
				guardDirection = "v"
				//visitedLocations = append(visitedLocations, [3]int{guardIndex[0], guardIndex[1], 2})
			} else {
				guardIndex[1]++
				visitedLocations = append(visitedLocations, [3]int{guardIndex[0], guardIndex[1], 1})
			}
		} else if guardDirection == "v" {
			if guardIndex[0] == len(guardMap)-1 {
				break
			} else if guardMap[guardIndex[0]+1][guardIndex[1]] == "#" {
				guardDirection = "<"
				//visitedLocations = append(visitedLocations, [3]int{guardIndex[0], guardIndex[1], 3})
			} else {
				guardIndex[0]++
				visitedLocations = append(visitedLocations, [3]int{guardIndex[0], guardIndex[1], 2})
			}
		} else if guardDirection == "<" {
			if guardIndex[1] == 0 {
				break
			} else if guardMap[guardIndex[0]][guardIndex[1]-1] == "#" {
				guardDirection = "^"
				//visitedLocations = append(visitedLocations, [3]int{guardIndex[0], guardIndex[1], 0})
			} else {
				guardIndex[1]--
				visitedLocations = append(visitedLocations, [3]int{guardIndex[0], guardIndex[1], 3})
			}
		}

		newVisitedMapLocationsSize := len(buildMapVisitedLocations(visitedLocations))
		if newVisitedMapLocationsSize == lastVisitedMapLocationsSize {
			noIterationsVisitedLocationsNotUpdated++
		} else {
			lastVisitedMapLocationsSize = newVisitedMapLocationsSize
			noIterationsVisitedLocationsNotUpdated = 0
		}
		lastVisitedMapLocationsSize = newVisitedMapLocationsSize
	}
	return visitedLocations
}

func searchForValidPath(uniqueVisitedLocationsWithDirection map[[2]int][]int, location [3]int, guardMap [][]string) bool {
	var searchDirection int
	if location[2] < 3 {
		searchDirection = location[2] + 1
	} else {
		searchDirection = 0
	}
	startRow := location[0]
	startColumn := location[1]

	if searchDirection == 0 {
		for i := startRow - 1; i >= 0; i-- {
			if guardMap[i][startColumn] == "#" {
				break
			}
			if value, exists := uniqueVisitedLocationsWithDirection[[2]int{i, startColumn}]; exists {
				for _, v := range value {
					if v == searchDirection {
						return true
					}
				}
			}
		}
	} else if searchDirection == 1 {
		for i := startColumn + 1; i < len(guardMap[0]); i++ {
			if guardMap[startRow][i] == "#" {
				break
			}
			if value, exists := uniqueVisitedLocationsWithDirection[[2]int{startRow, i}]; exists {
				for _, v := range value {
					if v == searchDirection {
						return true
					}
				}
			}
		}
	} else if searchDirection == 2 {
		for i := startRow + 1; i < len(guardMap); i++ {
			if guardMap[i][startColumn] == "#" {
				break
			}
			if value, exists := uniqueVisitedLocationsWithDirection[[2]int{i, startColumn}]; exists {
				for _, v := range value {
					if v == searchDirection {
						return true
					}
				}
			}
		}
	} else if searchDirection == 3 {
		for i := startColumn - 1; i >= 0; i-- {
			if guardMap[startRow][i] == "#" {
				break
			}
			if value, exists := uniqueVisitedLocationsWithDirection[[2]int{startRow, i}]; exists {
				for _, v := range value {
					if v == searchDirection {
						return true
					}
				}
			}
		}
	}

	return false
}

func getNewObstaclesCount(visitedLocations [][3]int, guardMap [][]string) int {
	newObstaclesCount := 0
	uniqueVisitedLocationsWithDirection := make(map[[2]int][]int)

	// We can't create a loop before 3 turns
	turnCounter := 0
	for i, location := range visitedLocations {
		currentDirection := location[2]
		uniqueVisitedLocationsWithDirection[[2]int{location[0], location[1]}] = append(uniqueVisitedLocationsWithDirection[[2]int{location[0], location[1]}], currentDirection)
		if i > 0 && currentDirection != visitedLocations[i-1][2] {
			turnCounter++
			//continue
		}
		if i < len(visitedLocations)-1 && currentDirection != visitedLocations[i+1][2] {
			// Ignore this as there is already an obstacle
			continue
		}
		if turnCounter < 3 {
			continue
		} else {
			if searchForValidPath(uniqueVisitedLocationsWithDirection, location, guardMap) {
				newObstaclesCount++
			}
		}
	}
	return newObstaclesCount
}

func part2BruteForce(guardMap [][]string) int {
	newObstacleCount := 0
	for i := 0; i < len(guardMap); i++ {
		tqdm.With(Interval(0, len(guardMap[0])), fmt.Sprintf("Processing Guard Map Row (%d)", i), func(v2 interface{}) (brk2 bool) {
			j := v2.(int)
			cell := guardMap[i][j]
			if cell == "." {
				guardMap[i][j] = "#"
				visitedLocations := getGuardSteps(guardMap)
				if len(visitedLocations) == 0 {
					// There is a cycle in the guard's steps
					newObstacleCount++
				}
				guardMap[i][j] = "."
			}
			return
		})
	}
	return newObstacleCount
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
	var guardMap [][]string

	// Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Read each character in the line and add it as a list of characters to guardMap
		var line []string
		for _, char := range scanner.Text() {
			line = append(line, string(char))
		}
		guardMap = append(guardMap, line)
	}

	visitedLocations := getGuardSteps(guardMap)
	// Count unique values in visitedLocations
	uniqueLocations := make(map[[2]int]bool)
	for _, location := range visitedLocations {
		uniqueLocations[[2]int{location[0], location[1]}] = true
	}
	fmt.Println("Guard Steps:", len(uniqueLocations))

	fmt.Println("New Obstacles Count:", part2BruteForce(guardMap))

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
