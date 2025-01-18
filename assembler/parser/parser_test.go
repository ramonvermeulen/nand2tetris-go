package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParser(t *testing.T) {
	input := "test input"
	reader := strings.NewReader(input)
	parser := NewParser(reader)

	assert.NotNil(t, parser)
	assert.Equal(t, reader, parser.reader)
}

func TestAdvance(t *testing.T) {
	input := `
// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/6/rect/Rect.asm

// Draws a rectangle at the top-left corner of the screen.
// The rectangle is 16 pixels wide and R0 pixels high.
// Usage: Before executing, put a value in R0.

   // If (R0 <= 0) goto END else n = R0
   @R0
   D=M
   @END
   D;JLE
   @n
   M=D
   // addr = base address of first screen row
   @SCREEN
   D=A
   @addr
   M=D
(LOOP)
   // RAM[addr] = -1
   @addr
   A=M
   M=-1
   // addr = base address of next screen row
   @addr
   D=M
   @32
   D=D+A
   @addr
   M=D
   // decrements n and loops
   @n
   MD=M-1
   @LOOP
   D;JGT
(END)
   @END
   0;JMP
`
	reader := strings.NewReader(input)
	parser := NewParser(reader)

	expectedLines := []ParsedLine{
		AInstruction{Symbol: "R0"},
		CInstruction{Dest: "D", Comp: "M", Jump: ""},
		AInstruction{Symbol: "END"},
		CInstruction{Dest: "", Comp: "D", Jump: "JLE"},
		AInstruction{Symbol: "n"},
		CInstruction{Dest: "M", Comp: "D", Jump: ""},
		AInstruction{Symbol: "SCREEN"},
		CInstruction{Dest: "D", Comp: "A", Jump: ""},
		AInstruction{Symbol: "addr"},
		CInstruction{Dest: "M", Comp: "D", Jump: ""},
		Label{Name: "LOOP"},
		AInstruction{Symbol: "addr"},
		CInstruction{Dest: "A", Comp: "M", Jump: ""},
		CInstruction{Dest: "M", Comp: "-1", Jump: ""},
		AInstruction{Symbol: "addr"},
		CInstruction{Dest: "D", Comp: "M", Jump: ""},
		AInstruction{Symbol: "32"},
		CInstruction{Dest: "D", Comp: "D+A", Jump: ""},
		AInstruction{Symbol: "addr"},
		CInstruction{Dest: "M", Comp: "D", Jump: ""},
		AInstruction{Symbol: "n"},
		CInstruction{Dest: "MD", Comp: "M-1", Jump: ""},
		AInstruction{Symbol: "LOOP"},
		CInstruction{Dest: "", Comp: "D", Jump: "JGT"},
		Label{Name: "END"},
		AInstruction{Symbol: "END"},
		CInstruction{Dest: "", Comp: "0", Jump: "JMP"},
	}

	for _, expected := range expectedLines {
		parsedLine, hasNext, err := parser.advance()
		assert.NoError(t, err)
		assert.True(t, hasNext)
		assert.Equal(t, expected, parsedLine)
	}

	_, hasNext, err := parser.advance()
	assert.NoError(t, err)
	assert.False(t, hasNext)
}

func TestParseAndAddSymbols(t *testing.T) {
	input := "@2\nD=A\n(LOOP)\n@LOOP\n0;JMP\n"
	reader := strings.NewReader(input)
	parser := NewParser(reader)

	symbolTable := map[string]int{
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
		"SCREEN": 16384,
		"KBD":    24576,
	}

	expectedLines := []ParsedLine{
		AInstruction{Symbol: "2"},
		CInstruction{Dest: "D", Comp: "A", Jump: ""},
		Label{Name: "LOOP"},
		AInstruction{Symbol: "LOOP"},
		CInstruction{Dest: "", Comp: "0", Jump: "JMP"},
	}

	parsedLines, err := parser.ParseAndAddSymbols(symbolTable)
	assert.NoError(t, err)
	assert.Equal(t, expectedLines, parsedLines)
	assert.Equal(t, 2, symbolTable["LOOP"])
}
