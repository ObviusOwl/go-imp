package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"terhaak.de/imp/pkg/asm"
	"terhaak.de/imp/pkg/vm"
)

func parseMemFlags() vm.Program {
	var prog vm.Program
	if len(os.Args) == 1 {
		return prog
	}
	r := regexp.MustCompile("^([0-9]+):(.*)$")
	for _, f := range os.Args[1:] {
		if m := r.FindStringSubmatch(f); m != nil {
			addr, _ := strconv.ParseInt(m[1], 10, 0)
			value, _ := strconv.ParseInt(m[2], 10, 0)
			prog = append(prog, vm.PushInt(value), vm.StoreMemory(addr))
		}
	}
	return prog
}

func execEmbedded() {
	if prog, err := asm.LoadEmbeddedAssembly(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	} else if prog != nil {
		prog = append(parseMemFlags(), prog...)
		if err := vm.RunProgram(prog); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(2)
		}
		os.Exit(0)
	}

}

func main() {
	execEmbedded()

	asmCmd := flag.NewFlagSet("asm", flag.ExitOnError)
	asmFile := asmCmd.String("f", "", "Path to the asm file to run")
	asmOutFile := asmCmd.String("embed", "", "Path to new file to create with VM and embedded code")

	switch os.Args[1] {
	case "asm":
		asmCmd.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if asmCmd.Parsed() {
		var err error
		if *asmFile != "" && *asmOutFile == "" {
			err = asm.RunAssemblyFile(*asmFile)
		} else if *asmFile != "" && *asmOutFile != "" {
			err = asm.EmbedAssemblyFile(*asmOutFile, *asmFile)
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
