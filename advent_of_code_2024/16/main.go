package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strings"
)

type State struct {
	x, y      int
	direction int      // 0: East, 1: South, 2: West, 3: North
	cost      int      // Total score so far
	path      []string // Path history (sequence of moves)
}

// Directions: East (0), South (1), West (2), North (3)
var directions = []struct {
	dx, dy int
	symbol byte
}{
	{0, 1, '>'},  // East
	{1, 0, 'v'},  // South
	{0, -1, '<'}, // West
	{-1, 0, '^'}, // North
}

type PriorityQueue []State

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(State))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func Dijkstra(maze [][]byte, startX, startY, endX, endY int) (int, []State) {
	// Priority queue for Dijkstra
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, State{startX, startY, 0, 0, []string{}}) // Start at (startX, startY) facing East with cost 0

	// Track minimum cost to reach each state (x, y, direction)
	visited := make(map[[3]int]int)

	// List to store all paths with the minimum cost to the end state
	var optimalPaths []State
	minCost := -1

	for pq.Len() > 0 {
		current := heap.Pop(pq).(State)

		// Check if we reached the end
		if maze[current.x][current.y] == 'E' {
			// If this is the first time reaching the end, record the cost
			if minCost == -1 {
				minCost = current.cost
			}
			// If this state has the same minimum cost, add it to the list
			if current.cost == minCost {
				optimalPaths = append(optimalPaths, current)
			}
			continue
		}

		// Skip if we've seen this state with a lower cost
		key := [3]int{current.x, current.y, current.direction}
		if prevCost, found := visited[key]; found && current.cost > prevCost {
			continue
		}
		visited[key] = current.cost

		// Move Forward
		nx := current.x + directions[current.direction].dx
		ny := current.y + directions[current.direction].dy
		if maze[nx][ny] != '#' { // Valid move
			newPath := append([]string{}, current.path...) // Copy current path
			newPath = append(newPath, "F")
			heap.Push(pq, State{nx, ny, current.direction, current.cost + 1, newPath})
		}

		// Turn Clockwise
		newPathR := append([]string{}, current.path...)
		newPathR = append(newPathR, "R")
		heap.Push(pq, State{current.x, current.y, (current.direction + 1) % 4, current.cost + 1000, newPathR})

		// Turn Counterclockwise
		newPathL := append([]string{}, current.path...)
		newPathL = append(newPathL, "L")
		heap.Push(pq, State{current.x, current.y, (current.direction + 3) % 4, current.cost + 1000, newPathL})
	}

	return minCost, optimalPaths
}

func countSeats(startX, startY int, pathStates []State) int {
	uniquePoints := make(map[[2]int]bool) // Map to store unique positions (x, y)

	for _, pathState := range pathStates {
		path := pathState.path
		x, y, direction := startX, startY, 0 // Start position and facing East
		uniquePoints[[2]int{x, y}] = true    // Mark start position as visited

		for _, move := range path {
			switch move {
			case "F": // Move Forward
				x += directions[direction].dx
				y += directions[direction].dy
				uniquePoints[[2]int{x, y}] = true
			case "R": // Turn Clockwise
				direction = (direction + 1) % 4
			case "L": // Turn Counterclockwise
				direction = (direction + 3) % 4 // Equivalent to (direction - 1 + 4) % 4
			}
		}
	}

	return len(uniquePoints) // Return the count of unique points
}

func printPathOnMap(maze [][]byte, startX, startY int, path []string) {
	// Start at the given position, facing East (direction = 0)
	x, y, direction := startX, startY, 0

	// Replace the starting position with the initial direction
	maze[x][y] = directions[direction].symbol

	// Process the path
	for _, move := range path {
		if move == "F" { // Move Forward
			x += directions[direction].dx
			y += directions[direction].dy
			if maze[x][y] != '#' { // Only mark valid positions
				maze[x][y] = directions[direction].symbol
			}
		} else if move == "R" { // Rotate Clockwise
			direction = (direction + 1) % 4
		} else if move == "L" { // Rotate Counterclockwise
			direction = (direction + 3) % 4 // (direction - 1 + 4) % 4
		}
	}

	// Print the updated maze
	for _, row := range maze {
		fmt.Println(string(row))
	}
}

func main() {
	file, err := os.Open("./16/input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	raceMap := make([][]byte, 0)
	startPosition := [2]int{0, 0}
	endPosition := [2]int{0, 0}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		raceMap = append(raceMap, []byte(line))
		startIndex := strings.Index(line, "S")
		if startIndex != -1 {
			startPosition[0] = len(raceMap) - 1
			startPosition[1] = startIndex
		}
		endIndex := strings.Index(line, "E")
		if endIndex != -1 {
			endPosition[0] = len(raceMap) - 1
			endPosition[1] = endIndex
		}
	}

	for _, row := range raceMap {
		fmt.Println(string(row))
	}

	cost, paths := Dijkstra(raceMap, startPosition[0], startPosition[1], endPosition[0], endPosition[1])
	fmt.Println("Cost:", cost)
	fmt.Println("Paths:", paths)
	//fmt.Println("Path:", path)

	// Print the path on the map
	printPathOnMap(raceMap, startPosition[0], startPosition[1], paths[0].path)
	spectatorSeats := countSeats(startPosition[0], startPosition[1], paths)
	fmt.Println("Spectator Seats:", spectatorSeats)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
