package vm

import "testing"

func runStrExec(tc execTestCase, subject Executer) error {
	vm := newExecTestVM(tc)
	return callExec(tc, vm, subject, func() error {
		return vm.expectStackStr(tc.exp.(string))
	})
}

func TestPushStr(t *testing.T) {
	cases := []execTestCase{
		{"teststr", "teststr", 0, "teststr", false},
		{"empty", "", 2, "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := runStrExec(tc, PushStr(tc.a.(string))); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestConcatStr(t *testing.T) {
	cases := []execTestCase{
		{"a+b", "a", "b", "ab", false},
		{"a+empty", "a", "", "a", false},
		{"a+int", "a", 5, "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := runStrExec(tc, ConcatStr{}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestFormatStr(t *testing.T) {
	cases := []execTestCase{
		{"%s-%d", "a", 5, "a-5", false},
		{"-%s-", "a", "b", "-a-", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := runStrExec(tc, FormatStr(tc.name)); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestLengthStr(t *testing.T) {
	cases := []execTestCase{
		{"one", "a", "a", 1, false},
		{"empty", "", "", 0, false},
		{"int", 5, 5, "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vm := newExecTestVM(tc)
			err := callExec(tc, vm, LengthStr{}, func() error {
				return vm.expectStackInt(tc.exp.(int))
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
