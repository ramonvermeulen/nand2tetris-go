package asm

import (
	"fmt"
	"os"
)

type Assembler struct {
	file        *os.File
	SymbolTable map[string]int
}

func createInitialSymbolTable() map[string]int {
	symbolTable := map[string]int{
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
		"SCREEN": 16384,
		"KBD":    24576,
	}
	for i := 0; i < 16; i++ {
		symbolTable[fmt.Sprintf("R%d", i)] = i
	}
	return symbolTable
}

func NewAssembler(hackFilePath string) (*Assembler, error) {
	file, err := os.Create(hackFilePath)
	if err != nil {
		return nil, err
	}
	return &Assembler{
		SymbolTable: createInitialSymbolTable(),
		file:        file,
	}, nil
}

func (a *Assembler) Close() error {
	return a.file.Close()
}
