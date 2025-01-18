package asm

import (
	"fmt"
	"io"
	"nand2tetris-go/parser"
	"strconv"
)

type Assembler struct {
	writer        io.Writer
	SymbolTable   map[string]int
	symbolCounter int
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

func NewAssembler(writer io.Writer) *Assembler {
	return &Assembler{
		SymbolTable: createInitialSymbolTable(),
		writer:      writer,
	}
}

func (a *Assembler) assembleAInstruction(aInst parser.AInstruction) (string, error) {
	var address int
	var err error

	if match, ok := a.SymbolTable[aInst.Symbol]; ok {
		address = match
	} else {
		address, err = strconv.Atoi(aInst.Symbol)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("0%015b", address), nil
}

func (a *Assembler) assembleCInstruction(cInst parser.CInstruction) (string, error) {
	output := "111"

	comp, ok := compMnemonics[cInst.Comp]
	if !ok {
		return "", fmt.Errorf("invalid comp mnemonic: %s", cInst.Comp)
	}
	dest, ok := destMnemonics[cInst.Dest]
	if !ok {
		return "", fmt.Errorf("invalid dest mnemonic: %s", cInst.Dest)
	}
	jump, ok := jumpMnemonics[cInst.Jump]
	if !ok {
		return "", fmt.Errorf("invalid jump mnemonic: %s", cInst.Jump)
	}

	output += comp
	output += dest
	output += jump
	return output, nil
}

func (a *Assembler) AssembleLine(parsedLine parser.ParsedLine) error {
	// skip labels, since they are not instructions
	if _, ok := parsedLine.(parser.Label); ok {
		return nil
	}
	var output string
	var err error

	if aInst, ok := parsedLine.(parser.AInstruction); ok {
		output, err = a.assembleAInstruction(aInst)
		if err != nil {
			return fmt.Errorf("failed to assemble A-instruction: %v", err)
		}
	}
	if cInst, ok := parsedLine.(parser.CInstruction); ok {
		output, err = a.assembleCInstruction(cInst)
		if err != nil {
			return fmt.Errorf("failed to assemble C-instruction: %v", err)
		}
	}
	if output != "" {
		_, err = a.writer.Write([]byte(fmt.Sprintf("%s\n", output)))
		if err != nil {
			return fmt.Errorf("failed to write output to writer: %v", err)
		}
	}

	return nil
}
