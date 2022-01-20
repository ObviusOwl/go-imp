package asm

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"terhaak.de/imp/pkg/vm"
)

var opCodeReg = regexp.MustCompile(`^\s*([a-zA-Z]+)(?:\s+([0-9]+))?\s*$`)

func coerceInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 0)
	return int(i)
}

func ParseMnemonic(s string) (vm.Executer, error) {
	m := opCodeReg.FindStringSubmatch(s)
	if m == nil {
		return nil, fmt.Errorf("invalid mnemonic: %s", s)
	}

	switch strings.ToLower(m[1]) {

	case "lab":
		return vm.Label(coerceInt(m[2])), nil
	case "jmp":
		return vm.Jump(coerceInt(m[2])), nil
	case "jnz":
		return vm.JumpNonZero(coerceInt(m[2])), nil
	case "jez":
		return vm.JumpZero(coerceInt(m[2])), nil
	case "stop":
		return vm.Stop{}, nil

	case "add":
		return vm.Add{}, nil
	case "min":
		return vm.Minus{}, nil
	case "div":
		return vm.Div{}, nil
	case "mul":
		return vm.Mult{}, nil

	case "eql":
		return vm.Equal{}, nil
	case "gtt":
		return vm.Greater{}, nil
	case "ltt":
		return vm.Lesser{}, nil

	case "psh":
		return vm.PushInt(coerceInt(m[2])), nil
	case "stm":
		return vm.StoreMemory(coerceInt(m[2])), nil
	case "ldm":
		return vm.LoadMemory(coerceInt(m[2])), nil
	case "out":
		return vm.Output(coerceInt(m[2])), nil

	default:
		return nil, fmt.Errorf("unknown opcode %s", m[0])
	}

}

func ParseAssemblyFile(file io.Reader) (vm.Program, error) {
	var program vm.Program
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

func LoadAssemblyFile(path string) (vm.Program, error) {
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
		return vm.RunProgram(program)
	}
	return err
}

func DumpAssemblyProgram(prog vm.Program) {
	for _, inst := range prog {
		fmt.Printf("%v\n", inst)
	}
}

const magicString = "embeddedcodecode"

func LoadEmbeddedAssembly() (vm.Program, error) {
	exeFile, err := os.Open("/proc/self/exe")
	if err != nil {
		return nil, err
	}
	defer exeFile.Close()

	// [binary:var][asm_text:var][asm_size:8Byte][magic:fixed]
	suffixSize := int64(len(magicString) + 8)

	if _, err := exeFile.Seek(-1*suffixSize, io.SeekEnd); err != nil {
		return nil, err
	}

	var codeSize int64
	if suffix, err := ioutil.ReadAll(exeFile); err != nil {
		return nil, err
	} else if !bytes.Equal(suffix[8:], []byte(magicString)) {
		// no embedded program
		return nil, nil
	} else {
		// network byte order
		codeSize = int64(binary.BigEndian.Uint64(suffix[:8]))
	}

	if _, err := exeFile.Seek(-1*(codeSize+suffixSize), io.SeekEnd); err != nil {
		return nil, err
	}

	if code, err := ioutil.ReadAll(exeFile); err != nil {
		return nil, err
	} else {
		codeReader := bytes.NewReader(code[:len(code)-int(suffixSize)])
		return ParseAssemblyFile(codeReader)
	}
}

func EmbedAssembly(targetFile io.Writer, sourceFile io.Reader) error {
	exeFile, err := os.Open("/proc/self/exe")
	if err != nil {
		return err
	}
	defer exeFile.Close()

	// copy ELF go binary file first
	if _, err = io.Copy(targetFile, exeFile); err != nil {
		return err
	}
	exeFile.Seek(0, io.SeekEnd)

	// append the source file
	codeSize, err := io.Copy(targetFile, sourceFile)
	if err != nil {
		return err
	}

	// write size of code in network byte order and fixed size
	if err := binary.Write(targetFile, binary.BigEndian, uint64(codeSize)); err != nil {
		return err
	}
	// write magic string to detect embedded code
	if _, err := targetFile.Write([]byte(magicString)); err != nil {
		return err
	}
	return nil
}

func EmbedAssemblyFile(target, source string) error {
	binTarget, err := os.Create(target)
	if err == nil {
		defer binTarget.Close()
		binSource, err := os.Open(source)
		if err == nil {
			defer binSource.Close()
			return EmbedAssembly(binTarget, binSource)
		}
	}
	return err
}
