package asm

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"terhaak.de/imp/pkg/vm"
)

var opNameReg = regexp.MustCompile(`^\s*([a-zA-Z]+)`)
var intArgReg = regexp.MustCompile(`^\s*(-?[0-9]+)`)
var strArgReg = regexp.MustCompile(`^\s*"([^"\\]*(?:\\.[^"\\]*)*)"`)
var paramReg = regexp.MustCompile(`^\s*@param\s+([a-zA-Z0-9]+)\s+(-?[0-9]+)\s+(str|int)\s+`)

func parseOpName(s string) (string, int) {
	m := opNameReg.FindStringSubmatch(s)
	if m == nil {
		return "", 0
	}
	return strings.ToLower(m[1]), len(m[0])
}

func parseIntArg(s string) (int, int) {
	m := intArgReg.FindStringSubmatch(s)
	if m == nil {
		return 0, 0
	}
	i, _ := strconv.ParseInt(m[1], 10, 0)
	return int(i), len(m[0])
}

func parseStrArg(s string) (string, int) {
	m := strArgReg.FindStringSubmatch(s)
	if m == nil {
		return "", 0
	}
	v, _ := strconv.Unquote(`"` + m[1] + `"`)
	return v, len(m[0])
}

func parseMnemonic(name string, line string, lineNum int) (vm.Executer, int, error) {
	var arg interface{}
	var l int = 0

	switch name {
	case "lab", "jmp", "jnz", "jez", "psh", "stm", "ldm", "out":
		arg, l = parseIntArg(line)
		if l == 0 {
			return nil, 0, fmt.Errorf("expected int argument on line %d", lineNum)
		}
	case "str", "fmt":
		arg, l = parseStrArg(line)
		if l == 0 {
			return nil, 0, fmt.Errorf("expected string argument on line %d", lineNum)
		}
	}

	switch name {
	case "lab":
		return vm.Label(arg.(int)), l, nil
	case "jmp":
		return vm.Jump(arg.(int)), l, nil
	case "jnz":
		return vm.JumpNonZero(arg.(int)), l, nil
	case "jez":
		return vm.JumpZero(arg.(int)), l, nil
	case "stop":
		return vm.Stop{}, l, nil

	case "add":
		return vm.Add{}, l, nil
	case "min":
		return vm.Minus{}, l, nil
	case "div":
		return vm.Div{}, l, nil
	case "mul":
		return vm.Mult{}, l, nil

	case "eql":
		return vm.Equal{}, l, nil
	case "gtt":
		return vm.Greater{}, l, nil
	case "ltt":
		return vm.Lesser{}, l, nil

	case "psh":
		return vm.PushInt(arg.(int)), l, nil
	case "stm":
		return vm.StoreMemory(arg.(int)), l, nil
	case "ldm":
		return vm.LoadMemory(arg.(int)), l, nil
	case "out":
		return vm.Output(arg.(int)), l, nil

	case "cat":
		return vm.ConcatStr{}, l, nil
	case "len":
		return vm.LengthStr{}, l, nil
	case "str":
		return vm.PushStr(arg.(string)), l, nil
	case "fmt":
		return vm.FormatStr(arg.(string)), l, nil

	default:
		return nil, l, fmt.Errorf("unknown opcode %s", name)
	}

}

func parseParam(line string, lineNum int) (*Parameter, error) {
	m := paramReg.FindStringSubmatch(line)
	if m == nil {
		return nil, nil
	}

	var p Parameter
	line = line[len(m[0]):]
	addr, _ := strconv.ParseInt(m[2], 10, 0)

	if m[3] == "str" {
		if arg, l := parseStrArg(line); l == 0 {
			return nil, fmt.Errorf("expected string argument on line %d", lineNum)
		} else {
			p = StringParameter{Name: m[1], Address: int(addr), Value: &arg}
		}
	} else if m[3] == "int" {
		if arg, l := parseIntArg(line); l == 0 {
			return nil, fmt.Errorf("expected int argument on line %d", lineNum)
		} else {
			p = IntParameter{Name: m[1], Address: int(addr), Value: &arg}
		}
	}

	return &p, nil
}

func ParseAssemblyFile(file io.Reader) (vm.Program, Metadata, error) {
	var program vm.Program
	var meta Metadata

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	isHeader := true
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if len(line) == 0 || (line[0] == ';' && !isHeader) {
			continue
		} else if line[0] == ';' && isHeader {
			param, err := parseParam(line[1:], lineNum)
			if err != nil {
				return nil, meta, err
			} else if param != nil {
				meta.Params = append(meta.Params, *param)
			}
			continue
		}

		line = strings.Split(line, ";")[0]

		opName, l := parseOpName(line)
		if l == 0 {
			return nil, meta, fmt.Errorf("expected opname on line %d", lineNum)
		}
		line = line[l:]

		op, _, err := parseMnemonic(opName, line, lineNum)
		if err != nil {
			return nil, meta, err
		}
		if op != nil {
			program = append(program, op)
		}

		isHeader = false
	}

	return program, meta, nil
}
