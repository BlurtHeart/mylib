package stack

import (
	"errors"
)

type stack []interface{}

func NewStack() *Stack {
	return &stack{}
}

func (s *stack) Push(data interface{}) {
	*s = append(*s, data)
}

func (s *stack) Pop() (interface{}, error) {
	if s.Len() == 0 {
		return nil, errors.New("empty stack")
	}
	data := (*s)[s.Len()-1]
	*s = (*s)[:s.Len()-1]
	return data, nil
}

func (s *stack) Top() (interface{}, error) {
	if s.Len() == 0 {
		return nil, errors.New("empty stack")
	}
	return (*s)[s.Len()-1], nil
}

func (s *stack) Len() int {
	return len(*s)
}
