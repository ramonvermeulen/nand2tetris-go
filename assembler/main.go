package main

import (
	"flag"
	"fmt"
	"log"
	"nand2tetris-go/asm"
	"nand2tetris-go/parser"
	"os"
	"path/filepath"
	"strings"
)

func Assemble(asmFilePath string, hackFilePath string) {
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

	parsedLines, err := prs.ParseAndAddSymbols(assembler.SymbolTable)
	if err != nil {
		log.Fatalf("Error processing symbols: %v", err)
	}

	for _, parsedLine := range parsedLines {
		if err := assembler.AssembleLine(parsedLine); err != nil {
			log.Fatalf("Error assembling line: %v", err)
		}
	}
}

func ParseArguments() string {
	flag.Parse()
	targetFilePath := flag.Arg(0)
	if targetFilePath == "" {
		log.Fatalf("No target .asm file provided. Use the first argument to specify a target .asm file path.")
	}
	if _, err := os.Stat(targetFilePath); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf(fmt.Sprintf("File does not exist: %s", targetFilePath), err)
		} else if os.IsPermission(err) {
			log.Fatalf(fmt.Sprintf("Permission denied accessing file: %s", targetFilePath), err)
		} else {
			log.Fatalf("Error accessing file: %s", err)
		}
	}
	return targetFilePath
}

func main() {
	targetFilePath := ParseArguments()
	targetFilePathWithoutExtension := strings.Trim(targetFilePath, filepath.Ext(targetFilePath))
	hackFilePath := fmt.Sprintf("%s.hack", targetFilePathWithoutExtension)
	Assemble(targetFilePath, hackFilePath)
}
