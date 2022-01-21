package vm

import (
	"fmt"

	"terhaak.de/imp/pkg/stack"
)

//
// interfaces
//

type MapMemory map[int]DataValue

type DefaultRunner struct {
	program Program
	pc      int
	stack   stack.Stack
}

// A Machine is the abstraction of the full VM, it _has_ a Memory and a Runner
type Machine struct {
	ctrl DefaultRunner
	mem  MapMemory
}

// An Executer implements an instuction using the given environment and resources.
type Executer interface {
	Exec(vm Runner, st stack.Stack, mem Memory) error
}

// A Program is a series (slice) of executable instructions
type Program []Executer

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
	Run(program Program, mem Memory) error
	Stop() error
	Jump(label Label) error
}

// A label is a special instruction for locating jump targets.
// The types value represents a unique identifier.
type Label int

func (inst Label) Exec(vm Runner, st stack.Stack, mem Memory) error { return nil }

// A DataValue represents a value stored in memory or on the stack.
// Operations (such as arithmetic instructions) operate on DataValues.
// Any value can be stored, operations must cast the value to their native type.
// Casting errors cause the VM to halt as it is the responsibility of the compiler
// to produce meaningful and valid VM code.
type DataValue interface{}

//
// VM implementation
//

func New() *Machine {
	var vm Machine
	vm.ctrl.stack = stack.New()
	vm.mem = make(MapMemory)
	return &vm
}

func (ctrl *DefaultRunner) Jump(label Label) error {
	for idx, inst := range ctrl.program {
		if value, ok := inst.(Label); ok && value == label {
			ctrl.pc = idx
			return nil
		}
	}
	return fmt.Errorf("segmentation fault: jump label not found %v", label)
}

func (ctrl *DefaultRunner) Run(program Program, mem Memory) error {
	ctrl.stack = stack.New()
	ctrl.program = program
	ctrl.pc = 0

	for ; ctrl.pc < len(ctrl.program); ctrl.pc++ {
		err := ctrl.program[ctrl.pc].Exec(ctrl, ctrl.stack, mem)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctrl *DefaultRunner) Stop() error {
	ctrl.pc = len(ctrl.program)
	return nil
}

func (mem MapMemory) Load(address int) DataValue {
	value, ok := mem[address]
	if !ok {
		return nil
	}
	return value
}

func (mem MapMemory) Store(address int, value DataValue) {
	mem[address] = value
}

func RunProgram(program Program) error {
	vm := New()
	return vm.ctrl.Run(program, vm.mem)
}
