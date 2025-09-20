package token

import (
	"encoding/json"
	"unicode"
)

// Kind represents Token Kind.
type Kind string

// Map is an alias of map[string]Kind, representing a map of keywords.
type Map map[string]Kind

type Token struct {
	Literal string `json:"literal"`
	Kind    Kind   `json:"kind"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
}

const (
	KIND_INSTRUCTION         Kind = "instruction"
	KIND_IDENTIFIER          Kind = "identifier"
	KIND_TYPE                Kind = "type"
	KIND_SEPARATOR           Kind = "separator"
	KIND_KEYWORD             Kind = "keyword"
	KIND_BOOL_CONSTANT       Kind = "bool_constant"
	KIND_ARITHMETIC_OPERATOR Kind = "operator"
	KIND_NUMBER              Kind = "number"
	KIND_STRING              Kind = "string"
	KIND_EOF                 Kind = "eof"
	KIND_ILLEGAL             Kind = "illegal"
)

var (
	Types = Map{
		"int8":   KIND_TYPE,
		"int16":  KIND_TYPE,
		"int32":  KIND_TYPE,
		"int64":  KIND_TYPE,
		"uint8":  KIND_TYPE,
		"uint16": KIND_TYPE,
		"uint32": KIND_TYPE,
		"uint64": KIND_TYPE,
		"string": KIND_TYPE,
		"bool":   KIND_TYPE,
	}

	Keywords = Map{
		"export": KIND_KEYWORD,
	}

	Separators = Map{
		",": KIND_SEPARATOR,
		"{": KIND_SEPARATOR,
		"}": KIND_SEPARATOR,
		"(": KIND_SEPARATOR,
		")": KIND_SEPARATOR,
	}

	BoolConstants = Map{
		"true":  KIND_BOOL_CONSTANT,
		"false": KIND_BOOL_CONSTANT,
	}

	Instructions = Map{
		"add":      KIND_INSTRUCTION,
		"minus":    KIND_INSTRUCTION,
		"divide":   KIND_INSTRUCTION,
		"multiply": KIND_INSTRUCTION,
		"mov":      KIND_INSTRUCTION,
		"make":     KIND_INSTRUCTION,
	}
)

// New return a pointer to Token.
func New(literal string, kind Kind, line int, column int) *Token {
	return &Token{literal, kind, line, column}
}

// String converts Token into string presentation.
func (t *Token) String() string {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// Is checks if value is present in map.
func (m Map) Is(value string) bool {
	if value == "" {
		return false
	}
	_, ok := m[value]
	return ok
}

// IsArithmeticOperator reports whether r is
//
//	'+', '-', '/', '*'
func IsArithmeticOperator(r rune) bool {
	switch r {
	case '+', '-', '/', '*':
		return true
	}
	return false
}

// IsIdentifier reports whether the value is valid identifier, returns false if:
//
// - value is empty;
//
// - value is type, keyword or base instruction;
//
// - value is separator;
//
// - value is #, or arithmetic operator.
func IsIdentifier(value string) bool {
	if len(value) == 0 || Types.Is(value) || Keywords.Is(value) || Instructions.Is(value) {
		return false
	}
	for index, r := range value {
		switch {
		case index == 0 && unicode.IsNumber(r):
			return false
		case index == 0 && r == '_':
			return false
		case r == '#' || IsArithmeticOperator(r):
			return false
		case Separators.Is(string(r)):
			return false
		}
	}
	return true
}
