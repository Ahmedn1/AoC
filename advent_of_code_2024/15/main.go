package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
)

type Box struct {
	positions [][2]int
}

var boxValues = []byte{'O', '[', ']'}
var boxShapePositionLenMap = map[int][]byte{1: {'O'}, 2: {'[', ']'}}

func canBoxMoveTentative(warehouseMap [][]byte, box Box, direction string) bool {
	switch direction {
	case "^":
		for _, position := range box.positions {
			if position[0] == 0 || warehouseMap[position[0]-1][position[1]] == '#' {
				return false
			}
		}
	case ">":
		position := box.positions[len(box.positions)-1]
		if position[1] == len(warehouseMap[0])-1 || warehouseMap[position[0]][position[1]+1] == '#' {
			return false
		}
	case "<":
		position := box.positions[0]
		if position[1] == 0 || warehouseMap[position[0]][position[1]-1] == '#' {
			return false
		}
	case "v":
		for _, position := range box.positions {
			if position[0] == len(warehouseMap)-1 || warehouseMap[position[0]+1][position[1]] == '#' {
				return false
			}
		}
	}
	return true
}

func canBoxMove(warehouseMap [][]byte, box Box, direction string) bool {
	switch direction {
	case "^":
		for _, position := range box.positions {
			if position[0] == 0 || warehouseMap[position[0]-1][position[1]] != '.' {
				return false
			}
		}
	case ">":
		position := box.positions[len(box.positions)-1]
		if position[1] == len(warehouseMap[0])-1 || warehouseMap[position[0]][position[1]+1] != '.' {
			return false
		}
	case "<":
		position := box.positions[0]
		if position[1] == 0 || warehouseMap[position[0]][position[1]-1] != '.' {
			return false
		}
	case "v":
		for _, position := range box.positions {
			if position[0] == len(warehouseMap)-1 || warehouseMap[position[0]+1][position[1]] != '.' {
				return false
			}
		}
	}
	return true
}

func checkBoxExists(boxes []Box, box Box) bool {
	for _, b := range boxes {
		if len(b.positions) == len(box.positions) {
			equal := true
			for i, position := range b.positions {
				if position[0] != box.positions[i][0] || position[1] != box.positions[i][1] {
					equal = false
					break
				}
			}
			if equal {
				return true
			}
		}
	}
	return false
}

func addBoxToList(boxes []Box, box Box) []Box {
	sort.Slice(box.positions, func(i, j int) bool {
		if box.positions[i][0] == box.positions[j][0] {
			return box.positions[i][1] < box.positions[j][1]
		}
		return box.positions[i][0] < box.positions[j][0]
	})
	if !checkBoxExists(boxes, box) {
		boxes = append(boxes, box)
	}
	return boxes
}

func getLinkedBoxes(warehouseMap [][]byte, box Box, direction string) []Box {
	linkedBoxes := make([]Box, 0)
	switch direction {
	case "^":
		for _, position := range box.positions {
			boxPart := warehouseMap[position[0]-1][position[1]]
			if position[0] > 0 && slices.Contains(boxValues, boxPart) {
				newBox := Box{}
				switch boxPart {
				case 'O':
					newBox.positions = append(newBox.positions, [2]int{position[0] - 1, position[1]})
				case '[':
					newBox.positions = append(newBox.positions, [2]int{position[0] - 1, position[1]})
					newBox.positions = append(newBox.positions, [2]int{position[0] - 1, position[1] + 1})
				case ']':
					newBox.positions = append(newBox.positions, [2]int{position[0] - 1, position[1]})
					newBox.positions = append(newBox.positions, [2]int{position[0] - 1, position[1] - 1})
				}
				linkedBoxes = addBoxToList(linkedBoxes, newBox)
			}
		}
	case ">":
		position := box.positions[len(box.positions)-1]
		if position[1] < len(warehouseMap[0])-1 && slices.Contains(boxValues, warehouseMap[position[0]][position[1]+1]) {
			boxPart := warehouseMap[position[0]][position[1]+1]
			if position[1] < len(warehouseMap[0])-1 && slices.Contains(boxValues, boxPart) {
				newBox := Box{}
				switch boxPart {
				case 'O':
					newBox.positions = append(newBox.positions, [2]int{position[0], position[1] + 1})
				case '[':
					newBox.positions = append(newBox.positions, [2]int{position[0], position[1] + 1})
					newBox.positions = append(newBox.positions, [2]int{position[0], position[1] + 2})
				}
				linkedBoxes = addBoxToList(linkedBoxes, newBox)
			}
		}
	case "<":
		position := box.positions[0]
		if position[1] > 0 && slices.Contains(boxValues, warehouseMap[position[0]][position[1]-1]) {
			boxPart := warehouseMap[position[0]][position[1]-1]
			if position[1] > 0 && slices.Contains(boxValues, boxPart) {
				newBox := Box{}
				switch boxPart {
				case 'O':
					newBox.positions = append(newBox.positions, [2]int{position[0], position[1] - 1})
				case ']':
					newBox.positions = append(newBox.positions, [2]int{position[0], position[1] - 1})
					newBox.positions = append(newBox.positions, [2]int{position[0], position[1] - 2})
				}
				linkedBoxes = addBoxToList(linkedBoxes, newBox)
			}
		}
	case "v":
		for _, position := range box.positions {
			boxPart := warehouseMap[position[0]+1][position[1]]
			if position[0] < len(warehouseMap)-1 && slices.Contains(boxValues, boxPart) {
				newBox := Box{}
				switch boxPart {
				case 'O':
					newBox.positions = append(newBox.positions, [2]int{position[0] + 1, position[1]})
				case '[':
					newBox.positions = append(newBox.positions, [2]int{position[0] + 1, position[1]})
					newBox.positions = append(newBox.positions, [2]int{position[0] + 1, position[1] + 1})
				case ']':
					newBox.positions = append(newBox.positions, [2]int{position[0] + 1, position[1]})
					newBox.positions = append(newBox.positions, [2]int{position[0] + 1, position[1] - 1})
				}
				linkedBoxes = addBoxToList(linkedBoxes, newBox)
			}
		}
	}

	return linkedBoxes
}

func findAllLinkedBoxes(warehouseMap [][]byte, box Box, direction string) []Box {
	allLinkedBoxes := make([]Box, 0)
	linkedBoxes := getLinkedBoxes(warehouseMap, box, direction)
	for _, linkedBox := range linkedBoxes {
		allLinkedBoxes = append(allLinkedBoxes, linkedBox)
		linkedBoxes2 := findAllLinkedBoxes(warehouseMap, linkedBox, direction)
		for _, linkedBox2 := range linkedBoxes2 {
			allLinkedBoxes = append(allLinkedBoxes, linkedBox2)
		}
	}
	return allLinkedBoxes
}

func moveBox(warehouseMap [][]byte, box Box, direction string) Box {
	allLinkedBoxes := findAllLinkedBoxes(warehouseMap, box, direction)
	allLinkedBoxes = append(allLinkedBoxes, box)
	for _, linkedBox := range allLinkedBoxes {
		if !canBoxMoveTentative(warehouseMap, linkedBox, direction) {
			return box
		}
	}

	sort.Slice(allLinkedBoxes, func(i, j int) bool {
		if allLinkedBoxes[i].positions[0][0] == allLinkedBoxes[j].positions[0][0] {
			return allLinkedBoxes[i].positions[0][1] < allLinkedBoxes[j].positions[0][1]
		}
		return allLinkedBoxes[i].positions[0][0] < allLinkedBoxes[j].positions[0][0]
	})

	switch direction {
	case "^":
		for _, linkedBox := range allLinkedBoxes {
			for i, position := range linkedBox.positions {
				warehouseMap[position[0]-1][position[1]] = boxShapePositionLenMap[len(linkedBox.positions)][i]
				warehouseMap[position[0]][position[1]] = '.'
				linkedBox.positions[i][0]--
			}
		}
	case ">":
		for j := len(allLinkedBoxes) - 1; j >= 0; j-- {
			linkedBox := allLinkedBoxes[j]
			for i := len(linkedBox.positions) - 1; i >= 0; i-- {
				position := linkedBox.positions[i]
				warehouseMap[position[0]][position[1]+1] = boxShapePositionLenMap[len(linkedBox.positions)][i]
				warehouseMap[position[0]][position[1]] = '.'
				linkedBox.positions[i][1]++
			}
		}
	case "<":
		for _, linkedBox := range allLinkedBoxes {
			for i, position := range linkedBox.positions {
				warehouseMap[position[0]][position[1]-1] = boxShapePositionLenMap[len(linkedBox.positions)][i]
				warehouseMap[position[0]][position[1]] = '.'
				linkedBox.positions[i][1]--
			}
		}
	case "v":
		for j := len(allLinkedBoxes) - 1; j >= 0; j-- {
			linkedBox := allLinkedBoxes[j]
			for i, position := range linkedBox.positions {
				warehouseMap[position[0]+1][position[1]] = boxShapePositionLenMap[len(linkedBox.positions)][i]
				warehouseMap[position[0]][position[1]] = '.'
				linkedBox.positions[i][0]++
			}
		}
	}
	return box
}

func moveRobot(warehouseMap [][]byte, robotMovements string, initialRobotPosition [2]int) Box {
	robotBox := Box{}
	robotBox.positions = append(robotBox.positions, [2]int{initialRobotPosition[0], initialRobotPosition[1]})
	for _, movement := range robotMovements {
		linkedBoxes := getLinkedBoxes(warehouseMap, robotBox, string(movement))
		if len(linkedBoxes) > 0 {
			for _, linkedBox := range linkedBoxes {
				linkedBox = moveBox(warehouseMap, linkedBox, string(movement))
			}
		}
		switch movement {
		case '^':
			if canBoxMove(warehouseMap, robotBox, string(movement)) {
				warehouseMap[robotBox.positions[0][0]-1][robotBox.positions[0][1]] = '@'
				warehouseMap[robotBox.positions[0][0]][robotBox.positions[0][1]] = '.'
				robotBox.positions[0][0]--
			}
		case '>':
			if canBoxMove(warehouseMap, robotBox, string(movement)) {
				warehouseMap[robotBox.positions[0][0]][robotBox.positions[0][1]+1] = '@'
				warehouseMap[robotBox.positions[0][0]][robotBox.positions[0][1]] = '.'
				robotBox.positions[0][1]++
			}
		case '<':
			if canBoxMove(warehouseMap, robotBox, string(movement)) {
				warehouseMap[robotBox.positions[0][0]][robotBox.positions[0][1]-1] = '@'
				warehouseMap[robotBox.positions[0][0]][robotBox.positions[0][1]] = '.'
				robotBox.positions[0][1]--
			}
		case 'v':
			if canBoxMove(warehouseMap, robotBox, string(movement)) {
				warehouseMap[robotBox.positions[0][0]+1][robotBox.positions[0][1]] = '@'
				warehouseMap[robotBox.positions[0][0]][robotBox.positions[0][1]] = '.'
				robotBox.positions[0][0]++
			}
		}
	}
	return robotBox
}

func calculateGPS(warehouseMap [][]byte) int {
	gps := 0
	for i := 0; i < len(warehouseMap); i++ {
		for j := 0; j < len(warehouseMap[0]); j++ {
			if warehouseMap[i][j] == 'O' || warehouseMap[i][j] == '[' {
				gps += 100*i + j
			}
		}
	}
	return gps
}

func main() {
	file, err := os.Open("./15/input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	warehouseMap := make([][]byte, 0)
	warehouseMap2 := make([][]byte, 0)
	initialRobot1Position := [2]int{0, 0}
	initialRobot2Position := [2]int{0, 0}
	mapFinished := false
	robotMovements := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			mapFinished = true
			continue
		}
		if mapFinished {
			robotMovements += line
			continue
		}
		warehouseMap = append(warehouseMap, []byte(line))
		robotIndex := strings.Index(line, "@")
		if robotIndex != -1 {
			initialRobot1Position[0] = len(warehouseMap) - 1
			initialRobot1Position[1] = robotIndex
		}

		transformedLine := strings.ReplaceAll(line, ".", "..")
		transformedLine = strings.ReplaceAll(transformedLine, "#", "##")
		transformedLine = strings.ReplaceAll(transformedLine, "@", "@.")
		transformedLine = strings.ReplaceAll(transformedLine, "O", "[]")
		warehouseMap2 = append(warehouseMap2, []byte(transformedLine))
		robotIndex = strings.Index(transformedLine, "@")
		if robotIndex != -1 {
			initialRobot2Position[0] = len(warehouseMap2) - 1
			initialRobot2Position[1] = robotIndex
		}
	}
	_ = moveRobot(warehouseMap, robotMovements, initialRobot1Position)
	fmt.Println(calculateGPS(warehouseMap))

	_ = moveRobot(warehouseMap2, robotMovements, initialRobot2Position)
	fmt.Println(calculateGPS(warehouseMap2))

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
