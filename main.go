package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"terhaak.de/imp/pkg/asm"
	"terhaak.de/imp/pkg/lexer"
	"terhaak.de/imp/pkg/vm"
)

func parseMemFlags(params []asm.Parameter) vm.Program {
	// Set up flag parser using the flag.Value interface
	for _, p := range params {
		if v, ok := p.(asm.StringParameter); ok {
			flag.Var(v, v.Name, "")
		}
		if v, ok := p.(asm.IntParameter); ok {
			flag.Var(v, v.Name, "")
		}
	}
	flag.Parse()

	// generate VM code to set the memory locations
	prog := make(vm.Program, 0, 2*len(params))
	for _, p := range params {
		if v, ok := p.(asm.StringParameter); ok {
			prog = append(prog, vm.PushStr(*v.Value), vm.StoreMemory(v.Address))
		}
		if v, ok := p.(asm.IntParameter); ok {
			prog = append(prog, vm.PushInt(*v.Value), vm.StoreMemory(v.Address))
		}
	}
	return prog
}

func execEmbedded() {
	if prog, meta, err := asm.LoadEmbeddedAssembly(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	} else if prog != nil {
		prog = append(parseMemFlags(meta.Params), prog...)
		if err := vm.RunProgram(prog); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(2)
		}
		os.Exit(0)
	}
}

func runLexer(fileName string) error {
	code, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	l := lexer.NewWithDefaultRules(string(code))
	for tok := l.Next(); !lexer.IsEof(tok); tok = l.Next() {
		fmt.Printf("%v\n", tok)
	}
	return nil
}

func main() {
	execEmbedded()

	asmCmd := flag.NewFlagSet("asm", flag.ExitOnError)
	asmFile := asmCmd.String("f", "", "Path to the asm file to run")
	asmOutFile := asmCmd.String("embed", "", "Path to new file to create with VM and embedded code")

	lexCmd := flag.NewFlagSet("lex", flag.ExitOnError)
	lexFile := lexCmd.String("f", "", "Path to the IMP code file to lex")

	switch os.Args[1] {
	case "asm":
		asmCmd.Parse(os.Args[2:])
	case "lex":
		lexCmd.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	var err error
	if asmCmd.Parsed() {
		if *asmFile != "" && *asmOutFile == "" {
			err = asm.RunAssemblyFile(*asmFile)
		} else if *asmFile != "" && *asmOutFile != "" {
			err = asm.EmbedAssemblyFile(*asmOutFile, *asmFile)
		}

	} else if lexCmd.Parsed() {
		if *lexFile != "" {
			err = runLexer(*lexFile)
		} else {
			err = fmt.Errorf("missing mandatory file parameter")
		}
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
