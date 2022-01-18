package vm

import (
	"fmt"

	"terhaak.de/imp/pkg/stack"
)

//
// interfaces
//

type VM struct {
	program Program
	pc      int
	stack   stack.Stack
	mem     map[int]DataValue
}

type Executer interface {
	Exec(st stack.Stack, vm Machine) error
}

// A Memory is a basic key-value store mapping an address to a DataValue.
// The VM has virtual memory: a memory access always succeeds, and there is
// infinite memory. If an unset address is accessed, the value is undefined.
type Memory interface {
	Load(address int) DataValue
	Store(address int, value DataValue)
}

// A Runner is a type that can run a Program (slice of Executable) sequentially.
// This is the interface to a basic control unit within the CPU.
type Runner interface {
	Run(program Program) error
	Jump(label Label) error
}

// A Machine is the abstraction of the full VM, it has Memory and a Runner
type Machine interface {
	Memory
	Runner
}

// A label is a special instruction for locating jump targets.
// The types value represents a unique identifier.
type Label int

func (inst Label) Exec(st stack.Stack, vm Machine) error { return nil }

// A Program is a series (slice) of executable instructions
type Program []Executer

// A DataValue represents a value stored in memory or on the stack.
// Operations (such as arithmetic instructions) operate on DataValues.
// Any value can be stored, operations must cast the value to their native type.
// Casting errors cause the VM to halt as it is the responsibility of the compiler
// to produce meaningful and valid VM code.
type DataValue interface{}

//
// VM implementation
//

func New() *VM {
	var vm VM
	vm.stack = stack.New()
	vm.mem = make(map[int]DataValue)
	return &vm
}

func (vm *VM) Jump(label Label) error {
	for idx, inst := range vm.program {
		if value, ok := inst.(Label); ok && value == label {
			vm.pc = idx
			return nil
		}
	}
	return fmt.Errorf("segmentation fault: jump label not found %v", label)

}

// Run the given program. The memory and stack are not reset.
func (vm *VM) Run(program Program) error {
	vm.program = program
	vm.pc = 0

	for ; vm.pc < len(vm.program); vm.pc++ {
		err := vm.program[vm.pc].Exec(vm.stack, vm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (vm *VM) Load(address int) DataValue {
	value, ok := vm.mem[address]
	if !ok {
		return nil
	}
	return value
}

func (vm *VM) Store(address int, value DataValue) {
	vm.mem[address] = value
}

func RunProgram(program Program) error {
	return New().Run(program)
}
