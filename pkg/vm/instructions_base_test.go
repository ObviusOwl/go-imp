package vm

import (
	"fmt"
	"testing"
)

type execTestCase struct {
	name string
	a    interface{}
	b    interface{}
	exp  interface{}
	err  bool
}

func newExecTestVM(tc execTestCase) *vmMock {
	vm := newMockVM()
	vm.stack.Push(tc.b)
	vm.stack.Push(tc.a)
	return vm
}

func callExec(tc execTestCase, vm *vmMock, subject Executer, valueTest func() error) error {
	err := subject.Exec(vm.stack, vm)
	if err != nil && tc.err {
		return nil
	} else if err != nil && !tc.err {
		return fmt.Errorf("Expected no error, but got %v", err)
	} else if err == nil && tc.err {
		return fmt.Errorf("Expected error, but got nothing")
	}
	return valueTest()
}

func runIntExec(tc execTestCase, subject Executer) error {
	vm := newExecTestVM(tc)
	return callExec(tc, vm, subject, func() error {
		return vm.expectStackInt(tc.exp.(int))
	})
}

func TestAdd(t *testing.T) {
	cases := []execTestCase{
		{"0+0", 0, 0, 0, false},
		{"0+2", 0, 2, 2, false},
		{"5+8", 5, 8, 13, false},
		{"-5+15", -5, 15, 10, false},
		{"4+-20", 4, -20, -16, false},
		{"5+'y'", 5, "y", 0, true},
		{"'z'+7", "z", 7, 0, true},
		{"true+7", true, 7, 0, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := runIntExec(tc, Add{}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestLesser(t *testing.T) {
	cases := []execTestCase{
		{name: "0<0", a: 0, b: 0, exp: 0, err: false},
		{name: "5<8", a: 5, b: 8, exp: 1, err: false},
		{name: "4<-20", a: 4, b: -20, exp: 0, err: false},
		{name: "5<'y'", a: 5, b: "y", exp: 0, err: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := runIntExec(tc, Lesser{}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGreater(t *testing.T) {
	cases := []execTestCase{
		{name: "0<0", a: 0, b: 0, exp: 0, err: false},
		{name: "5<8", a: 5, b: 8, exp: 0, err: false},
		{name: "4<-20", a: 4, b: -20, exp: 1, err: false},
		{name: "8<5", a: 4, b: -20, exp: 1, err: false},
		{name: "5<'y'", a: 5, b: "y", exp: 0, err: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := runIntExec(tc, Greater{}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	cases := []execTestCase{
		{name: "0=0", a: 0, b: 0, exp: 1, err: false},
		{name: "0=1", a: 0, b: 1, exp: 0, err: false},
		{name: "5=5", a: 5, b: 5, exp: 1, err: false},
		{name: "8=5", a: 8, b: 5, exp: 0, err: false},
		{name: "5='y'", a: 5, b: "y", exp: 0, err: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := runIntExec(tc, Equal{}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestJump(t *testing.T) {
	vm := newMockVM()
	Jump(5).Exec(vm.stack, vm)
	if err := vm.expectJump(true, 5); err != nil {
		t.Fatal(err)
	}
}

func TestJumpNonZero(t *testing.T) {
	cases := []execTestCase{
		{name: "0", a: 0, b: 0, exp: false, err: false},
		{name: "1", a: 1, b: 1, exp: true, err: false},
		{name: "7", a: 7, b: 7, exp: true, err: false},
		{name: "x", a: "x", b: "x", exp: false, err: true},
	}

	for _, tc := range cases {
		vm := newExecTestVM(tc)
		err := callExec(tc, vm, JumpNonZero(5), func() error {
			return vm.expectJump(tc.exp.(bool), 5)
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
func TestJumpZero(t *testing.T) {
	cases := []execTestCase{
		{name: "0", a: 0, b: 0, exp: true, err: false},
		{name: "1", a: 1, b: 1, exp: false, err: false},
		{name: "7", a: 7, b: 7, exp: false, err: false},
		{name: "x", a: "x", b: "x", exp: false, err: true},
	}

	for _, tc := range cases {
		vm := newExecTestVM(tc)
		err := callExec(tc, vm, JumpZero(5), func() error {
			return vm.expectJump(tc.exp.(bool), 5)
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAsg(t *testing.T) {
	cases := []execTestCase{
		{name: "5", a: 5, b: 7, exp: 5, err: false},
		{name: "0", a: 0, b: 3, exp: 0, err: false},
	}

	for _, tc := range cases {
		vm := newExecTestVM(tc)
		err := callExec(tc, vm, StoreMemory(5), func() error {
			return vm.expectMemInt(tc.exp.(int))
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestDeref(t *testing.T) {
	vm := newMockVM()
	vm.mem = 99
	err := callExec(execTestCase{err: false}, vm, LoadMemory(5), func() error {
		return vm.expectStackInt(99)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestToString(t *testing.T) {
	cases := []struct {
		expected string
		value    fmt.Stringer
	}{
		{"lab 5", Label(5)},
		{"jmp 5", Jump(5)},
		{"jnz 5", JumpNonZero(5)},
		{"jez 5", JumpZero(5)},
		{"add", Add{}},
		{"min", Minus{}},
		{"div", Div{}},
		{"mul", Mult{}},
		{"eql", Equal{}},
		{"gtt", Greater{}},
		{"ltt", Lesser{}},
		{"psh 5", PushInt(5)},
		{"stm 5", StoreMemory(5)},
		{"ldm 5", LoadMemory(5)},
		{"out 5", Output(5)},
	}

	for _, tc := range cases {
		t.Run(tc.expected, func(t *testing.T) {
			if actual := tc.value.String(); actual != tc.expected {
				t.Fatalf("Expected '%s', but got '%v'", tc.expected, actual)
			}
		})
	}
}
