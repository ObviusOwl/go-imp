package vm

import (
	"fmt"
	"strings"

	"terhaak.de/imp/pkg/stack"
)

type PushStr string

func (inst PushStr) Exec(vm Runner, st stack.Stack, mem Memory) error {
	st.Push(string(inst))
	return nil
}

type ConcatStr struct{}

func (inst ConcatStr) Exec(vm Runner, st stack.Stack, mem Memory) error {
	values, err := popStrings(st, 2)
	if err == nil {
		st.Push(values[0] + values[1])
	}
	return err
}

type FormatStr string

func (inst FormatStr) Exec(vm Runner, st stack.Stack, mem Memory) error {
	format := string(inst)
	// count how many % we have, minus the escaped ones
	argc := strings.Count(format, "%") - 2*strings.Count(format, "%%")

	// dont cast anything
	values := make([]interface{}, argc)
	err := stack.Process(st, func(itcount int, item interface{}) (bool, error) {
		values[itcount-1] = item
		return (itcount < argc), nil
	})

	if err == nil {
		st.Push(fmt.Sprintf(format, values...))
	}
	return err
}

type LengthStr struct{}

func (inst LengthStr) Exec(vm Runner, st stack.Stack, mem Memory) error {
	values, err := popStrings(st, 1)
	if err == nil {
		st.Push(len(values[0]))
	}
	return err
}

func (inst ConcatStr) String() string { return "cat" }
func (inst LengthStr) String() string { return "len" }

func (inst PushStr) String() string   { return fmt.Sprintf("str \"%s\"", string(inst)) }
func (inst FormatStr) String() string { return fmt.Sprintf("fmt \"%s\"", string(inst)) }
