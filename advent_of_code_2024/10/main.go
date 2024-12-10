package main

import (
	"bufio"
	"fmt"
	"os"
)

func getTrailHeads(topographicMap [][]int) [][]int {
	var trailHeads [][]int
	for i, row := range topographicMap {
		for j, cell := range row {
			if cell == 0 {
				trailHeads = append(trailHeads, []int{i, j})
			}
		}
	}
	return trailHeads
}

func getSurroundingPoints(topographicMap [][]int, i, j, val int) [][]int {
	var surroundingCells [][]int
	if i > 0 && topographicMap[i-1][j] == val+1 {
		surroundingCells = append(surroundingCells, []int{i - 1, j})
	}
	if i < len(topographicMap)-1 && topographicMap[i+1][j] == val+1 {
		surroundingCells = append(surroundingCells, []int{i + 1, j})
	}
	if j > 0 && topographicMap[i][j-1] == val+1 {
		surroundingCells = append(surroundingCells, []int{i, j - 1})
	}
	if j < len(topographicMap[0])-1 && topographicMap[i][j+1] == val+1 {
		surroundingCells = append(surroundingCells, []int{i, j + 1})
	}

	return surroundingCells
}

func traceTrails(topographicMap [][]int, startingPoint []int) [][][]int {
	trails := make([][][]int, 0)
	validTrails := make([][][]int, 0)

	// Get surrounding cells
	i, j := startingPoint[0], startingPoint[1]
	surroundingCells := getSurroundingPoints(topographicMap, i, j, topographicMap[i][j])
	if len(surroundingCells) == 0 {
		trail := [][]int{startingPoint}
		trails = append(trails, trail)
	}
	for _, cell := range surroundingCells {
		trail := [][]int{startingPoint}
		branchingTrails := traceTrails(topographicMap, cell)
		for _, branch := range branchingTrails {
			if len(branch) > 0 {
				trails = append(trails, append(trail, branch...))
			}
		}
	}

	// Filter out trails with no valid end
	for _, trail := range trails {
		if len(trail) > 0 && topographicMap[trail[len(trail)-1][0]][trail[len(trail)-1][1]] == 9 {
			validTrails = append(validTrails, trail)
		}
	}
	return validTrails
}

func calculateUniqueTrailsHeadTails(trails [][][]int) int {
	trailsMap := make(map[[2]int]map[[2]int]bool)
	for _, trail := range trails {
		trailHead := [2]int{trail[0][0], trail[0][1]}
		end := [2]int{trail[len(trail)-1][0], trail[len(trail)-1][1]}
		if _, ok := trailsMap[trailHead]; !ok {
			trailsMap[trailHead] = make(map[[2]int]bool)
		}
		trailsMap[trailHead][end] = true
	}

	// Count the number of unique trails
	uniqueTrails := 0
	for _, trailEnds := range trailsMap {
		uniqueTrails += len(trailEnds)
	}
	return uniqueTrails
}

func main() {
	// read data from test_input.txt
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var topographicMap [][]int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var row []int
		for _, char := range line {
			// convert char to int
			row = append(row, int(char-'0'))
		}
		topographicMap = append(topographicMap, row)
	}

	// print topographic map
	for _, row := range topographicMap {
		for _, cell := range row {
			fmt.Print(cell)
		}
		fmt.Println()
	}

	// get trail heads
	trailHeads := getTrailHeads(topographicMap)
	totalTrails := 0
	totalUniqueTrails := 0
	for _, trailHead := range trailHeads {
		trails := traceTrails(topographicMap, trailHead)
		//fmt.Println(trails)
		totalTrails += calculateUniqueTrailsHeadTails(trails)
		totalUniqueTrails += len(trails)
	}
	fmt.Println("Total number of unique trail heads and tails: ", totalTrails)
	fmt.Println("Total number of unique trails: ", totalUniqueTrails)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
