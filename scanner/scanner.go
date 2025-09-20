package scanner

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/dywoq/dywoqlang/token"
)

type Scanner struct {
	input             string
	pos, column, line int
	tokenizers        []tokenizerFunc
	setupOn           bool
}

func New(input string) *Scanner {
	return &Scanner{input: input, pos: 0, column: 1, line: 1, tokenizers: []tokenizerFunc{}, setupOn: false}
}

type tokenizerFunc func(r rune) (*token.Token, error)

var (
	errNoMatch = errors.New("no match")
)

// current returns the current character,
// if s.pos is not greater than or equal to len(s.input)
func (s *Scanner) current() rune {
	if s.pos >= len(s.input) {
		return 0 // Return null character or EOF
	}
	return rune(s.input[s.pos])
}

// peek returns the future character,
// if s.pos+1 is not greater than or equal to len(s.input).
func (s *Scanner) peek() rune {
	if s.pos+1 >= len(s.input) {
		return 0
	}
	return rune(s.input[s.pos+1])
}

// advance advances to the next position by i,
// if i is not negative and s.pos+i is not greater than len(s.input).
func (s *Scanner) advance() {
	if s.pos < len(s.input) {
		if s.current() == '\n' {
			s.line++
			s.column = 1
		} else {
			s.column++
		}
		s.pos++
	}
}

// slice returns a slice of string with the start and the end,
// returns error if:
//
// - if start is negative;
//
// - if start exceeds end;
//
// - if end exceeds len(s.input)
func (s *Scanner) slice(start, end int) (string, error) {
	switch {
	case start < 0:
		return "", errors.New("start position is negative")
	case start > end:
		return "", errors.New("start position is after end position")
	case end > len(s.input):
		return "", errors.New("end position exceeds input length")
	}
	return s.input[start:end], nil
}

func (s *Scanner) illegal(line, column int) *token.Token {
	return token.New("", token.KIND_ILLEGAL, line, column)
}

func (s *Scanner) tokenizeType(r rune) (*token.Token, error) {
	if !unicode.IsLetter(r) {
		return nil, errNoMatch
	}

	startPos, startLine, startCol := s.pos, s.line, s.column

	s.advance()
	for unicode.IsLetter(s.current()) {
		s.advance()
	}

	strType, err := s.slice(startPos, s.pos)
	if err != nil {
		s.pos, s.line, s.column = startPos, startLine, startCol
		return nil, err
	}

	switch strType {
	case "int", "uint":
		if !unicode.IsNumber(s.current()) {
			return nil, errNoMatch
		}

		numStart := s.pos
		for unicode.IsNumber(s.current()) {
			s.advance()
		}

		numberPart, err := s.slice(numStart, s.pos)
		if err != nil {
			return nil, err
		}
		result := strType + numberPart

		if !token.Types.Is(result) {
			return nil, fmt.Errorf("wrong type: %s", result)
		}
		return token.New(result, token.KIND_TYPE, s.line, s.column), nil

	case "string", "bool":
		return token.New(strType, token.KIND_TYPE, s.line, s.column), nil

	default:
		return nil, errNoMatch
	}
}

func (s *Scanner) tokenizeSeparator(r rune) (*token.Token, error) {
	if !token.Separators.Is(string(r)) {
		return nil, errNoMatch
	}
	t := token.New(string(r), token.KIND_SEPARATOR, s.line, s.column)
	s.advance()
	return t, nil
}

func (s *Scanner) tokenizeKeyword(r rune) (*token.Token, error) {
	if !unicode.IsLetter(r) {
		return nil, errNoMatch
	}
	startPos, startLine, startCol := s.pos, s.line, s.column
	s.advance()
	for unicode.IsLetter(s.current()) {
		s.advance()
	}
	str, err := s.slice(startPos, s.pos)
	if err != nil {
		s.pos, s.line, s.column = startPos, startLine, startCol
		return nil, err
	}
	if !token.Keywords.Is(str) {
		s.pos, s.line, s.column = startPos, startLine, startCol
		return nil, errNoMatch
	}
	s.advance()
	return token.New(str, token.KIND_KEYWORD, s.line, s.column), nil
}

func (s *Scanner) tokenizeString(r rune) (*token.Token, error) {
	if r != '"' {
		return nil, errNoMatch
	}
	startPos, startLine, startCol := s.pos, s.line, s.column
	s.advance()
	for s.current() != '"' {
		s.advance()
	}
	str, err := s.slice(startPos+1, s.pos)
	if err != nil {
		s.pos, s.line, s.column = startPos, startLine, startCol
		return nil, err
	}
	s.advance()
	return token.New(str, token.KIND_STRING, s.line, s.column), nil
}

// setup adds default tokenizers and sets s.setupOn to true.
func (s *Scanner) setup() {
	s.tokenizers = append(s.tokenizers, s.tokenizeKeyword)
	s.tokenizers = append(s.tokenizers, s.tokenizeType)
	s.tokenizers = append(s.tokenizers, s.tokenizeSeparator)
	s.tokenizers = append(s.tokenizers, s.tokenizeString)
	s.setupOn = true
}

// tokenize returns a token, got from the tokenizers.
//
// If len(s.tokenizers) is 0, it returns nil and error.
//
// If tokenizer returns errNoMatch, the scanner continues iterate over tokenizers slice.
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
	return s.illegal(s.line, s.column), nil
}

func (s *Scanner) reset() {
	s.pos = 0
	s.line = 1
	s.column = 1
}

func (s *Scanner) skipWhitespace() {
	for s.current() != 0 && unicode.IsSpace(s.current()) {
		s.advance()
	}
}

func (s *Scanner) Scan() ([]*token.Token, error) {
	s.reset()
	if !s.setupOn {
		s.setup()
	}
	var tokens []*token.Token
	for s.pos < len(s.input) {
		s.skipWhitespace()
		if s.current() == 0 {
			break
		}
		t, err := s.tokenize()
		if err != nil {
			return nil, err
		}

		if t.Kind == token.KIND_ILLEGAL {
			return nil, fmt.Errorf("illegal character at line %d, column %d", t.Line, t.Column)
		}

		tokens = append(tokens, t)
	}
	tokens = append(tokens, token.New("", token.KIND_EOF, s.line, s.column))
	return tokens, nil
}
