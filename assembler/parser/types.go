package parser

import "fmt"

type ParsedLine interface {
	String() string
}

type AInstruction struct {
	Address int
}

func (inst AInstruction) String() string {
	return fmt.Sprintf("@%d\n", inst.Address)
}

type CInstruction struct {
	Comp string
	Dest string
	Jump string
}

func (inst CInstruction) String() string {
	// TODO(ramon) implement cases with no jump or dest
	return fmt.Sprintf("%s=%s;%s\n", inst.Dest, inst.Comp, inst.Jump)
}

type Label struct {
	Name string
}

func (inst Label) String() string {
	return fmt.Sprintf("(%s)\n", inst.Name)
}
