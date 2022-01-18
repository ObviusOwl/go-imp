package vm

import (
	"fmt"

	"terhaak.de/imp/pkg/stack"
)

func popInts(st stack.Stack, count int) ([]int, error) {
	var typeerr error = nil
	values := make([]int, count)

	stackerr := stack.Process(st, func(itcount int, item interface{}) bool {
		if value, ok := item.(int); ok {
			values[itcount-1] = value
		} else {
			typeerr = fmt.Errorf("expected int from stack, got %v", item)
		}
		return (itcount < count) && (typeerr == nil)
	})

	if stackerr != nil {
		return values, stackerr
	}
	return values, typeerr
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
