package scanner

import (
	"errors"
	"fmt"
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

// TokenizeString tokenizes a string.
//
// If the string was unterminated, the tokenizer
// returns an error.
//
// Returns an error if scanner reached End Of File (EOF).
func TokenizeString(c Context) (*token.Token, error) {
	if c.Eof() {
		return nil, ErrEof
	}
	if r, _ := c.Current(); r != '"' {
		return nil, ErrNoMatch
	}

	err := c.Advance(1)
	if err != nil {
		return nil, err
	}

	start := c.Position().Position
	for {
		if c.Eof() {
			return nil, fmt.Errorf("unterminated string at line %d, column %d", c.Position().Line, c.Position().Column)
		}

		r, _ := c.Current()
		if r == '\n' {
			return nil, fmt.Errorf("unterminated string at line %d, column %d", c.Position().Line, c.Position().Column)
		}

		if r == '"' {
			break
		}

		if err := c.Advance(1); err != nil {
			return nil, err
		}
	}

	substr, err := c.Slice(start, c.Position().Position)
	if err != nil {
		return nil, err
	}

	err = c.Advance(1)
	if err != nil {
		return nil, err
	}

	return c.New(substr, token.String), nil
}

// TokenizeKeyword tokenizes a keyword.
//
// Returns an error if the scanner reached End Of File (EOF).
//
// If it doesn't match, the function returns ErrNoMatch
// and advances to the initial position.
func TokenizeKeyword(c Context) (*token.Token, error) {
	if c.Eof() {
		return nil, ErrEof
	}
	if r, _ := c.Current(); !unicode.IsLetter(r) {
		return nil, ErrNoMatch
	}

	start := c.Position().Position
	pos := start
	for {
		if pos >= len(c.Input()) {
			break
		}
		r := rune(c.Input()[pos])
		if !unicode.IsLetter(r) {
			break
		}
		pos++
	}

	substr, err := c.Slice(start, pos)
	if err != nil {
		return nil, err
	}

	if !token.KeywordsMap.Is(substr) {
		return nil, ErrNoMatch
	}

	if err := c.Advance(pos - start); err != nil {
		return nil, err
	}

	return c.New(substr, token.Keyword), nil
}

// TokenizeSeparator tokenizes a separator.
//
// Returns an error if the scanner reached End Of File (EOF).
//
// If it doesn't match, the function returns ErrNoMatch
// and advances to the initial position.
func TokenizeSeparator(c Context) (*token.Token, error) {
	if c.Eof() {
		return nil, ErrEof
	}
	r, _ := c.Current()
	if !token.SeparatorsMap.Is(string(r)) {
		return nil, ErrNoMatch
	}

	if err := c.Advance(1); err != nil {
		return nil, err
	}

	return c.New(string(r), token.Separator), nil
}

// TokenizeSpecial tokenizes a special names.
//
// Returns an error if the scanner reached End Of File (EOF).
//
// If it doesn't match, the function returns ErrNoMatch
// and advances to the initial position.
func TokenizeSpecial(c Context) (*token.Token, error) {
	if c.Eof() {
		return nil, ErrEof
	}
	if r, _ := c.Current(); !unicode.IsLetter(r) {
		return nil, ErrNoMatch
	}

	start := c.Position().Position
	for {
		err := c.Advance(1)
		if err != nil {
			return nil, err
		}

		if c.Eof() {
			break
		}

		r, _ := c.Current()
		if !unicode.IsLetter(r) {
			break
		}
	}

	substr, err := c.Slice(start, c.Position().Position)
	if err != nil {
		return nil, err
	}

	if !token.SpecialMap.Is(substr) {
		c.Position().Position = start
		return nil, ErrNoMatch
	}
	return c.New(substr, token.Special), nil
}

// TokenizeBaseInstruction tokenizes base instructions.
//
// Returns an error if the scanner reached End Of File (EOF).
//
// If it doesn't match, the function returns ErrNoMatch
// and advances to the initial position.
func TokenizeBaseInstruction(c Context) (*token.Token, error) {
	if c.Eof() {
		return nil, ErrEof
	}
	if r, _ := c.Current(); !unicode.IsLetter(r) {
		return nil, ErrNoMatch
	}

	start := c.Position().Position
	for {
		err := c.Advance(1)
		if err != nil {
			return nil, err
		}

		if c.Eof() {
			break
		}

		r, _ := c.Current()
		if !unicode.IsLetter(r) {
			break
		}
	}

	substr, err := c.Slice(start, c.Position().Position)
	if err != nil {
		return nil, err
	}

	if !token.BaseInstructionsMap.Is(substr) {
		c.Position().Position = start
		return nil, ErrNoMatch
	}
	return c.New(substr, token.BaseInstruction), nil
}

// TokenizeBinaryOperator tokenizes binary operators
//
// Returns an error if the scanner reached End Of File (EOF).
//
// If it doesn't match, the function returns ErrNoMatch
// and advances to the initial position.
func TokenizeBinaryOperator(c Context) (*token.Token, error) {
	if c.Eof() {
		return nil, ErrEof
	}
	r, _ := c.Current()
	if !token.BinaryOperatorsMap.Is(string(r)) {
		return nil, ErrNoMatch
	}

	if err := c.Advance(1); err != nil {
		return nil, err
	}
	return c.New(string(r), token.BinaryOperator), nil
}

// TokenizeIdentifier tokenizes identifiers.
//
// Returns an error if the scanner reached End Of File (EOF).
//
// If it doesn't match, the function returns ErrNoMatch
// and advances to the initial position.
func TokenizeIdentifier(c Context) (*token.Token, error) {
	if c.Eof() {
		return nil, ErrEof
	}

	r, _ := c.Current()
	if !unicode.IsLetter(r) && r != '_' {
		return nil, ErrNoMatch
	}

	start := c.Position().Position
	pos := start

	for pos < len(c.Input()) {
		r := rune(c.Input()[pos])
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			break
		}
		pos++
	}

	substr, err := c.Slice(start, pos)
	if err != nil {
		return nil, err
	}

	if !token.IsIdentifier(substr) {
		return nil, ErrNoMatch
	}

	if err := c.Advance(pos - start); err != nil {
		return nil, err
	}

	return c.New(substr, token.Identifier), nil
}

// TokenizeTypes tokenizes types.
//
// Returns an error if the scanner reached End Of File (EOF).
//
// If it doesn't match, the function returns ErrNoMatch
// and advances to the initial position.
func TokenizeTypes(c Context) (*token.Token, error) {
	if c.Eof() {
		return nil, ErrEof
	}

	r, _ := c.Current()
	if !unicode.IsLetter(r) {
		return nil, ErrNoMatch
	}

	start := c.Position().Position
	for {
		err := c.Advance(1)
		if err != nil {
			return nil, err
		}

		if c.Eof() {
			break
		}

		r, _ := c.Current()
		if !unicode.IsLetter(r) {
			break
		}
	}

	substr, err := c.Slice(start, c.Position().Position)
	if err != nil {
		return nil, err
	}

	switch substr {
	case "str", "void", "bool":
		return c.New(substr, token.Type), nil

	case "i", "u", "f":
		for {
			c.Advance(1)
			if c.Eof() {
				break
			}
			r, _ = c.Current()
			if !unicode.IsNumber(r) {
				break
			}
		}

		substr, err = c.Slice(start, c.Position().Position)
		if err != nil {
			return nil, err
		}

		if !token.TypesMap.Is(substr) {
			return nil, fmt.Errorf("wrong numeric type: %s", substr)
		}

		return c.New(substr, token.Type), nil

	case "[]":
		return c.New(substr, token.Type), nil
	}

	c.Position().Position = start
	return nil, ErrNoMatch
}

// TokenizeTypes tokenizes bool constants.
//
// Returns an error if the scanner reached End Of File (EOF).
//
// If it doesn't match, the function returns ErrNoMatch
// and advances to the initial position.
func TokenizeBoolConstant(c Context) (*token.Token, error) {
	if c.Eof() {
		return nil, ErrEof
	}

	r, _ := c.Current()
	if !unicode.IsLetter(r) {
		return nil, ErrNoMatch
	}

	start := c.Position().Position
	for {
		err := c.Advance(1)
		if err != nil {
			return nil, err
		}

		if c.Eof() {
			break
		}

		r, _ := c.Current()
		if !unicode.IsLetter(r) {
			break
		}
	}

	substr, err := c.Slice(start, c.Position().Position)
	if err != nil {
		return nil, err
	}

	if !token.BoolConstantsMap.Is(substr) {
		c.Position().Position = start
		return nil, ErrNoMatch
	}

	return c.New(substr, token.BoolConstant), nil
}

// TokenizeTypes tokenizes comments that start with hash (#).
//
// Returns an error if the scanner reached End Of File (EOF).
//
// If it doesn't match, the function returns ErrNoMatch
// and advances to the initial position.
func TokenizeComment(c Context) (*token.Token, error) {
	r, err := c.Current()
	if err != nil {
		return nil, err
	}

	if r != '#' {
		return nil, ErrNoMatch
	}

	start := c.Position().Position
	for {
		if c.Eof() {
			break
		}
		r, _ := c.Current()
		if r == '\n' {
			break
		}
		c.Advance(1)
	}

	end := c.Position().Position
	literal, _ := c.Slice(start, end)
	return c.New(literal, token.Comment), nil
}
