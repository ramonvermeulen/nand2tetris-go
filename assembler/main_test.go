package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssemble(t *testing.T) {
	asmContent := `
// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/6/max/Max.asm

// Computes R2 = max(R0, R1)  (R0,R1,R2 refer to RAM[0],RAM[1],RAM[2])
// Usage: Before executing, put two values in R0 and R1.

  // D = R0 - R1
  @R0
  D=M
  @R1
  D=D-M
  // If (D > 0) goto ITSR0
  @ITSR0
  D;JGT
  // Its R1
  @R1
  D=M
  @OUTPUT_D
  0;JMP
(ITSR0)
  @R0
  D=M
(OUTPUT_D)
  @R2
  M=D
(END)
  @END
  0;JMP
	`
	asmFile, err := os.CreateTemp("", "test.asm")
	assert.NoError(t, err)
	defer os.Remove(asmFile.Name())

	_, err = asmFile.WriteString(asmContent)
	assert.NoError(t, err)
	asmFile.Close()

	hackFile, err := os.CreateTemp("", "test.hack")
	assert.NoError(t, err)
	defer os.Remove(hackFile.Name())
	hackFile.Close()

	Assemble(asmFile.Name(), hackFile.Name())

	hackContent, err := os.ReadFile(hackFile.Name())
	assert.NoError(t, err)

	expectedHackContent := `0000000000000000
1111110000010000
0000000000000001
1111010011010000
0000000000001010
1110001100000001
0000000000000001
1111110000010000
0000000000001100
1110101010000111
0000000000000000
1111110000010000
0000000000000010
1110001100001000
0000000000001110
1110101010000111
`
	assert.Equal(t, expectedHackContent, string(hackContent))
}

func TestParseArguments(t *testing.T) {
	asmFile, err := os.CreateTemp("", "test.asm")
	assert.NoError(t, err)
	defer os.Remove(asmFile.Name())

	os.Args = []string{"cmd", asmFile.Name()}
	targetFilePath := ParseArguments()

	assert.Equal(t, asmFile.Name(), targetFilePath)
}
