package main

import (
	"flag"
	"fmt"
	"log"
	"nand2tetris-go/asm"
	"nand2tetris-go/parser"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Assemble(asmFilePath string, hackFilePath string) {
	parsedLines := []parser.ParsedLine{}
	prs, err := parser.NewParser(asmFilePath)
	if err != nil {
		log.Fatalf("Error creating parser: %v", err)
	}
	defer prs.Close()
	assembler, err := asm.NewAssembler(hackFilePath)
	if err != nil {
		log.Fatalf("Error creating assembler: %v", err)
	}
	defer assembler.Close()

	// first pass
	lineCounter := 0
	symbolCounter := 16
	hasMoreLines := true
	for hasMoreLines {
		lineCounter++
		parsedLine, hasNext := prs.Advance()
		hasMoreLines = hasNext
		if parsedLine != nil {
			if l, ok := parsedLine.(parser.Label); ok {
				lineCounter--
				assembler.SymbolTable[l.Name] = lineCounter
			}
			if a, ok := parsedLine.(parser.AInstruction); ok {
				if _, exists := assembler.SymbolTable[a.Symbol]; !exists {
					if _, err := strconv.Atoi(a.Symbol); err != nil {
						assembler.SymbolTable[a.Symbol] = symbolCounter
						symbolCounter++
					}
				}
			}
			parsedLines = append(parsedLines, parsedLine)
		}
	}

	// second pass

}

func ParseArguments() string {
	flag.Parse()
	targetFilePath := flag.Arg(0)
	if targetFilePath == "" {
		fmt.Println("No target .asm file provided. Use the first argument to specify a target .asm file path.")
		os.Exit(1)
	}
	if _, err := os.Stat(targetFilePath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist:", targetFilePath)
		} else if os.IsPermission(err) {
			fmt.Println("Permission denied:", targetFilePath)
		} else {
			fmt.Println("Error accessing file:", err)
		}
		os.Exit(1)
	}
	return targetFilePath
}

func main() {
	targetFilePath := ParseArguments()
	targetFilePathWithoutExtension := strings.Trim(targetFilePath, filepath.Ext(targetFilePath))
	hackFilePath := fmt.Sprintf("%s.hack", targetFilePathWithoutExtension)
	Assemble(targetFilePath, hackFilePath)
}
