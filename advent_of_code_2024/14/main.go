package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Robot struct {
	x  int
	y  int
	vX int
	vY int
}

func moveRobot(robot *Robot, time int, spaceX int, spaceY int) {
	dX := (robot.x + robot.vX*time) % spaceX
	if dX < 0 {
		dX += spaceX
	}
	dY := (robot.y + robot.vY*time) % spaceY
	if dY < 0 {
		dY += spaceY
	}
	robot.x = dX
	robot.y = dY
}

func moveRobots(robots []Robot, time int, spaceX int, spaceY int) {
	for i := range robots {
		moveRobot(&robots[i], time, spaceX, spaceY)
	}
}

func countRobotsinQuadrants(robots []Robot, spaceX int, spaceY int) [4]int {
	quadrants := [4]int{0, 0, 0, 0}
	middleX, middleY := spaceX/2, spaceY/2
	for _, robot := range robots {
		if robot.x < middleX && robot.y < middleY {
			quadrants[0]++
		} else if robot.x > middleX && robot.y < middleY {
			quadrants[1]++
		} else if robot.x < middleX && robot.y > middleY {
			quadrants[2]++
		} else if robot.x > middleX && robot.y > middleY {
			quadrants[3]++
		}
	}
	return quadrants
}

func visualizeRobots(robots []Robot, width, height int) {
	// Create an empty grid initialized with '.'
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}

	// Place '*' at the points
	for _, robot := range robots {
		x, y := robot.x, robot.y
		if x >= 0 && x < width && y >= 0 && y < height { // Check bounds
			grid[y][x] = '*'
		}
	}

	// Move cursor to the top-left corner
	fmt.Print("\033[H\033[2J") // Clear the screen
	fmt.Print("\033[H")        // Move cursor to top-left

	// Print the grid
	for _, row := range grid {
		for _, cell := range row {
			fmt.Print(string(cell))
		}
		fmt.Println()
	}
}

func largeColExists(robots []Robot, spaceX, threshold int) bool {
	robotsXs := make(map[int]int)
	for _, robot := range robots {
		robotsXs[robot.x]++
	}

	for _, count := range robotsXs {
		if count >= threshold {
			return true
		}
	}
	return false
}

func largeRowExists(robots []Robot, spaceY, threshold int) bool {
	robotsYs := make(map[int]int)
	for _, robot := range robots {
		robotsYs[robot.y]++
	}

	for _, count := range robotsYs {
		if count >= threshold {
			return true
		}
	}
	return false
}

func isXmasTree(robots []Robot, spaceX, spaceY int) bool {
	largeRow := largeRowExists(robots, spaceY, 30)
	largeCol := largeColExists(robots, spaceX, 30)
	return largeRow && largeCol
}

func findXmasTree(robots []Robot, spaceX int, spaceY int) int {
	for time := 0; time < 100000; time++ {
		moveRobots(robots, 1, spaceX, spaceY)
		if isXmasTree(robots, spaceX, spaceY) {
			return time + 1
		}
	}
	return -1
}

func main() {
	file, err := os.Open("./14/input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	//const spaceX, spaceY = 11, 7
	const spaceX, spaceY = 101, 103
	robotRegex := regexp.MustCompile(`-?\d+`)

	robots := make([]Robot, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		extracted := robotRegex.FindAllString(line, -1)
		x, _ := strconv.Atoi(extracted[0])
		y, _ := strconv.Atoi(extracted[1])
		vX, _ := strconv.Atoi(extracted[2])
		vY, _ := strconv.Atoi(extracted[3])
		robot := Robot{
			x:  x,
			y:  y,
			vX: vX,
			vY: vY,
		}
		robots = append(robots, robot)
	}

	robotsCopy := make([]Robot, len(robots))
	copy(robotsCopy, robots)
	moveRobots(robots, 100, spaceX, spaceY)
	quadrants := countRobotsinQuadrants(robots, spaceX, spaceY)
	mul := 1
	for _, q := range quadrants {
		if q > 0 {
			mul *= q
		}
	}
	fmt.Println(mul)

	xMasTreeTime := findXmasTree(robotsCopy, spaceX, spaceY)
	fmt.Println(xMasTreeTime)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
