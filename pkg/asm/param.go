package asm

import (
	"fmt"
	"strconv"
)

// StringParameter implements also the flag.Value interface
type StringParameter struct {
	Name    string
	Address int
	Value   *string
}

// IntParameter implements also the flag.Value interface
type IntParameter struct {
	Name    string
	Address int
	Value   *int
}

func (p StringParameter) String() string {
	if p.Value == nil {
		return ""
	}
	return *p.Value
}

func (p StringParameter) Set(s string) error {
	*p.Value = s
	return nil
}

func (p IntParameter) String() string {
	if p.Value == nil {
		return ""
	}
	return fmt.Sprint(*p.Value)
}

func (p IntParameter) Set(s string) error {
	i, err := strconv.ParseInt(s, 10, 0)
	if err == nil {
		*p.Value = int(i)
	}
	return err
}
