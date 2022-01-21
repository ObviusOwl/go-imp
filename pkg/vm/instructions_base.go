package vm

import (
	"fmt"

	"terhaak.de/imp/pkg/stack"
)

//
// control flow instructions
//

type Jump Label

func (inst Jump) Exec(vm Runner, st stack.Stack, mem Memory) error {
	return vm.Jump(Label(inst))
}

type JumpNonZero Label

func (inst JumpNonZero) Exec(vm Runner, st stack.Stack, mem Memory) error {
	values, err := popInts(st, 1)
	if err == nil && IntToBool(values[0]) {
		return vm.Jump(Label(inst))
	}
	return err
}

type JumpZero Label

func (inst JumpZero) Exec(vm Runner, st stack.Stack, mem Memory) error {
	values, err := popInts(st, 1)
	if err == nil && !IntToBool(values[0]) {
		return vm.Jump(Label(inst))
	}
	return err
}

type Stop struct{}

func (inst Stop) Exec(vm Runner, st stack.Stack, mem Memory) error {
	vm.Stop()
	return nil
}

//
// Arithmetic instructions
//

type Add struct{}

func (inst Add) Exec(vm Runner, st stack.Stack, mem Memory) error {
	return stackIntReduce(st, func(a, b int) (DataValue, error) { return a + b, nil })
}

type Minus struct{}

func (inst Minus) Exec(vm Runner, st stack.Stack, mem Memory) error {
	return stackIntReduce(st, func(a, b int) (DataValue, error) { return a - b, nil })
}

type Div struct{}

func (inst Div) Exec(vm Runner, st stack.Stack, mem Memory) error {
	return stackIntReduce(st, func(a, b int) (DataValue, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	})
}

type Mult struct{}

func (inst Mult) Exec(vm Runner, st stack.Stack, mem Memory) error {
	return stackIntReduce(st, func(a, b int) (DataValue, error) { return a * b, nil })
}

//
// Logic instructions
//

type Equal struct{}

func (inst Equal) Exec(vm Runner, st stack.Stack, mem Memory) error {
	op1, err1 := st.Pop()
	op2, err2 := st.Pop()

	if err1 == nil && err2 == nil {
		st.Push(BoolToInt(op1 == op2))
	} else if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	}
	return nil
}

type Lesser struct{}

func (inst Lesser) Exec(vm Runner, st stack.Stack, mem Memory) error {
	return stackIntReduce(st, func(a, b int) (DataValue, error) { return BoolToInt(a < b), nil })
}

type Greater struct{}

func (inst Greater) Exec(vm Runner, st stack.Stack, mem Memory) error {
	return stackIntReduce(st, func(a, b int) (DataValue, error) { return BoolToInt(a > b), nil })
}

//
// data move instructions
//

type PushInt int

func (inst PushInt) Exec(vm Runner, st stack.Stack, mem Memory) error {
	st.Push(int(inst))
	return nil
}

type StoreMemory int

func (inst StoreMemory) Exec(vm Runner, st stack.Stack, mem Memory) error {
	if value, err := st.Pop(); err != nil {
		return err
	} else {
		mem.Store(int(inst), value)
	}
	return nil
}

type LoadMemory int

func (inst LoadMemory) Exec(vm Runner, st stack.Stack, mem Memory) error {
	st.Push(mem.Load(int(inst)))
	return nil
}

type Output int

func (inst Output) Exec(vm Runner, st stack.Stack, mem Memory) error {
	fmt.Printf("%v\n", mem.Load(int(inst)))
	return nil
}

//
// Mnemonics
//

func (inst Label) String() string       { return fmt.Sprintf("lab %d", int(inst)) }
func (inst Jump) String() string        { return fmt.Sprintf("jmp %d", int(inst)) }
func (inst JumpNonZero) String() string { return fmt.Sprintf("jnz %d", int(inst)) }
func (inst JumpZero) String() string    { return fmt.Sprintf("jez %d", int(inst)) }
func (inst Stop) String() string        { return "stp" }

func (inst Add) String() string   { return "add" }
func (inst Minus) String() string { return "min" }
func (inst Div) String() string   { return "div" }
func (inst Mult) String() string  { return "mul" }

func (inst Equal) String() string   { return "eql" }
func (inst Greater) String() string { return "gtt" }
func (inst Lesser) String() string  { return "ltt" }

func (inst PushInt) String() string     { return fmt.Sprintf("psh %d", int(inst)) }
func (inst StoreMemory) String() string { return fmt.Sprintf("stm %d", int(inst)) }
func (inst LoadMemory) String() string  { return fmt.Sprintf("ldm %d", int(inst)) }
func (inst Output) String() string      { return fmt.Sprintf("out %d", int(inst)) }
