package vm

import (
	"fmt"

	"terhaak.de/imp/pkg/stack"
)

func popInts(st stack.Stack, count int) ([]int, error) {
	values := make([]int, count)
	err := stack.Process(st, func(itcount int, item interface{}) (bool, error) {
		if value, ok := item.(int); ok {
			values[itcount-1] = value
		} else {
			return (itcount < count), fmt.Errorf("expected int from stack, got %v", item)
		}
		return (itcount < count), nil
	})
	return values, err
}

func popStrings(st stack.Stack, count int) ([]string, error) {
	values := make([]string, count)
	err := stack.Process(st, func(itcount int, item interface{}) (bool, error) {
		if value, ok := item.(string); ok {
			values[itcount-1] = value
		} else {
			return (itcount < count), fmt.Errorf("expected string from stack, got %v", item)
		}
		return (itcount < count), nil
	})
	return values, err
}

func stackIntReduce(st stack.Stack, f func(int, int) (DataValue, error)) error {
	var err error
	operands, err := popInts(st, 2)
	if err == nil {
		result, err := f(operands[0], operands[1])
		if err == nil {
			st.Push(result)
		}
	}
	return err
}

func BoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func IntToBool(value int) bool {
	return value != 0
}
