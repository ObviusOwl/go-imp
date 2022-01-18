package vm

import (
	"reflect"
	"testing"
)

func TestLoadAssemblyFile(t *testing.T) {
	expected := Program{
		PushInt(6),
		PushInt(5),
		Lesser{},
		JumpNonZero(1),
		PushInt(20),
		Jump(2),
		Label(1),
		PushInt(10),
		Label(2),
		Asg(1),
		Output(1),
		PushInt(4),
		PushInt(5),
		Add{},
	}

	actual, err := LoadAssemblyFile("testdata/test1.asm")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v, but got %v", expected, actual)
	}
}
