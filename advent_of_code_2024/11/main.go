package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func blink(stone int) []int {
	if stone == 0 {
		return []int{1}
	} else if len(strconv.Itoa(stone))%2 == 0 {
		stoneName := strconv.Itoa(stone)
		fistHalf, _ := strconv.Atoi(stoneName[:len(stoneName)/2])
		secondHalf, _ := strconv.Atoi(stoneName[len(stoneName)/2:])
		return []int{fistHalf, secondHalf}
	} else {
		return []int{stone * 2024}
	}
}

func blinkStoneNTimes(stone int, cache *sync.Map, originalBlinks int, nBlinks int, mu *sync.Mutex) int {
	if v, ok := cache.Load(stone); ok {
		if stoneResults, ok := v.([]int); ok {
			if stoneResults[nBlinks-1] != 0 {
				return stoneResults[nBlinks-1]
			}
		}
	} else {
		cache.Store(stone, make([]int, originalBlinks))
	}

	if nBlinks == 1 {
		nStones := len(blink(stone))
		mu.Lock()
		v, _ := cache.Load(stone)
		stoneResults, _ := v.([]int)
		stoneResults[0] = nStones
		cache.Store(stone, stoneResults)
		mu.Unlock()
		return nStones
	}

	nStones := 0
	for _, stone := range blink(stone) {
		nStones += blinkStoneNTimes(stone, cache, originalBlinks, nBlinks-1, mu)
	}
	mu.Lock()
	v, _ := cache.Load(stone)
	stoneResults, _ := v.([]int)
	stoneResults[nBlinks-1] = nStones
	mu.Unlock()
	return nStones
}

func reduceBlinks(nStones []int) int {
	totalStones := 0
	for _, nStone := range nStones {
		totalStones += nStone
	}
	return totalStones
}

func runBlinks(stones []int, nBlinks int) int {
	var mu sync.Mutex
	var stonesCache sync.Map
	totalStones := 0
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	numWorkers := runtime.NumCPU() - 1

	var wg sync.WaitGroup
	taskQueue := make(chan int, len(stones))
	results := make([]int, len(stones))

	// Create a worker pool
	for j := 0; j < numWorkers; j++ {
		go func() {
			for index := range taskQueue {
				results[index] = blinkStoneNTimes(stones[index], &stonesCache, nBlinks, nBlinks, &mu)
				wg.Done()
			}
		}()
	}

	// Add tasks to the taskQueue
	for index := range stones {
		wg.Add(1)
		taskQueue <- index
	}

	close(taskQueue)
	wg.Wait()
	totalStones = reduceBlinks(results)
	return totalStones
}

func main() {
	// read data from test_input.txt
	file, err := os.Open("./11/input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var stones []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// split line by space
		stonesNames := strings.Split(line, " ")
		for _, stoneName := range stonesNames {
			stoneNumber, _ := strconv.Atoi(stoneName)
			stones = append(stones, stoneNumber)
		}
	}
	fmt.Println(stones)
	part1Stones := runBlinks(stones, 25)
	//fmt.Println(part1Stones)
	fmt.Println(part1Stones)

	part2Stones := runBlinks(stones, 75)
	//fmt.Println(part2Stones)
	fmt.Println(part2Stones)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
