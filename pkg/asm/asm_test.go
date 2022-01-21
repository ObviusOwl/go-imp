package asm

import (
	"reflect"
	"testing"

	"terhaak.de/imp/pkg/vm"
)

func TestLoadAssemblyFile(t *testing.T) {
	expected := vm.Program{
		vm.PushInt(6),
		vm.PushInt(5),
		vm.Lesser{},
		vm.JumpNonZero(1),
		vm.PushInt(20),
		vm.Jump(2),
		vm.Label(1),
		vm.PushInt(10),
		vm.Label(2),
		vm.StoreMemory(1),
		vm.Output(1),
		vm.PushInt(4),
		vm.PushInt(5),
		vm.Add{},
	}

	actual, _, err := LoadAssemblyFile("testdata/test1.asm")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v, but got %v", expected, actual)
	}
}
