package asm

import (
	"fmt"
	"os"

	"terhaak.de/imp/pkg/vm"
)

type Metadata struct {
	Params []Parameter
}

type Parameter interface{}

func LoadAssemblyFile(path string) (vm.Program, Metadata, error) {
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		return ParseAssemblyFile(file)
	}
	return nil, Metadata{}, err
}

func RunAssemblyFile(path string) error {
	program, _, err := LoadAssemblyFile(path)
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
