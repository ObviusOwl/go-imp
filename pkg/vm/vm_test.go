package vm

import (
	"fmt"
	"testing"

	"terhaak.de/imp/pkg/stack"
)

type vmMock struct {
	mem     interface{}
	stack   stack.Stack
	pc      Label
	pcSet   bool
	stopSet bool
}

func newMockVM() *vmMock {
	vm := &vmMock{}
	vm.stack = stack.New()
	vm.pcSet = false
	vm.stopSet = false
	return vm
}

func (vm *vmMock) Run(program Program) error {
	return nil
}
func (vm *vmMock) Stop() error {
	vm.stopSet = true
	return nil
}

func (vm *vmMock) Jump(label Label) error {
	vm.pc = label
	vm.pcSet = true
	return nil
}

func (vm *vmMock) Load(address int) DataValue {
	return vm.mem
}

func (vm *vmMock) Store(address int, value DataValue) {
	vm.mem = value
}

func (vm *vmMock) expectMemInt(expected int) error {
	if actual, ok := vm.mem.(int); !ok {
		return fmt.Errorf("Expected VM memory to be int, but got %v", vm.mem)
	} else if actual != expected {
		return fmt.Errorf("Expected VM memory to be %v, but got %v", expected, actual)
	}
	return nil
}

func (vm *vmMock) expectStackInt(expected int) error {
	item, err := vm.stack.Peek()
	if err != nil {
		return err
	}
	actual, ok := item.(int)
	if !ok {
		return fmt.Errorf("Expected top of the stack to be int, got %v", item)
	}
	if actual != expected {
		return fmt.Errorf("Expected top of stack to be %d, got %d", expected, actual)
	}
	return nil
}

func (vm *vmMock) expectStackStr(expected string) error {
	item, err := vm.stack.Peek()
	if err == nil {
		if actual, ok := item.(string); !ok {
			return fmt.Errorf("Expected top of the stack to be string, got %v", item)
		} else if actual != expected {
			return fmt.Errorf("Expected top of the stack to be '%s', got '%s'", expected, actual)
		}
	}
	return err
}

func (vm *vmMock) expectJump(expected bool, label Label) error {
	if !expected && vm.pcSet {
		return fmt.Errorf("Expected no jump, but got label %v", vm.pc)
	} else if expected && !vm.pcSet {
		return fmt.Errorf("Expected jump to label %v, but got nothing", label)
	} else if expected && vm.pc != label {
		return fmt.Errorf("Expected jump to label %v, but got %v", label, vm.pc)
	}
	return nil
}

func TestVMJump(t *testing.T) {
	vm := New()
	vm.program = Program{PushInt(4), Label(5), Label(6)}
	vm.pc = 0
	vm.Jump(6)
	if vm.pc != 2 {
		t.Fatalf("Expected program counter to be %d, but got %d", 2, vm.pc)
	}
}

func TestVMLoadStore(t *testing.T) {
	vm := New()
	vm.Store(5, 99)
	if actual := vm.Load(5); actual != 99 {
		t.Fatalf("Expected memory value to be %d, but got %d", 99, actual)
	}
}

func TestVMRun(t *testing.T) {
	prog := Program{
		PushInt(1),
		PushInt(1),
		Jump(5),
		Add{},
		Label(5),
		PushInt(99),
	}

	vm := New()
	if err := vm.Run(prog); err != nil {
		t.Fatal(err)
	}
	if values, _ := popInts(vm.stack, 1); values[0] != 99 {
		t.Fatalf("Expected stack top for %d, but got %d", 99, values[0])
	}
}

func TestVMStop(t *testing.T) {
	prog := Program{
		PushInt(99),
		Stop{},
		PushInt(1),
	}

	vm := New()
	if err := vm.Run(prog); err != nil {
		t.Fatal(err)
	}
	if values, _ := popInts(vm.stack, 1); values[0] != 99 {
		t.Fatalf("Expected stack top to be %d, but got %d", 99, values[0])
	}
}
