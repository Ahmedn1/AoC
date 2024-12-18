package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var registers []uint64
var instructionPointer int
var output []int

func getComboOperand(operand uint64) uint64 {
	if operand <= 3 {
		return operand
	}
	return registers[operand-4]
}

func adv(operand uint64) {
	numerator := registers[0]
	denom := uint64(math.Pow(2, float64(getComboOperand(operand))))
	registers[0] = numerator / denom
}

func bxl(operand uint64) {
	registers[1] = registers[1] ^ operand
}

func bst(operand uint64) {
	combo := getComboOperand(operand)
	registers[1] = combo % 8
}

func jnz(operand uint64) {
	if registers[0] != 0 {
		instructionPointer = int(operand - 2)
		//instructionPointerIncrease = 1
	}
}

func bxc(operand uint64) {
	registers[1] = registers[1] ^ registers[2]
}

func out(operand uint64) {
	combo := getComboOperand(operand)
	//fmt.Print(combo%8, ",")
	output = append(output, int(combo%8))
}

func bdv(operand uint64) {
	numerator := registers[0]
	denom := uint64(math.Pow(2, float64(getComboOperand(operand))))
	registers[1] = numerator / denom
}

func cdv(operand uint64) {
	numerator := registers[0]
	denom := uint64(math.Pow(2, float64(getComboOperand(operand))))
	registers[2] = numerator / denom
}

var instructions = map[int]func(uint64){
	0: adv,
	1: bxl,
	2: bst,
	3: jnz,
	4: bxc,
	5: out,
	6: bdv,
	7: cdv,
}

func runProgram(program []int) {
	for instructionPointer < len(program) {
		//fmt.Printf("A: %08b\n", registers[0])
		opcode := program[instructionPointer]
		operand := program[instructionPointer+1]
		//fmt.Println("Instruction:", opcode, operand)
		instructions[opcode](uint64(operand))
		instructionPointer += 2
	}
}

// I was able to complete the missing piece of this logic by the help of Gravitar64 on the AoC subreddit[https://www.reddit.com/r/adventofcode/comments/1hg38ah/comment/m2k2ka5/]
func findA(program []int, a, b, c uint64, prgPos int) uint64 {
	if prgPos < 0 || prgPos >= len(program) {
		return a
	}

	// Try all possible values for the last 3 bits (i from 0 to 7)
	for i := 0; i < 8; i++ {
		// Run the program with the updated value of A (shifted by 3 bits)
		registers = []uint64{a*8 + uint64(i), b, c}

		// Reset the instruction pointer and output
		instructionPointer = 0
		output = []int{}

		runProgram(program)
		firstDigitOut := output[0]

		// If the output matches the program's value at the current position, recurse
		if firstDigitOut == program[prgPos] {
			// Recursively call findA with the updated value of 'a' and the next position
			e := findA(program, a*8+uint64(i), b, c, prgPos-1)
			if e != 0 {
				return e
			}
		}
	}

	// Return 0 if no valid value of 'a' is found
	return 0
}

func main() {
	file, err := os.Open("./17/test_input_2.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	program := make([]int, 0)
	registers = []uint64{}
	programStart := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			programStart = true
			continue
		}
		if programStart {
			parts := strings.Split(line[9:], ",")
			for i := 0; i < len(parts); i++ {
				val, _ := strconv.Atoi(parts[i])
				program = append(program, val)
			}
		} else {
			register, _ := strconv.Atoi(line[12:])
			registers = append(registers, uint64(register))
		}
	}
	instructionPointer = 0

	fmt.Println(program)
	runProgram(program)
	fmt.Println(findA(program, 0, registers[1], registers[2], len(program)-1))

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
