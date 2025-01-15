package parser

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

type Parser struct {
	file           *os.File
	scanner        *bufio.Scanner
	isCommentBlock bool
}

func NewParser(filePath string) (*Parser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &Parser{file: file, scanner: bufio.NewScanner(file)}, nil
}

func (p *Parser) Close() error {
	return p.file.Close()
}

func (p *Parser) isSingleWord(s string) bool {
	re := regexp.MustCompile(`^\S+$`)
	return re.MatchString(s)
}

func (p *Parser) Advance() (ParsedLine, bool) {
	// TODO(ramon) refactor into multiple functions per instruction type
	if !p.scanner.Scan() {
		return nil, false
	}
	line := p.scanner.Text()
	trimmedLine := strings.TrimSpace(line)

	// Multi-line comment block
	if p.isCommentBlock {
		if strings.HasSuffix(trimmedLine, "*/") {
			p.isCommentBlock = false
		}
		return p.Advance()
	}
	if strings.HasPrefix(trimmedLine, "/*") {
		// Multi-line syntax on single line
		if !strings.HasSuffix(trimmedLine, "*/") {
			p.isCommentBlock = true
		}
		return p.Advance()
	}

	// Whitespace or single-line comment
	if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
		return p.Advance()
	}

	firstChar := trimmedLine[:1]

	// Label
	if firstChar == "(" && strings.HasSuffix(trimmedLine, ")") && p.isSingleWord(trimmedLine) {
		return Label{Name: strings.Trim(trimmedLine, "()")}, true
	}

	// A-Instruction
	if firstChar == "@" && p.isSingleWord(trimmedLine) {
		return AInstruction{Symbol: strings.Replace(trimmedLine, "@", "", 1)}, true
	}

	// C-Instruction
	if strings.ContainsAny(firstChar, "01AD!-M") {
		re := regexp.MustCompile(`^(?:(\w+)=)?([^;]+)(?:;(\w+))?$`)
		matches := re.FindStringSubmatch(trimmedLine)
		dest, comp, jump := matches[1], matches[2], matches[3]
		return CInstruction{Dest: dest, Comp: comp, Jump: jump}, true
	}

	// Unparsable instruction meaning incorrect "hack" assembly
	log.Fatalf("Unparsable line: \"%s\"", line)
	os.Exit(1)
	return nil, false
}
