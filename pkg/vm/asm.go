package vm

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var opCodeReg = regexp.MustCompile(`^\s*([a-zA-Z]+)(?:\s+([0-9]+))?\s*$`)

func coerceInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 0)
	return int(i)
}

func ParseMnemonic(s string) (Executer, error) {
	m := opCodeReg.FindStringSubmatch(s)
	if m == nil {
		return nil, fmt.Errorf("invalid mnemonic: %s", s)
	}

	switch strings.ToLower(m[1]) {

	case "lab":
		return Label(coerceInt(m[2])), nil
	case "jmp":
		return Jump(coerceInt(m[2])), nil
	case "jnz":
		return JumpNonZero(coerceInt(m[2])), nil
	case "jez":
		return JumpZero(coerceInt(m[2])), nil

	case "add":
		return Add{}, nil
	case "min":
		return Minus{}, nil
	case "div":
		return Div{}, nil
	case "mul":
		return Mult{}, nil

	case "eql":
		return Equal{}, nil
	case "gtt":
		return Greater{}, nil
	case "ltt":
		return Lesser{}, nil

	case "psh":
		return PushInt(coerceInt(m[2])), nil
	case "stm":
		return Asg(coerceInt(m[2])), nil
	case "ldm":
		return Deref(coerceInt(m[2])), nil
	case "out":
		return Output(coerceInt(m[2])), nil

	default:
		return nil, fmt.Errorf("unknown opcode %s", m[0])
	}

}

func ParseAssemblyFile(file *os.File) (Program, error) {
	var program Program
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == ';' {
			continue
		}
		line = strings.Split(line, ";")[0]

		inst, err := ParseMnemonic(line)
		if err != nil {
			return nil, err
		}

		program = append(program, inst)
	}

	return program, nil
}

func LoadAssemblyFile(path string) (Program, error) {
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		return ParseAssemblyFile(file)
	}
	return nil, err
}

func RunAssemblyFile(path string) error {
	program, err := LoadAssemblyFile(path)
	if err == nil {
		return RunProgram(program)
	}
	return err
}

func DumpAssemblyProgram(prog Program) {
	for _, inst := range prog {
		fmt.Printf("%v\n", inst)
	}
}
