package vm

import (
	"fmt"

	"terhaak.de/imp/pkg/stack"
)

//
// control flow instructions
//

type Jump Label

func (inst Jump) Exec(st stack.Stack, vm Machine) error {
	return vm.Jump(Label(inst))
}

type JumpNonZero Label

func (inst JumpNonZero) Exec(st stack.Stack, vm Machine) error {
	values, err := popInts(st, 1)
	if err == nil && IntToBool(values[0]) {
		return vm.Jump(Label(inst))
	}
	return err
}

type JumpZero Label

func (inst JumpZero) Exec(st stack.Stack, vm Machine) error {
	values, err := popInts(st, 1)
	if err == nil && !IntToBool(values[0]) {
		return vm.Jump(Label(inst))
	}
	return err
}

//
// Arithmetic instructions
//

type Add struct{}

func (inst Add) Exec(st stack.Stack, vm Machine) error {
	return stackIntReduce(st, func(a, b int) (DataValue, error) { return a + b, nil })
}

type Minus struct{}

func (inst Minus) Exec(st stack.Stack, vm Machine) error {
	return stackIntReduce(st, func(a, b int) (DataValue, error) { return a - b, nil })
}

type Div struct{}

func (inst Div) Exec(st stack.Stack, vm Machine) error {
	return stackIntReduce(st, func(a, b int) (DataValue, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	})
}

type Mult struct{}

func (inst Mult) Exec(st stack.Stack, vm Machine) error {
	return stackIntReduce(st, func(a, b int) (DataValue, error) { return a * b, nil })
}

//
// Logic instructions
//

type Equal struct{}

func (inst Equal) Exec(st stack.Stack, vm Machine) error {
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

func (inst Lesser) Exec(st stack.Stack, vm Machine) error {
	return stackIntReduce(st, func(a, b int) (DataValue, error) { return BoolToInt(a < b), nil })
}

type Greater struct{}

func (inst Greater) Exec(st stack.Stack, vm Machine) error {
	return stackIntReduce(st, func(a, b int) (DataValue, error) { return BoolToInt(a > b), nil })
}

//
// data move instructions
//

type PushInt int

func (inst PushInt) Exec(st stack.Stack, vm Machine) error {
	st.Push(int(inst))
	return nil
}

type Asg int

func (inst Asg) Exec(st stack.Stack, vm Machine) error {
	if value, err := st.Pop(); err != nil {
		return err
	} else {
		vm.Store(int(inst), value)
	}
	return nil
}

type Deref int

func (inst Deref) Exec(st stack.Stack, vm Machine) error {
	st.Push(vm.Load(int(inst)))
	return nil
}

type Output int

func (inst Output) Exec(st stack.Stack, vm Machine) error {
	fmt.Printf("%v\n", vm.Load(int(inst)))
	return nil
}

//
// Mnemonics
//

func (inst Label) String() string       { return fmt.Sprintf("lab %d", int(inst)) }
func (inst Jump) String() string        { return fmt.Sprintf("jmp %d", int(inst)) }
func (inst JumpNonZero) String() string { return fmt.Sprintf("jnz %d", int(inst)) }
func (inst JumpZero) String() string    { return fmt.Sprintf("jez %d", int(inst)) }

func (inst Add) String() string   { return "add" }
func (inst Minus) String() string { return "min" }
func (inst Div) String() string   { return "div" }
func (inst Mult) String() string  { return "mul" }

func (inst Equal) String() string   { return "eql" }
func (inst Greater) String() string { return "gtt" }
func (inst Lesser) String() string  { return "ltt" }

func (inst PushInt) String() string { return fmt.Sprintf("psh %d", int(inst)) }
func (inst Asg) String() string     { return fmt.Sprintf("stm %d", int(inst)) }
func (inst Deref) String() string   { return fmt.Sprintf("ldm %d", int(inst)) }
func (inst Output) String() string  { return fmt.Sprintf("out %d", int(inst)) }
