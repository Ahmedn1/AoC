package main

import "fmt"

// Stack represents a stack data structure
type Stack struct {
	elements []int
}

// Push adds an element to the top of the stack
func (s *Stack) Push(value int) {
	s.elements = append(s.elements, value)
}

// Pop removes and returns the top element of the stack
func (s *Stack) Pop() (int, error) {
	if len(s.elements) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}

	// Get the last element
	value := s.elements[len(s.elements)-1]

	// Remove the last element
	s.elements = s.elements[:len(s.elements)-1]

	return value, nil
}

// Peek returns the top element without removing it
func (s *Stack) Peek() (int, error) {
	if len(s.elements) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}

	return s.elements[len(s.elements)-1], nil
}

// IsEmpty checks if the stack is empty
func (s *Stack) IsEmpty() bool {
	return len(s.elements) == 0
}

// Size returns the number of elements in the stack
func (s *Stack) Size() int {
	return len(s.elements)
}
