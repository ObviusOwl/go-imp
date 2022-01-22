package asm

import (
	"fmt"

	"terhaak.de/imp/pkg/vm"
)

// parses lab, jmp, jnz, jez from basic instructions set
type CtrlInstrParser struct{}

func (p CtrlInstrParser) Parse(name string, line string, lineNum int) (vm.Executer, int, error) {
	if name == "stop" {
		return vm.Stop{}, 0, nil
	}

	arg, l := parseIntArg(line)
	if l == 0 {
		return nil, 0, fmt.Errorf("expected int argument on line %d", lineNum)
	}

	switch name {
	case "lab":
		return vm.Label(arg), l, nil
	case "jmp":
		return vm.Jump(arg), l, nil
	case "jnz":
		return vm.JumpNonZero(arg), l, nil
	case "jez":
		return vm.JumpZero(arg), l, nil
	}
	return nil, 0, nil
}

// parses add, min, div, mul from basic instructions set
type MathInstrParser struct{}

func (p MathInstrParser) Parse(name string, line string, lineNum int) (vm.Executer, int, error) {
	switch name {
	case "add":
		return vm.Add{}, 0, nil
	case "min":
		return vm.Minus{}, 0, nil
	case "div":
		return vm.Div{}, 0, nil
	case "mul":
		return vm.Mult{}, 0, nil
	}
	return nil, 0, nil
}

// parses eql, gtt, ltt from basic instructions set
type LogicInstrParser struct{}

func (p LogicInstrParser) Parse(name string, line string, lineNum int) (vm.Executer, int, error) {
	switch name {
	case "eql":
		return vm.Equal{}, 0, nil
	case "gtt":
		return vm.Greater{}, 0, nil
	case "ltt":
		return vm.Lesser{}, 0, nil
	}
	return nil, 0, nil
}

// parses psh, stm, ldm, out from basic instructions set
type DataInstrParser struct{}

func (p DataInstrParser) Parse(name string, line string, lineNum int) (vm.Executer, int, error) {
	arg, l := parseIntArg(line)
	if l == 0 {
		return nil, 0, fmt.Errorf("expected int argument on line %d", lineNum)
	}

	switch name {
	case "psh":
		return vm.PushInt(arg), l, nil
	case "stm":
		return vm.StoreMemory(arg), l, nil
	case "ldm":
		return vm.LoadMemory(arg), l, nil
	case "out":
		return vm.Output(arg), l, nil
	}
	return nil, 0, nil
}

// parses cat, len, str, fmt from extended instructions set
type StrInstrParser struct{}

func (p StrInstrParser) Parse(name string, line string, lineNum int) (vm.Executer, int, error) {
	switch name {
	case "cat":
		return vm.ConcatStr{}, 0, nil
	case "len":
		return vm.LengthStr{}, 0, nil
	}

	arg, l := parseStrArg(line)
	if l == 0 {
		return nil, 0, fmt.Errorf("expected string argument on line %d", lineNum)
	}

	switch name {
	case "str":
		return vm.PushStr(arg), l, nil
	case "fmt":
		return vm.FormatStr(arg), l, nil
	}
	return nil, 0, nil
}
