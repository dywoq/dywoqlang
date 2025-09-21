package token

import "unicode"

// Kind is a alias of string.
// Represents token kind.
type Kind string

// Map is alias of map[string]Kind.
// Used to create maps of keywords.
type Map map[string]Kind

// Token representing a token with literal,
// kind, line and column.
type Token struct {
	Literal string
	Kind    Kind
	Line    int
	Column  int
}

const (
	KIND_KEYWORD          Kind = "keyword"
	KIND_IDENTIFIER       Kind = "identifier"
	KIND_BASE_INSTRUCTION Kind = "instruction"
	KIND_INTEGER          Kind = "integer"
	KIND_STRING           Kind = "string"
	KIND_SEPARATOR        Kind = "separator"
	KIND_TYPE             Kind = "type"
	KIND_BOOL_CONSTANT    Kind = "bool_constant"
)

var (
	Keywords = Map{
		"export": KIND_KEYWORD,
	}

	BaseInstructions = Map{
		"mov":   KIND_BASE_INSTRUCTION,
		"add":   KIND_BASE_INSTRUCTION,
		"div":   KIND_BASE_INSTRUCTION,
		"mul":   KIND_BASE_INSTRUCTION,
		"sub":   KIND_BASE_INSTRUCTION,
		"ret":   KIND_BASE_INSTRUCTION,
		"write": KIND_BASE_INSTRUCTION,
	}

	Separators = Map{
		",": KIND_SEPARATOR,
		"(": KIND_SEPARATOR,
		")": KIND_SEPARATOR,
		":": KIND_SEPARATOR,
		"[": KIND_SEPARATOR,
		"]": KIND_SEPARATOR,
	}

	Types = Map{
		"int8":   KIND_TYPE,
		"int16":  KIND_TYPE,
		"int32":  KIND_TYPE,
		"int64":  KIND_TYPE,
		"uint8":  KIND_TYPE,
		"uint16": KIND_TYPE,
		"uint32": KIND_TYPE,
		"uint64": KIND_TYPE,
		"bool":   KIND_TYPE,
		"string": KIND_TYPE,
	}

	BoolConstants = Map{
		"true":  KIND_BOOL_CONSTANT,
		"false": KIND_BOOL_CONSTANT,
	}
)

// Is checks if value is present in the map.
func (m Map) Is(value string) bool {
	_, ok := m[value]
	return ok
}

// New returns a new pointer to Token.
func New(literal string, kind Kind, line, column int) *Token {
	return &Token{literal, kind, line, column}
}

// IsIdentifier reports whether the value is a valid identifier,
// returns false if:
//
// - value is keyword, base instruction, separator, type or bool constant;
//
// - first character is number;
//
// - the characters of value include #, / or any of separators.
func IsIdentifier(value string) bool {
	switch {
	case len(value) == 0,
		Keywords.Is(value),
		BaseInstructions.Is(value),
		Separators.Is(value),
		Types.Is(value),
		BoolConstants.Is(value):
		return false
	}
	for i, r := range value {
		if i == 0 && unicode.IsNumber(r) {
			return false
		}
		if Separators.Is(string(r)) {
			return false
		}
		if r == '#' || r == '/' {
			return false
		}
	}
	return true
}
