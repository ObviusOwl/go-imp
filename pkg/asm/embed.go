package asm

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"

	"terhaak.de/imp/pkg/vm"
)

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
		defer binTarget.Chmod(0755)
		binSource, err := os.Open(source)
		if err == nil {
			defer binSource.Close()
			return EmbedAssembly(binTarget, binSource)
		}
	}
	return err
}
