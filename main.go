package main

import (
	"flag"
	"fmt"
	"os"

	"terhaak.de/imp/pkg/vm"
)

func main() {
	asmCmd := flag.NewFlagSet("asm", flag.ExitOnError)
	asmFile := asmCmd.String("f", "", "Path to the asm file to run")

	switch os.Args[1] {
	case "asm":
		asmCmd.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if asmCmd.Parsed() {
		err := vm.RunAssemblyFile(*asmFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
