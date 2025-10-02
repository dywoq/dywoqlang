package scanner

import (
	"errors"
	"unicode"

	"github.com/dywoq/dywoqlang/token"
)

type Scanner struct {
	input    string
	position *token.Position

	setupOn    bool
	tokenizers []TokenizerFunc
}

// New returns a new pointer to Scanner.
func New() *Scanner {
	return &Scanner{input: "", position: &token.Position{Line: 1, Column: 1}, setupOn: false, tokenizers: make([]TokenizerFunc, 0)}
}

// Advance advances to the next position by n
//
// If newline character is met,
// scanner increases line and column beside the current position.
//
// Returns an error if scanner reached End Of File (EOF),
// or the current position+n will make the position out of the input.
//
// Does nothing if n is zero.
func (s *Scanner) Advance(n int) error {
	if s.Eof() {
		return ErrEof
	}
	if s.position.Position+n > len(s.input) {
		return errors.New("the current position+n is higher than the length")
	}
	if n == 0 {
		return nil
	}

	for i := 0; i < n; i++ {
		s.position.Position++
		if s.Eof() {
			return nil
		}
		r, _ := s.Current()
		if r == '\n' {
			s.position.Line++
			s.position.Column = 1
		} else {
			s.position.Column++
		}
	}
	return nil
}

// Slice takes a substring from the input surrounded by start and end.
//
// Returns an error if start is negative, start is higher than end,
// or end is higher than the input.
func (s *Scanner) Slice(start, end int) (string, error) {
	switch {
	case start > end:
		return "", errors.New("start is higher than the end")
	case start < 0:
		return "", errors.New("start is negative")
	case end > len(s.input):
		return "", errors.New("end is higher than the input")
	}
	return s.input[start:end], nil
}

// Eof returns true if scanner reached End Of File (EOF).
func (s *Scanner) Eof() bool {
	return s.position.Position >= len(s.input)
}

// Peek returns the future character.
//
// Returns an error if the scanner reached EOF,
// or the current position+1 will make the scanner position out of the input.
func (s *Scanner) Peek() (rune, error) {
	switch {
	case s.Eof():
		return 0, errors.New("reached eof")
	case s.position.Position+1 >= len(s.input):
		return 0, errors.New("current position+1 is higher than the input")
	}
	return rune(s.input[s.position.Position+1]), nil
}

// Current returns the current character.
// Returns an error if the scanner reached EOF.
func (s *Scanner) Current() (rune, error) {
	if s.Eof() {
		return 0, errors.New("reached eof")
	}
	return rune(s.input[s.position.Position]), nil
}

// Position returns a current position.
func (s *Scanner) Position() *token.Position {
	return s.position
}

// New returns a new token.
//
// The only difference from token.NewToken is
// that scanner automatically inserts the position.
func (s *Scanner) New(literal string, kind token.Kind) *token.Token {
	return token.NewToken(literal, kind, s.position)
}

// Scan scans input, turning characters into the tokens.
//
// If there are no tokenizers even after the setup of default tokenizers (which is rare),
// Scan returns an error.
//
// If input is empty, the function returns an error.
func (s *Scanner) Scan(input string) ([]*token.Token, error) {
	if len(input) == 0 {
		return nil, errors.New("input is empty")
	}
	s.setup()
	if len(s.tokenizers) == 0 {
		return nil, errors.New("there are no tokenizers")
	}
	s.reset(input)

	var result []*token.Token
	for !s.Eof() {
		err := s.skipWhitespace()
		if err != nil {
			return nil, err
		}
		t, err := s.tokenize()
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	result = append(result, token.NewToken("", token.Eof, s.position))
	return result, nil
}

func (s *Scanner) skipWhitespace() error {
	for {
		if s.Eof() {
			return nil
		}
		r, _ := s.Current()
		if !unicode.IsSpace(r) {
			return nil
		}
		if err := s.Advance(1); err != nil {
			return err
		}
	}
}

func (s *Scanner) reset(input string) {
	s.input = input
	s.position.Position = 0
	s.position.Line = 1
	s.position.Column = 1
}

func (s *Scanner) tokenize() (*token.Token, error) {
	for _, t := range s.tokenizers {
		tok, err := t(s)
		if err != nil {
			if errors.Is(err, ErrNoMatch) {
				continue
			}
			return nil, err
		}
		return tok, nil
	}
	return token.NewToken("illegal", token.Illegal, s.position), nil
}

func (s *Scanner) setup() {
	if !s.setupOn {
		s.tokenizers = []TokenizerFunc{
			TokenizeNumber,
			TokenizeString,
		}
		s.setupOn = true
	}
}
