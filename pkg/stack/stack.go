package stack

import "errors"

type Stack interface {
	Push(item interface{})
	Pop() (interface{}, error)
	Empty() bool
	Peek() (interface{}, error)
}

type stack struct {
	items []interface{}
	tos   int
}

// New creates a new stack
func New() *stack {
	return &stack{tos: -1}
}

// Push adds an element on the top of the stack
func (s *stack) Push(item interface{}) {
	s.items = append(s.items, item)
	s.tos++
}

// Pop removes the topmost element from the stack.
// If the stack is empty an error is returned.
func (s *stack) Pop() (interface{}, error) {
	if s.tos == -1 {
		return nil, errors.New("pop from empty stack")
	}
	item := s.items[s.tos]
	s.items = s.items[0:s.tos]
	s.tos--
	return item, nil
}

// Empty return true if the stack has no element, else false.
func (s *stack) Empty() bool {
	return s.tos == -1
}

// Peek returns the topmost element without removing it.
// If the stack is empty an error is returned.
func (s *stack) Peek() (interface{}, error) {
	if s.tos == -1 {
		return nil, errors.New("peek empty stack")
	}
	return s.items[s.tos], nil
}
