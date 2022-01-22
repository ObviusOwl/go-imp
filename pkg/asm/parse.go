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

// MnemonicParser provides the method Parse to parse a Mnemonic
//
// First parameter the mnemonic name in lower case, second parameter
// is the text following directly the mnemonic name on the line including comments
// and whitespace. Third parameter is the current line number.
//
// The return values are first the parsed vm.Executer (nil if the parser did not match).
// Second the number of caracters consumed (may be 0 for 0-argument instructions)
// and last an error (nil for no error, not matching the line is not an error).
//
// All parsers are called until one returns a non-nil Executer.
// The error (if any) from the last parser that was run is used to abort the parsing.
// If no parser matches (produces a non-nil Executer) the parsing is aborted with an error.
//
// Parsers may be stateful. There may be multiple parsers for the same instruction
// as long as it is deterministic what parser handles which variant. For example
// an instruction name may be overloaded with an int or a string argument, which
// can be handled by different parsers.
type MnemonicParser interface {
	Parse(name string, line string, lineNum int) (vm.Executer, int, error)
}

var opNameReg = regexp.MustCompile(`^(?:\s*([a-zA-Z]+|;))|(?:\s*$)`)
var intArgReg = regexp.MustCompile(`^\s*(-?[0-9]+)`)
var strArgReg = regexp.MustCompile(`^\s*"([^"\\]*(?:\\.[^"\\]*)*)"`)
var paramReg = regexp.MustCompile(`^\s*@param\s+([a-zA-Z0-9]+)\s+(-?[0-9]+)\s+(str|int)\s+`)

func parseOpName(s string) (string, int, bool) {
	m := opNameReg.FindStringSubmatch(s)
	if m == nil {
		// no match = error, consume none
		return "", 0, false
	} else if m[1] == ";" {
		// comment only line, consume all
		return "", len(s), true
	} else if m[1] == "" {
		// whitespace only line, consume fullmatch
		return "", len(m[0]), true
	} else {
		// found name, consume name only
		return strings.ToLower(m[1]), len(m[1]), true
	}
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

// ParseAssemblyFile uses the default mnemonic parsers.
// Use ParseAssembly() for control over the used parsers
func ParseAssemblyFile(file io.Reader) (vm.Program, Metadata, error) {
	parsers := []MnemonicParser{
		CtrlInstrParser{},
		MathInstrParser{},
		LogicInstrParser{},
		DataInstrParser{},
		StrInstrParser{},
	}
	return ParseAssembly(file, parsers)
}

func ParseAssembly(file io.Reader, parsers []MnemonicParser) (vm.Program, Metadata, error) {
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

		opName, l, ok := parseOpName(line)
		if !ok {
			return nil, meta, fmt.Errorf("expected opname on line %d", lineNum)
		} else if opName == "" {
			continue
		}
		line = line[l:]

		var err error = nil
		var op vm.Executer = nil
		for _, parser := range parsers {
			if op, _, err = parser.Parse(opName, line, lineNum); op != nil {
				break
			}
		}
		if op == nil {
			return nil, meta, fmt.Errorf("unknown opcode %s", opName)
		} else if err != nil {
			return nil, meta, err
		} else {
			program = append(program, op)
		}

		isHeader = false
	}

	return program, meta, nil
}
