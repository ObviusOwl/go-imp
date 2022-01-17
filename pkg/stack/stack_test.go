package stack

import (
	"fmt"
	"testing"
)

func expectTosInt(s *stack, expectedTos int, expectedValue int) (string, bool) {
	if s.tos != expectedTos {
		return fmt.Sprintf("Expected stack.tos to be %d, got %d", expectedTos, s.tos), false
	}

	if expectedTos == -1 {
		if len(s.items) != 0 {
			return fmt.Sprintf("Expected empty stack, got %d items", len(s.items)), false
		}
	} else {
		if len(s.items) == 0 {
			return "Expected non-empty stack", false
		}

		actual, ok := s.items[s.tos].(int)
		if !ok {
			return "Failed to convert top of stack to int", false
		}
		if actual != expectedValue {
			return fmt.Sprintf("Expected top of stack to be %d, got %d", expectedValue, actual), false
		}
	}
	return "", true
}

func excpectInt(item interface{}, expected int) (string, bool) {
	actual, ok := item.(int)
	if !ok {
		return fmt.Sprintf("Failed to convert value %v to int", item), false
	}
	if actual != expected {
		return fmt.Sprintf("Expected value to be %d, got %d", 8, actual), false
	}
	return "", true
}

func pushValues(s *stack, items []int) {
	for _, item := range items {
		s.Push(item)
	}
}

func TestNew(t *testing.T) {
	s := New()
	if err, ok := expectTosInt(s, -1, -1); !ok {
		t.Fatal(err)
	}
}

func TestPush(t *testing.T) {
	s := New()

	t.Run("first", func(t *testing.T) {
		s.Push(5)
		if err, ok := expectTosInt(s, 0, 5); !ok {
			t.Fatal(err)
		}
	})
	t.Run("second", func(t *testing.T) {
		s.Push(6)
		if err, ok := expectTosInt(s, 1, 6); !ok {
			t.Fatal(err)
		}
	})
}

func TestPop(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := New()
		_, err := s.Pop()
		if err == nil {
			t.Fatalf("Expected error on pop empty stack")
		}
	})

	cases := []struct {
		name     string
		items    []int
		tosIndex int
		tosValue int
		expected int
	}{
		{"one", []int{5}, -1, -1, 5},
		{"two", []int{5, 6}, 0, 5, 6},
		{"three", []int{5, 6, 9}, 1, 6, 9},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := New()
			pushValues(s, tc.items)
			item, err := s.Pop()
			if err != nil {
				t.Fatalf("Expected no error, but got %v", err)
			}

			if err, ok := excpectInt(item, tc.expected); !ok {
				t.Fatal(err)
			}

			if err, ok := expectTosInt(s, tc.tosIndex, tc.tosValue); !ok {
				t.Fatal(err)
			}
		})
	}

}

func TestEmpty(t *testing.T) {
	cases := []struct {
		name     string
		items    []int
		tosIndex int
		tosValue int
		expected bool
	}{
		{"zero", nil, -1, -1, true},
		{"one", []int{5}, 0, 5, false},
		{"two", []int{5, 6}, 1, 6, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := New()
			pushValues(s, tc.items)

			actual := s.Empty()
			if actual != tc.expected {
				t.Fatalf("Expected %v, but got %v", tc.expected, actual)
			}

			if err, ok := expectTosInt(s, tc.tosIndex, tc.tosValue); !ok {
				t.Fatal(err)
			}
		})
	}
}

func TestPeek(t *testing.T) {
	s := New()

	t.Run("empty", func(t *testing.T) {
		if _, err := s.Peek(); err == nil {
			t.Fatal("Expected error on peek empty stack")
		}
	})

	s.Push(8)

	item, err := s.Peek()

	t.Run("no-side-effect", func(t *testing.T) {
		if err, ok := expectTosInt(s, 0, 8); !ok {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if err, ok := excpectInt(item, 8); !ok {
		t.Fatal(err)
	}
}

func TestPushPop(t *testing.T) {
	s := New()
	s.Push(99)

	for i := 1; i < 2; i++ {
		t.Run(fmt.Sprintf("pass%d", i), func(t *testing.T) {
			s.Push(i)
			item, err := s.Pop()

			if err != nil {
				t.Fatalf("Expected no error, got %s", err)
			}
			if err, ok := excpectInt(item, i); !ok {
				t.Fatal(err)
			}
		})
	}
}
