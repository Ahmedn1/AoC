package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

func solveEquations(x1, y1, c1, x2, y2, c2, limit int) (int, int, error) {
	if limit == -1 {
		limit = math.MaxInt64
	}
	// Step 1: Calculate the determinants
	deltaX := x1*y2 - x2*y1
	deltaC := c1*y2 - c2*y1

	// Step 2: Check if deltaX divides deltaC perfectly
	if deltaX == 0 || deltaC%deltaX != 0 {
		return 0, 0, fmt.Errorf("no solution found within the limits") // No solution
	}

	// Step 3: Solve for A
	A := deltaC / deltaX
	if A < 0 || A > limit { // Check range for A
		return 0, 0, fmt.Errorf("no solution found within the limits")
	}

	// Step 4: Solve for B using the first equation
	if (c1-x1*A)%y1 != 0 {
		return 0, 0, fmt.Errorf("B is not an integer")
	}
	B := (c1 - x1*A) / y1
	if B < 0 || B > limit {
		return 0, 0, fmt.Errorf("no solution found within the limits")
	}

	return A, B, nil
}

func getTokens(clawMachines []map[string][2]int, prizeOffset, limit int) int {
	tokens := 0
	for _, machine := range clawMachines {
		aX, aY := machine["A"][0], machine["A"][1]
		bX, bY := machine["B"][0], machine["B"][1]
		prizeX, prizeY := machine["Prize"][0]+prizeOffset, machine["Prize"][1]+prizeOffset

		var a, b, cost int
		a, b, err := solveEquations(aX, bX, prizeX, aY, bY, prizeY, limit)
		cost = a*3 + b

		if err == nil {
			tokens += cost
		}
	}
	return tokens
}

func main() {
	file, err := os.Open("./13/input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	coordsRegex := regexp.MustCompile(`[XY][+=]\d+`)

	var clawMachines []map[string][2]int

	scanner := bufio.NewScanner(file)
	var machineLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			machineLines = append(machineLines, line)
		} else {
			clawMachines = append(clawMachines, createMachine(coordsRegex, machineLines))
			machineLines = []string{}
		}
	}
	clawMachines = append(clawMachines, createMachine(coordsRegex, machineLines))

	part1Tokens := getTokens(clawMachines, 0, 100)
	fmt.Println("Part 1: ", part1Tokens)

	part2Tokens := getTokens(clawMachines, 10000000000000, -1)
	fmt.Println("Part 2: ", part2Tokens)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}

func createMachine(coordsRegex *regexp.Regexp, machineLines []string) map[string][2]int {
	buttonAValues := coordsRegex.FindAllString(machineLines[0], -1)
	aX, _ := strconv.Atoi(buttonAValues[0][2:])
	aY, _ := strconv.Atoi(buttonAValues[1][2:])
	buttonBValues := coordsRegex.FindAllString(machineLines[1], -1)
	bX, _ := strconv.Atoi(buttonBValues[0][2:])
	bY, _ := strconv.Atoi(buttonBValues[1][2:])
	prizeValues := coordsRegex.FindAllString(machineLines[2], -1)
	prizeX, _ := strconv.Atoi(prizeValues[0][2:])
	prizeY, _ := strconv.Atoi(prizeValues[1][2:])
	return map[string][2]int{
		"A":     [2]int{aX, aY},
		"B":     [2]int{bX, bY},
		"Prize": [2]int{prizeX, prizeY},
	}
}
