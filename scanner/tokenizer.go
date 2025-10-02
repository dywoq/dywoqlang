package scanner

import (
	"errors"
	"unicode"

	"github.com/dywoq/dywoqlang/token"
)

// TokenizerFunc is a function alias for functions what tokenize symbols,
// using c.
type TokenizerFunc func(c Context) (*token.Token, error)

// TokenizeNumber tokenizes a number.
//
// If current symbol doesn't meet the requirements (such as not number),
// then it returns an error ErrNoMatch.
//
// Returns an error if scanner reached End Of File (EOF).
//
// If there's a point after the number, the tokenizer will mark it as a float number,
// if there's still a point, but after it there's no number, the tokenizer will return an error.
func TokenizeNumber(c Context) (*token.Token, error) {
	if c.Eof() {
		return nil, ErrEof
	}
	if r, _ := c.Current(); !unicode.IsNumber(r) {
		return nil, ErrNoMatch
	}

	start := c.Position().Position
	r, _ := c.Current()

	for unicode.IsNumber(r) {
		if err := c.Advance(1); err != nil {
			return nil, err
		}
		if c.Eof() {
			break
		}
		r, _ = c.Current()
	}

	if r, _ := c.Current(); r != '.' {
		substr, err := c.Slice(start, c.Position().Position)
		if err != nil {
			return nil, err
		}
		return c.New(substr, token.Integer), nil
	}

	if err := c.Advance(1); err != nil {
		return nil, err
	}
	if c.Eof() {
		return nil, errors.New("expected a number after point")
	}

	r, _ = c.Current()
	if !unicode.IsNumber(r) {
		return nil, errors.New("expected a number after point")
	}

	for unicode.IsNumber(r) {
		if err := c.Advance(1); err != nil {
			return nil, err
		}
		if c.Eof() {
			break
		}
		r, _ = c.Current()
	}

	substr, err := c.Slice(start, c.Position().Position)
	if err != nil {
		return nil, err
	}
	return c.New(substr, token.Float), nil
}
