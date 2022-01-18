package vm

import (
	"reflect"
	"testing"

	"terhaak.de/imp/pkg/stack"
)

func TestPopInts(t *testing.T) {
	if _, err := popInts(stack.NewWithItems(1, 2, "aaa", 5), 3); err == nil {
		t.Fatalf("Expected error, but got nothing")
	}

	if _, err := popInts(stack.NewWithItems(1, 2, 3), 10); err == nil {
		t.Fatalf("Expected error, but got nothing")
	}

	expected := []int{8, 6}
	actual, err := popInts(stack.NewWithItems(8, 6), 2)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v, but got %v", expected, actual)
	}
}
