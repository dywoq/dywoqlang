package token

import "unicode"

// Kind represents the token kind.
type Kind string

// Map represents the map of tokens.
type Map map[string]Kind

// Position is a token position.
type Position struct {
	Line     int `json:"line"`
	Column   int `json:"column"`
	Position int `json:"position"`
}

type Token struct {
	Literal  string    `json:"literal"`
	Kind     Kind      `json:"kind"`
	Position *Position `json:"position"`
}

const (
	Illegal         Kind = "illegal"
	Keyword         Kind = "keyword"
	Type            Kind = "type"
	Separator       Kind = "separator"
	Identifier      Kind = "identifier"
	Float           Kind = "float"
	Integer         Kind = "integer"
	String          Kind = "string"
	BaseInstruction Kind = "base_instruction"
	Special         Kind = "special"
	BinaryOperator  Kind = "binary_operator"
	BoolConstant    Kind = "bool_constant"
	Eof             Kind = "eof"
)

var (
	KeywordsMap = Map{
		"export":    Keyword,
		"import":    Keyword,
		"declare":   Keyword,
		"link":      Keyword,
		"consteval": Keyword,
		"copy":      Keyword,
		"meta":      Keyword,
		"array":     Keyword,
	}

	SpecialMap = Map{
		"nil": Special,
	}

	SeparatorsMap = Map{
		",": Separator,
		"{": Separator,
		"}": Separator,
		"(": Separator,
		")": Separator,
		";": Separator,
		"[": Separator,
		"]": Separator,
		":": Separator,
	}

	TypesMap = Map{
		"str":  Type,
		"bool": Type,
		"i8":   Type,
		"i16":  Type,
		"i32":  Type,
		"i64":  Type,
		"u8":   Type,
		"u16":  Type,
		"u32":  Type,
		"u64":  Type,
		"void": Type,
		"f32":  Type,
		"f64":  Type,
	}

	BaseInstructionsMap = Map{
		"stdout": BaseInstruction,
		"stderr": BaseInstruction,
		"mov":    BaseInstruction,
		"ret":    BaseInstruction,
		"add":    BaseInstruction,
		"div":    BaseInstruction,
		"mul":    BaseInstruction,
		"sub":    BaseInstruction,
	}

	BinaryOperatorsMap = Map{
		"+": BinaryOperator,
		"-": BinaryOperator,
		"/": BinaryOperator,
		"*": BinaryOperator,
	}

	BoolConstantsMap = Map{
		"true":  BoolConstant,
		"false": BoolConstant,
	}
)

// IsIdentifier reports whether value is a valid identifier,
// meaning value can't be keyword, separator, type,
// contain hash, left and right paren, slash or start with the digit.
func IsIdentifier(value string) bool {
	switch {
	case KeywordsMap.Is(value), SeparatorsMap.Is(value), TypesMap.Is(value):
		return false
	}
	for i, r := range value {
		// immediately check if first rune is digit
		if i == 0 && unicode.IsDigit(r) {
			return false
		}
		switch r {
		case '#':
			return false
		case '/':
			return false
		case '(':
			return false
		case ')':
			return false
		}
	}
	return true
}

// NewPosition returns a pointer to new token position.
func NewPosition(line, column, position int) *Position {
	return &Position{line, column, position}
}

// NewToken returns a pointer to new token.
func NewToken(literal string, kind Kind, position *Position) *Token {
	return &Token{Literal: literal, Kind: kind, Position: position}
}

// Equal reports whether x and y are equal by their literal and kind.
// If they're nil, it returns false.
func Equal(x *Token, y *Token) bool {
	if x == nil || y == nil {
		return false
	}
	if x.Literal == y.Literal && x.Kind == y.Kind {
		return true
	}
	return false
}

// Is reports whether the value exists in the tokens map.
func (m Map) Is(value string) bool {
	_, ok := m[value]
	return ok
}
