package scanner

import (
	"errors"
	"fmt"
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
	s.pos += i
	s.col++
	if s.current() == '\n' {
		s.col = 1
		s.line++
	}
}

func (s *Scanner) setup() {
	s.setupOn = true
	s.tokenizers = append(s.tokenizers, s.tokenizeKeyword)
	s.tokenizers = append(s.tokenizers, s.tokenizeIdentifier)
	s.tokenizers = append(s.tokenizers, s.tokenizeString)
	s.tokenizers = append(s.tokenizers, s.tokenizeSeparator)
	s.tokenizers = append(s.tokenizers, s.tokenizeType)
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
		goto restore
	}

	switch strType {
	case "int", "uint":
		if !unicode.IsNumber(s.current()) {
			goto restore
		}

		numStart := s.pos
		for unicode.IsNumber(s.current()) {
			s.advance(1)
		}

		numberPart, err := s.slice(numStart, s.pos)
		if err != nil {
			goto restore
		}
		result := strType + numberPart

		if !token.Types.Is(result) {
			goto restore
		}
		return token.New(result, token.KIND_TYPE, token.NewPosition(s.line, s.col, s.pos)), nil

	case "string", "bool":
		return token.New(strType, token.KIND_TYPE, token.NewPosition(s.line, s.col, s.pos)), nil

	default:
		goto restore
	}

restore:
	s.pos, s.line, s.col = startPos, startLine, startCol
	return nil, err
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
	return token.New(str, token.KIND_KEYWORD, token.NewPosition(s.line, s.col, s.pos)), nil

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
		return token.New(str, token.KIND_IDENTIFIER, token.NewPosition(s.pos, s.col, s.pos)), nil
	}
	s.pos, s.line, s.col = startPos, startLine, startCol
	return nil, errNoMatch
}

func (s *Scanner) skipWhitespace() {
	if unicode.IsSpace(s.current()) {
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
