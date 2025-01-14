package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Assemble(asmFilePath string, hackFilePath string) {
	return
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
