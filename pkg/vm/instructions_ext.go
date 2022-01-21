package vm

import (
	"fmt"
	"strings"

	"terhaak.de/imp/pkg/stack"
)

type PushStr string

func (inst PushStr) Exec(st stack.Stack, vm Machine) error {
	st.Push(string(inst))
	return nil
}

type ConcatStr struct{}

func (inst ConcatStr) Exec(st stack.Stack, vm Machine) error {
	values, err := popStrings(st, 2)
	if err == nil {
		st.Push(values[0] + values[1])
	}
	return err
}

type FormatStr string

func (inst FormatStr) Exec(st stack.Stack, vm Machine) error {
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

func (inst LengthStr) Exec(st stack.Stack, vm Machine) error {
	values, err := popStrings(st, 1)
	if err == nil {
		st.Push(len(values[0]))
	}
	return err
}

func (inst ConcatStr) String() string { return "cat" }
func (inst LengthStr) String() string { return "len" }

func (inst PushStr) String() string   { return fmt.Sprintf("psh \"%s\"", string(inst)) }
func (inst FormatStr) String() string { return fmt.Sprintf("fmt \"%s\"", string(inst)) }
