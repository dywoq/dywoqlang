package parser

import (
	"errors"

	"github.com/dywoq/dywoqlang/token"
)

// MiniParser is for parsing tokens into AST nodes.
// If it returns ErrNoMatch, it means the parser will try other mini-parsers.
type MiniParser func(c Context, t *token.Token) (Node, error)

// ParseInstruction parses the instruction into the AST node.
func ParseInstruction(c Context, t *token.Token) (Node, error) {
	if t.Kind != token.KIND_BASE_INSTRUCTION && !token.IsIdentifier(t.Literal) {
		return nil, ErrNoMatch
	}
	identifier := t
 	if identifier.Kind == token.KIND_KEYWORD {
		return nil, ErrNoMatch
	}
	c.Advance(1)

	var args []InstructionArgumentStatement
	for !c.Eof() {
		current := c.Current()
		if current.Literal == ";" {
			c.Advance(1)
			return InstructionStatement{
				Identifier: identifier.Literal,
				Arguments:  args,
			}, nil
		}
		if current.Literal == "," {
			c.Advance(1)
			continue
		}
		args = append(args, InstructionArgumentStatement{
			Type:  current.Kind,
			Value: current.Literal,
		})
		c.Advance(1)
	}

	return nil, errors.New("expected ';' at the end of instruction call")
}

// ParseModuleStatement parses the module token into the AST node.
func ParseModuleStatement(c Context, t *token.Token) (Node, error) {
	if t.Kind != token.KIND_KEYWORD && (t.Literal != "import" || t.Literal != "module") {
		return nil, ErrNoMatch
	}
	tType := t.Literal
	moduleType := ModuleStatementType("")
	switch tType {
	case "import":
		moduleType = ModuleStatementImporting
	case "module":
		moduleType = ModuleStatementDeclaration
	}

	c.Advance(1)

	identifier := c.Current()
	if identifier.Kind != token.KIND_IDENTIFIER {
		return nil, errors.New("expected a identifier after the module declaration")
	}

	return ModuleStatement{
		Identifier: identifier.Literal,
		Type:       moduleType,
	}, nil
}
