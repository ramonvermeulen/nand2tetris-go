package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ParsingError struct {
	Line string
	Err  error
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("error parsing line \"%s\": %v", e.Line, e.Err)
}

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

func (p *Parser) advance() (ParsedLine, bool, error) {
	if !p.scanner.Scan() {
		if err := p.scanner.Err(); err != nil {
			return nil, false, fmt.Errorf("scanner error: %w", err)
		}
		return nil, false, nil // EOF
	}

	line := p.scanner.Text()
	trimmedLine := strings.TrimSpace(line)
	if p.handleComments(trimmedLine) {
		return p.advance()
	}
	parsedLine, valid := p.parseLine(trimmedLine)
	if !valid {
		return nil, false, &ParsingError{Line: line, Err: fmt.Errorf("invalid instruction format")}
	}

	return parsedLine, true, nil
}

func (p *Parser) handleComments(trimmedLine string) bool {
	if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
		return true
	}
	if p.isCommentBlock {
		if strings.HasSuffix(trimmedLine, "*/") {
			p.isCommentBlock = false
		}
		return true
	}
	if strings.HasPrefix(trimmedLine, "/*") {
		if !strings.HasSuffix(trimmedLine, "*/") {
			p.isCommentBlock = true
		}
		return true
	}
	return false
}

func (p *Parser) parseLine(trimmedLine string) (ParsedLine, bool) {
	if p.isLabel(trimmedLine) {
		return Label{Name: strings.Trim(trimmedLine, "()")}, true
	}
	if p.isAInstruction(trimmedLine) {
		return AInstruction{Symbol: strings.Replace(trimmedLine, "@", "", 1)}, true
	}
	if p.isCInstruction(trimmedLine) {
		return p.parseCInstruction(trimmedLine)
	}
	return nil, false
}

func (p *Parser) isLabel(trimmedLine string) bool {
	return trimmedLine[:1] == "(" && strings.HasSuffix(trimmedLine, ")") && p.isSingleWord(trimmedLine)
}

func (p *Parser) isAInstruction(trimmedLine string) bool {
	return strings.HasPrefix(trimmedLine, "@") && p.isSingleWord(trimmedLine)
}

func (p *Parser) isCInstruction(trimmedLine string) bool {
	firstChar := trimmedLine[:1]
	return strings.ContainsAny(firstChar, "01AD!-M")
}

func (p *Parser) parseCInstruction(trimmedLine string) (ParsedLine, bool) {
	re := regexp.MustCompile(`^(?:(\w+)=)?([^;]+)(?:;(\w+))?$`)
	matches := re.FindStringSubmatch(trimmedLine)
	if len(matches) < 3 {
		return nil, false
	}
	dest, comp, jump := matches[1], matches[2], matches[3]
	return CInstruction{Dest: dest, Comp: comp, Jump: jump}, true
}

func (p *Parser) ParseAndAddSymbols(symbolTable map[string]int) ([]ParsedLine, error) {
	parsedLines, err := p.firstPass(symbolTable)
	if err != nil {
		return nil, err
	}

	if err := p.secondPass(parsedLines, symbolTable); err != nil {
		return nil, err
	}

	return parsedLines, nil
}

func (p *Parser) firstPass(symbolTable map[string]int) ([]ParsedLine, error) {
	var parsedLines []ParsedLine
	instructionCounter := 0
	hasMoreLines := true

	for hasMoreLines {
		parsedLine, hasNext, err := p.advance()
		hasMoreLines = hasNext
		if err != nil {
			return nil, fmt.Errorf("error during first pass: %w", err)
		}

		if parsedLine != nil {
			instructionCounter++
			if label, ok := parsedLine.(Label); ok {
				instructionCounter-- // label is not an instruction
				symbolTable[label.Name] = instructionCounter
			}

			parsedLines = append(parsedLines, parsedLine)
		}
	}

	return parsedLines, nil
}

func (p *Parser) secondPass(parsedLines []ParsedLine, symbolTable map[string]int) error {
	symbolCounter := 16

	for _, parsedLine := range parsedLines {
		if aInst, ok := parsedLine.(AInstruction); ok {
			if _, exists := symbolTable[aInst.Symbol]; !exists {
				if _, err := strconv.Atoi(aInst.Symbol); err != nil {
					symbolTable[aInst.Symbol] = symbolCounter
					symbolCounter++
				}
			}
		}
	}

	return nil
}
