package scanner

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/dywoq/dywoqlang/token"
)

type Scanner struct {
	input          string
	tokenizers     []tokenizer
	line, col, pos int
	setupOn        bool
}

// tokenizer is a function type, representing a tokenizer.
type tokenizer func(rune) (*token.Token, error)

var (
	errNoMatch = errors.New("no match")
)

func New(input string) *Scanner {
	return &Scanner{input: input, tokenizers: []tokenizer{}, line: 1, col: 1, pos: 0, setupOn: false}
}

// current returns the current token.
func (s *Scanner) current() rune {
	if s.pos >= len(s.input) {
		return 0
	}
	return rune(s.input[s.pos])
}

// advance moves the scanner's position forward by i.
func (s *Scanner) advance(i int) {
	if s.pos >= len(s.input) {
		return
	}
	for j := 0; j < i; j++ {
		if s.pos >= len(s.input) {
			return
		}

		char := s.current()
		s.pos++
		s.col++

		switch char {
		case '\n':
			s.col = 1
			s.line++
		case '\r':
			s.col = 1
		}
	}
}

func (s *Scanner) setup() {
	s.setupOn = true
	s.tokenizers = append(s.tokenizers, s.tokenizeComment)
	s.tokenizers = append(s.tokenizers, s.tokenizeType)
	s.tokenizers = append(s.tokenizers, s.tokenizeBaseInstruction)
	s.tokenizers = append(s.tokenizers, s.tokenizeKeyword)
	s.tokenizers = append(s.tokenizers, s.tokenizeString)
	s.tokenizers = append(s.tokenizers, s.tokenizeNumber)
	s.tokenizers = append(s.tokenizers, s.tokenizeSeparator)
	s.tokenizers = append(s.tokenizers, s.tokenizeBoolConstants)
	s.tokenizers = append(s.tokenizers, s.tokenizeIdentifier)
}

func (s *Scanner) reset() {
	s.line = 1
	s.col = 1
	s.pos = 0
}

func (s *Scanner) illegal(line, col, pos int) *token.Token {
	return token.New("", token.KIND_ILLEGAL, token.NewPosition(line, col, pos))
}

// slice returns s.input[start:end], returns error if:
//
// - start is negative;
//
// - start exceeds end;
//
// - end exceeds len(s.input).
func (s *Scanner) slice(start, end int) (string, error) {
	switch {
	case start < 0:
		return "", errors.New("start is negative")
	case start > end:
		return "", errors.New("start exceeds end")
	case end > len(s.input):
		return "", errors.New("end exceeds len(s.input)")
	}
	return s.input[start:end], nil
}

func (s *Scanner) tokenize() (*token.Token, error) {
	if len(s.tokenizers) == 0 {
		return nil, errors.New("there are no tokenizers")
	}
	r := s.current()
	for _, tokenizer := range s.tokenizers {
		t, err := tokenizer(r)
		if err != nil {
			if err == errNoMatch {
				continue
			}
			return nil, err
		}
		return t, err
	}
	return s.illegal(s.line, s.col, s.pos), nil
}

func (s *Scanner) tokenizeSeparator(r rune) (*token.Token, error) {
	if !token.Separators.Is(string(r)) {
		return nil, errNoMatch
	}
	t := token.New(string(r), token.KIND_SEPARATOR, token.NewPosition(s.line, s.col, s.pos))
	s.advance(1)
	return t, nil
}

func (s *Scanner) tokenizeType(r rune) (*token.Token, error) {
	if !unicode.IsLetter(r) {
		return nil, errNoMatch
	}

	startPos, startLine, startCol := s.pos, s.line, s.col
	s.advance(1)
	for unicode.IsLetter(s.current()) {
		s.advance(1)
	}

	strType, err := s.slice(startPos, s.pos)
	if err != nil {
		s.pos, s.line, s.col = startPos, startLine, startCol
		return nil, errNoMatch
	}

	switch strType {
	case "int", "uint":
		if !unicode.IsNumber(s.current()) {
			s.pos, s.line, s.col = startPos, startLine, startCol
			return nil, errNoMatch
		}

		numStart := s.pos
		for unicode.IsNumber(s.current()) {
			s.advance(1)
		}

		numberPart, err := s.slice(numStart, s.pos)
		if err != nil {
			s.pos, s.line, s.col = startPos, startLine, startCol
			return nil, errNoMatch
		}
		result := strType + numberPart

		if !token.Types.Is(result) {
			s.pos, s.line, s.col = startPos, startLine, startCol
			return nil, errNoMatch
		}
		return token.New(result, token.KIND_TYPE, token.NewPosition(startLine, startCol, startPos)), nil

	case "string", "bool":
		return token.New(strType, token.KIND_TYPE, token.NewPosition(startLine, startCol, startPos)), nil

	case "void":
		return token.New(strType, token.KIND_TYPE, token.NewPosition(startLine, startCol, startPos)), nil

	default:
		s.pos, s.line, s.col = startPos, startLine, startCol
		return nil, errNoMatch
	}
}

func (s *Scanner) tokenizeString(r rune) (*token.Token, error) {
	if r != '"' {
		return nil, errNoMatch
	}
	startPos, startLine, startCol := s.pos, s.line, s.col
	s.advance(1)

	for s.current() != '"' && s.pos < len(s.input) {
		s.advance(1)
	}
	if s.pos >= len(s.input) {
		s.pos, s.line, s.col = startPos, startLine, startCol
		return nil, fmt.Errorf("unterminated string at line %d, column %d", startLine, startCol)
	}

	str, err := s.slice(startPos+1, s.pos)
	if err != nil {
		goto restore
	}
	s.advance(1)
	return token.New(str, token.KIND_STRING, token.NewPosition(s.line, s.col, s.pos)), nil

restore:
	s.pos, s.line, s.col = startPos, startLine, startCol
	return nil, err
}

func (s *Scanner) tokenizeKeyword(r rune) (*token.Token, error) {
	if !unicode.IsLetter(r) {
		return nil, errNoMatch
	}
	startPos, startLine, startCol := s.pos, s.line, s.col
	s.advance(1)
	for unicode.IsLetter(s.current()) {
		s.advance(1)
	}
	str, err := s.slice(startPos, s.pos)
	if err != nil {
		goto restore
	}
	if !token.Keywords.Is(str) {
		goto restore
	}
	s.advance(1)
	return token.New(str, token.KIND_KEYWORD, token.NewPosition(startLine, startCol, startPos)), nil

restore:
	s.pos, s.line, s.col = startPos, startLine, startCol
	return nil, errNoMatch
}

func (s *Scanner) tokenizeIdentifier(r rune) (*token.Token, error) {
	if !unicode.IsLetter(r) {
		return nil, errNoMatch
	}
	startPos, startLine, startCol := s.pos, s.line, s.col
	s.advance(1)
	for unicode.IsLetter(s.current()) || unicode.IsDigit(s.current()) || s.current() == '_' {
		if r == '#' || token.Separators.Is(string(s.current())) {
			s.pos, s.line, s.col = startPos, startLine, startCol
			return nil, errNoMatch
		}
		s.advance(1)
	}
	str, err := s.slice(startPos, s.pos)
	if err != nil {
		return nil, err
	}
	if token.IsIdentifier(str) {
		return token.New(str, token.KIND_IDENTIFIER, token.NewPosition(startLine, startCol, startPos)), nil
	}
	s.pos, s.line, s.col = startPos, startLine, startCol
	return nil, errNoMatch
}

func (s *Scanner) tokenizeBaseInstruction(r rune) (*token.Token, error) {
	if !unicode.IsLetter(r) {
		return nil, errNoMatch
	}
	startLine, startCol, startPos := s.line, s.col, s.pos
	s.advance(1)
	for unicode.IsLetter(s.current()) {
		s.advance(1)
	}
	str, err := s.slice(startPos, s.pos)
	if err != nil {
		return nil, err
	}
	if !token.BaseInstructions.Is(str) {
		s.line, s.col, s.pos = startLine, startCol, startPos
		return nil, errNoMatch
	}
	s.advance(1)
	return token.New(str, token.KIND_BASE_INSTRUCTION, token.NewPosition(startLine, startCol, startPos)), nil
}

func (s *Scanner) tokenizeBoolConstants(r rune) (*token.Token, error) {
	if !unicode.IsLetter(r) {
		return nil, errNoMatch
	}
	startLine, startCol, startPos := s.line, s.col, s.pos
	s.advance(1)
	for unicode.IsLetter(s.current()) {
		s.advance(1)
	}
	str, err := s.slice(startPos, s.pos)
	if err != nil {
		return nil, err
	}
	if !token.BoolConstants.Is(str) {
		s.line, s.col, s.pos = startLine, startCol, startPos
		return nil, errNoMatch
	}
	s.advance(1)
	return token.New(str, token.KIND_BOOL_CONSTANT, token.NewPosition(startLine, startCol, startPos)), nil
}

func (s *Scanner) tokenizeNumber(r rune) (*token.Token, error) {
	startLine, startCol, startPos := s.line, s.col, s.pos
	hasDot := false

	if r == '-' {
		if s.pos+1 >= len(s.input) || !unicode.IsNumber(rune(s.input[s.pos+1])) {
			return nil, errNoMatch
		}
		s.advance(1)
		r = s.current()
	}

	if !unicode.IsNumber(r) {
		return nil, errNoMatch
	}

	for {
		c := s.current()
		if unicode.IsNumber(c) {
			s.advance(1)
			continue
		}
		if c == '.' {
			if hasDot {
				break
			}
			hasDot = true
			s.advance(1)
			if !unicode.IsNumber(s.current()) {
				s.pos, s.line, s.col = startPos, startLine, startCol
				return nil, fmt.Errorf("invalid float literal at line %d, column %d", startLine, startCol)
			}
			continue
		}
		break
	}

	str, err := s.slice(startPos, s.pos)
	if err != nil {
		s.pos, s.line, s.col = startPos, startLine, startCol
		return nil, err
	}

	if hasDot {
		return token.New(str, token.KIND_FLOAT, token.NewPosition(startLine, startCol, startPos)), nil
	}
	return token.New(str, token.KIND_INTEGER, token.NewPosition(startLine, startCol, startPos)), nil
}

func (s *Scanner) tokenizeComment(r rune) (*token.Token, error) {
	if r != '#' {
		return nil, errNoMatch
	}
	startLine, startCol, startPos := s.line, s.col, s.pos
	s.advance(1)

	commentStart := s.pos
	for c := s.current(); c != 0 && c != '\n' && c != '\r'; c = s.current() {
		s.advance(1)
	}

	str, err := s.slice(commentStart, s.pos)
	if err != nil {
		s.pos, s.line, s.col = startPos, startLine, startCol
		return nil, err
	}

	str = strings.TrimLeft(str, " \t")

	return token.New(str, token.KIND_COMMENT, token.NewPosition(startLine, startCol, startPos)), nil
}
func (s *Scanner) skipWhitespace() {
	for unicode.IsSpace(s.current()) {
		s.advance(1)
	}
}

func (s *Scanner) Scan() ([]*token.Token, error) {
	s.reset()
	if !s.setupOn {
		s.setup()
	}
	result := []*token.Token{}
	var scanErr error

	for s.pos < len(s.input) {
		s.skipWhitespace()
		if s.current() == 0 {
			break
		}
		t, err := s.tokenize()
		if err != nil {
			scanErr = err
			goto cleanup
		}
		if t.Kind == token.KIND_ILLEGAL {
			scanErr = fmt.Errorf("illegal character at line %d, column %d", t.Position.Line, t.Position.Column)
			goto cleanup
		}
		result = append(result, t)
	}
	result = append(result, token.New("", token.KIND_EOF, token.NewPosition(s.line, s.col, s.pos)))
	return result, nil

cleanup:
	result = append(result, token.New("", token.KIND_EOF, token.NewPosition(s.line, s.col, s.pos)))
	return result, scanErr
}
