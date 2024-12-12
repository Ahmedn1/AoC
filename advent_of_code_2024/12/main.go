package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func getSurroundingPlants(garden [][]string, plantLocation [2]int) [][2]int {
	var surroundingCells [][2]int
	i, j := plantLocation[0], plantLocation[1]
	plant := garden[i][j]

	if i > 0 && garden[i-1][j] == plant {
		surroundingCells = append(surroundingCells, [2]int{i - 1, j})
	}
	if i < len(garden)-1 && garden[i+1][j] == plant {
		surroundingCells = append(surroundingCells, [2]int{i + 1, j})
	}
	if j > 0 && garden[i][j-1] == plant {
		surroundingCells = append(surroundingCells, [2]int{i, j - 1})
	}
	if j < len(garden[0])-1 && garden[i][j+1] == plant {
		surroundingCells = append(surroundingCells, [2]int{i, j + 1})
	}

	return surroundingCells
}

func findRegion(garden [][]string, plantLocation [2]int) map[[2]int]bool {
	region := make(map[[2]int]bool)

	region[plantLocation] = true
	lastRegionSize := 0
	for {
		for p, _ := range region {
			surroundingPlants := getSurroundingPlants(garden, p)
			for _, plantLoc := range surroundingPlants {
				if _, ok := region[plantLoc]; !ok {
					region[plantLoc] = true
				}
			}
		}
		if len(region) == lastRegionSize {
			break
		} else {
			lastRegionSize = len(region)
		}
	}

	return region
}

func checkPlantInRegion(regions []map[[2]int]bool, plantLocation [2]int) bool {
	for _, region := range regions {
		if _, ok := region[plantLocation]; ok {
			return true
		}
	}
	return false
}

func addPointToPerimeter(perimeter [][][2]int, i, j int, valueToCompare string, garden [][]string, plantLocation [2]int) [][][2]int {
	if len(perimeter) > 0 {
		condition := false
		perimeterIdx := -1
		if valueToCompare == "row" {
			for idx, p := range perimeter {
				lastPoint := p[len(p)-1]
				if lastPoint[0] == i && lastPoint[1] == j-1 && ((lastPoint[0] == -1 || lastPoint[1] == -1) || garden[plantLocation[0]][plantLocation[1]-1] == garden[plantLocation[0]][plantLocation[1]]) {
					condition = true
					perimeterIdx = idx
					break
				}
			}
		} else {
			for idx, p := range perimeter {
				lastPoint := p[len(p)-1]
				if lastPoint[1] == j && lastPoint[0] == i-1 && ((lastPoint[0] == -1 || lastPoint[1] == -1) || garden[plantLocation[0]-1][plantLocation[1]] == garden[plantLocation[0]][plantLocation[1]]) {
					condition = true
					perimeterIdx = idx
					break
				}
			}
		}
		if condition {
			perimeter[perimeterIdx] = append(perimeter[perimeterIdx], [2]int{i, j})
		} else {
			perimeter = append(perimeter, [][2]int{{i, j}})
		}
	} else {
		perimeter = append(perimeter, [][2]int{{i, j}})
	}
	return perimeter
}

func calcRegionPerimeter(region map[[2]int]bool, garden [][]string) map[string][][][2]int {
	perimeters := make(map[string][][][2]int)
	perimeters["row"] = make([][][2]int, 0)
	perimeters["col"] = make([][][2]int, 0)

	sortedRegion := make([][2]int, 0, len(region))
	for k := range region {
		sortedRegion = append(sortedRegion, k)
	}
	sort.Slice(sortedRegion, func(i, j int) bool {
		if sortedRegion[i][0] == sortedRegion[j][0] {
			return sortedRegion[i][1] < sortedRegion[j][1]
		}
		return sortedRegion[i][0] < sortedRegion[j][0]
	})

	for _, p := range sortedRegion {
		i, j := p[0], p[1]
		if i == 0 {
			perimeters["row"] = addPointToPerimeter(perimeters["row"], i-1, j, "row", garden, p)
		} else if i == len(garden)-1 {
			perimeters["row"] = addPointToPerimeter(perimeters["row"], i, j, "row", garden, p)
		}
		if j == 0 {
			perimeters["col"] = addPointToPerimeter(perimeters["col"], i, j-1, "col", garden, p)
		} else if j == len(garden[0])-1 {
			perimeters["col"] = addPointToPerimeter(perimeters["col"], i, j, "col", garden, p)
		}
		if _, ok := region[[2]int{i - 1, j}]; !ok && i != 0 {
			perimeters["row"] = addPointToPerimeter(perimeters["row"], i-1, j, "row", garden, p)
		}
		if _, ok := region[[2]int{i + 1, j}]; !ok && i != len(garden)-1 {
			perimeters["row"] = addPointToPerimeter(perimeters["row"], i, j, "row", garden, p)
		}
		if _, ok := region[[2]int{i, j - 1}]; !ok && j != 0 {
			perimeters["col"] = addPointToPerimeter(perimeters["col"], i, j-1, "col", garden, p)
		}
		if _, ok := region[[2]int{i, j + 1}]; !ok && j != len(garden[0])-1 {
			perimeters["col"] = addPointToPerimeter(perimeters["col"], i, j, "col", garden, p)
		}
	}
	return perimeters
}

func findRegions(garden [][]string) map[string][]map[[2]int]bool {
	regions := make(map[string][]map[[2]int]bool)

	for i := 0; i < len(garden); i++ {
		for j := 0; j < len(garden[i]); j++ {
			plant := garden[i][j]
			if _, ok := regions[plant]; !ok {
				regions[plant] = []map[[2]int]bool{}
			}
			if !checkPlantInRegion(regions[plant], [2]int{i, j}) {
				plantRegion := findRegion(garden, [2]int{i, j})
				regions[plant] = append(regions[plant], plantRegion)
			}
		}
	}

	return regions
}

func countPerimetersAndSides(allPerimeters map[string][][][2]int) (int, int) {
	perimeterCount := 0
	sideCount := 0
	for _, perimeters := range allPerimeters {
		sideCount += len(perimeters)
		for _, perimeter := range perimeters {
			perimeterCount += len(perimeter)
		}
	}
	return perimeterCount, sideCount
}

func calcRegionsPrice(garden [][]string, regions map[string][]map[[2]int]bool, includePerimeters bool) int {
	price := 0
	for _, plantRegions := range regions {
		for _, region := range plantRegions {
			regionPerimeter := calcRegionPerimeter(region, garden)
			perimeterCount, sideCount := countPerimetersAndSides(regionPerimeter)
			regionArea := len(region)
			if includePerimeters {
				price += regionArea * perimeterCount
			} else {
				price += regionArea * sideCount
			}
		}
	}
	return price
}

func main() {
	file, err := os.Open("./12/input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var garden [][]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		garden = append(garden, strings.Split(line, ""))
	}

	regions := findRegions(garden)
	//fmt.Println(regions)
	price := calcRegionsPrice(garden, regions, true)
	fmt.Println(price)
	price = calcRegionsPrice(garden, regions, false)
	fmt.Println(price)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
