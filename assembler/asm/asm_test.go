package asm

import (
	"bytes"
	"fmt"
	"nand2tetris-go/parser"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssembleAIInstruction(t *testing.T) {
	assembler := &Assembler{
		SymbolTable: map[string]int{
			"R0":     0,
			"R4":     4,
			"SCREEN": 16384,
			"KBD":    24576,
		},
	}

	testCases := []struct {
		name     string
		input    parser.AInstruction
		expected string
		err      error
	}{
		{"ValidNumber", parser.AInstruction{Symbol: "72"}, "0000000001001000", nil},
		{"ValidSymbolR0", parser.AInstruction{Symbol: "R0"}, "0000000000000000", nil},
		{"ValidSymbolR4", parser.AInstruction{Symbol: "R4"}, "0000000000000100", nil},
		{"ValidSymbolSCREEN", parser.AInstruction{Symbol: "SCREEN"}, "0100000000000000", nil},
		{"ValidSymbolKBD", parser.AInstruction{Symbol: "KBD"}, "0110000000000000", nil},
		{"InvalidSymbol", parser.AInstruction{Symbol: "NOT_IN_TABLE"}, "", &strconv.NumError{Func: "Atoi", Num: "NOT_IN_TABLE", Err: strconv.ErrSyntax}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := assembler.assembleAInstruction(testCase.input)
			if testCase.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, testCase.err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestAssembleCInstruction(t *testing.T) {
	assembler := &Assembler{}

	testCases := []struct {
		name     string
		input    parser.CInstruction
		expected string
		err      error
	}{
		{"Valid", parser.CInstruction{Dest: "M", Comp: "D+M", Jump: "JGT"}, "1111000010001001", nil},
		{"ValidNoJump", parser.CInstruction{Dest: "D", Comp: "M", Jump: ""}, "1111110000010000", nil},
		{"ValidNoDest", parser.CInstruction{Dest: "", Comp: "M", Jump: "JGT"}, "1111110000000001", nil},
		{"InvalidDest", parser.CInstruction{Dest: "INVALID", Comp: "D+M", Jump: "JGT"}, "", fmt.Errorf("invalid dest mnemonic: INVALID")},
		{"InvalidComp", parser.CInstruction{Dest: "M", Comp: "INVALID", Jump: "JGT"}, "", fmt.Errorf("invalid comp mnemonic: INVALID")},
		{"InvalidJump", parser.CInstruction{Dest: "M", Comp: "D+M", Jump: "INVALID"}, "", fmt.Errorf("invalid jump mnemonic: INVALID")},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := assembler.assembleCInstruction(testCase.input)
			if testCase.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, testCase.err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestAssembleLine(t *testing.T) {
	testCases := []struct {
		name       string
		parsedLine parser.ParsedLine
		err        error
		expected   string
	}{
		{
			name:       "Label",
			parsedLine: parser.Label{},
			expected:   "",
			err:        nil,
		},
		{
			name:       "Valid AInstruction",
			parsedLine: parser.AInstruction{Symbol: "R0"},
			expected:   "0000000000000000\n",
			err:        nil,
		},
		{
			name:       "Invalid AInstruction",
			parsedLine: parser.AInstruction{Symbol: "NOT_IN_TABLE"},
			expected:   "",
			err:        fmt.Errorf("failed to assemble A-instruction: %v", &strconv.NumError{Func: "Atoi", Num: "NOT_IN_TABLE", Err: strconv.ErrSyntax}),
		},
		{
			name:       "Valid CInstruction",
			parsedLine: parser.CInstruction{Dest: "M", Comp: "D+M", Jump: "JGT"},
			expected:   "1111000010001001\n",
			err:        nil,
		},
		{
			name:       "Invalid CInstruction",
			parsedLine: parser.CInstruction{Dest: "INVALID", Comp: "D+M", Jump: "JGT"},
			expected:   "",
			err:        fmt.Errorf("failed to assemble C-instruction: invalid dest mnemonic: INVALID"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var buffer bytes.Buffer
			assembler := NewAssembler(&buffer)

			err := assembler.AssembleLine(testCase.parsedLine)
			if testCase.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, testCase.err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expected, buffer.String())
		})
	}
}
