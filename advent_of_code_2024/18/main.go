package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// const w, h, byteCount = 7, 7, 12
const w, h, byteCount = 71, 71, 1024

type Point struct {
	x, y int
}

type Node struct {
	point    Point
	distance int
}

type State struct {
	x, y      int
	direction int      // 0: East, 1: South, 2: West, 3: North
	cost      int      // Total score so far
	path      []string // Path history (sequence of moves)
}

type PriorityQueue []Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].distance < pq[j].distance }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Node))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func isValid(x, y, rows, cols int, grid [][]byte, visited [][]bool) bool {
	return x >= 0 && x < rows && y >= 0 && y < cols && grid[x][y] != '#' && !visited[x][y]
}

func dijkstra(grid [][]byte, startPoint, endPoint Point) int {
	rows, cols := len(grid), len(grid[0])

	directions := []Point{
		{0, 1}, {1, 0}, {0, -1}, {-1, 0}, // Right, Down, Left, Up
	}

	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, Node{startPoint, 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(Node)
		cx, cy := current.point.x, current.point.y

		if visited[cx][cy] {
			continue
		}

		visited[cx][cy] = true

		// If we reach the end, return the distance
		if cx == endPoint.x && cy == endPoint.y {
			return current.distance
		}

		// Explore neighbors
		for _, d := range directions {
			nx, ny := cx+d.x, cy+d.y
			if isValid(nx, ny, rows, cols, grid, visited) {
				heap.Push(pq, Node{Point{nx, ny}, current.distance + 1})
			}
		}
	}

	// If there's no path to the end, return -1
	return -1
}

func findBlockingByte(memoryMap [][]byte, startPoint, endPoint Point, fallingBytes []Point) Point {
	for i := byteCount; i < len(fallingBytes); i++ {
		memoryMap[fallingBytes[i].y][fallingBytes[i].x] = '#'
		distance := dijkstra(memoryMap, startPoint, endPoint)
		if distance == -1 {
			return fallingBytes[i]
		}
	}
	return Point{-1, -1}
}

func main() {
	file, err := os.Open("./18/input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	memoryMap := make([][]byte, h)
	for i := 0; i < h; i++ {
		memoryMap[i] = make([]byte, w)
		for j := 0; j < w; j++ {
			memoryMap[i][j] = '.'
		}
	}
	startPosition := Point{0, 0}
	memoryMap[startPosition.x][startPosition.y] = 'S'
	endPosition := Point{h - 1, w - 1}
	memoryMap[endPosition.x][endPosition.y] = 'E'

	fallingBytes := make([]Point, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		fallingBytes = append(fallingBytes, Point{x, y})
	}

	for i := 0; i < byteCount; i++ {
		memoryMap[fallingBytes[i].y][fallingBytes[i].x] = '#'
	}

	for _, row := range memoryMap {
		fmt.Println(string(row[:]))
	}

	distance := dijkstra(memoryMap, startPosition, endPosition)
	fmt.Println("Distance:", distance)

	breakingPoint := findBlockingByte(memoryMap, startPosition, endPosition, fallingBytes)
	fmt.Println("Breaking Point:", breakingPoint)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
