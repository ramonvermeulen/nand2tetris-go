package parser

import "fmt"

type ParsedLine interface {
	String() string
}

type Instruction interface {
	BinaryString(symbolTable *map[string]int) string
}

type AInstruction struct {
	Symbol string
}

func (inst AInstruction) String() string {
	return fmt.Sprintf("@%s\n", inst.Symbol)
}

func (inst AInstruction) BinaryString(symbolTable *map[string]int) string {
	return ""
}

type CInstruction struct {
	Comp string
	Dest string
	Jump string
}

func (inst CInstruction) String() string {
	if inst.Jump == "" {
		return fmt.Sprintf("%s=%s\n", inst.Dest, inst.Comp)
	}
	if inst.Dest == "" {
		return fmt.Sprintf("%s;%s\n", inst.Comp, inst.Jump)
	}
	return fmt.Sprintf("%s=%s;%s\n", inst.Dest, inst.Comp, inst.Jump)
}

func (inst CInstruction) BinaryString(symbolTable *map[string]int) string {
	return ""
}

type Label struct {
	Name string
}

func (inst Label) String() string {
	return fmt.Sprintf("(%s)\n", inst.Name)
}
