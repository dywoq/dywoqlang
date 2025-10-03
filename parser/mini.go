package parser

import (
	"strconv"

	"github.com/dywoq/dywoqlang/ast"
	"github.com/dywoq/dywoqlang/token"
)

// MiniFunc is an alias for functions that represent mini parser.
type MiniFunc func(c Context) (ast.Node, error)

func ParseDeclaration(c Context) (ast.Node, error) {
	var (
		exported, declared, linked, canBeLinked bool
		linkedFrom                              string
	)
loop:
	for !c.Eof() {
		t, _ := c.Current()
		switch t.Literal {
		case "link":
			_, _ = c.ExpectLiteral("link")
			_, _ = c.ExpectLiteral("(")
			v, err := c.ExpectMultiple(token.String, token.BoolConstant)
			if err != nil {
				return nil, err
			}

			if v.Kind == token.BoolConstant {
				canBeLinked, _ = strconv.ParseBool(v.Literal)
			}

			if v.Kind == token.String {
				linkedFrom = v.Literal
				linked = true
			}

			_, _ = c.ExpectLiteral(")")

		case "export":
			if !canBeLinked {
				return nil, c.Error("symbols that can be linked can't be exported")
			}
			_, _ = c.ExpectLiteral("export")
			exported = true

		case "declare":
			_, _ = c.ExpectLiteral("declare")
			declared = true

		default:
			break loop
		}
	}

	identifier, _ := c.Expect(token.Identifier)
	tType, _ := c.Expect(token.Type)
	value, err := ParseValue(c, declared)
	if err != nil {
		return nil, err
	}

	return ast.Declaration{
		Name:        identifier.Literal,
		Kind:        tType.Literal,
		Exported:    exported,
		Declared:    declared,
		Linked:      linked,
		LinkedFrom:  linkedFrom,
		CanBeLinked: canBeLinked,
		Value:       value,
	}, nil
}

func ParseValue(c Context, declared bool) (ast.Node, error) {
	t, err := c.Current()
	if err != nil {
		return nil, err
	}

	switch {
	case t.Kind == token.Integer, t.Kind == token.Float, t.Kind == token.String:
		_, _ = c.Expect(t.Kind)
		return ast.Value{Value: t.Literal, Kind: t.Kind}, nil

	case t.Kind == token.Identifier:
		_, _ = c.Expect(token.Identifier)
		return ast.Value{Value: t.Literal, Kind: t.Kind}, nil

	case t.Literal == "(":
		_, _ = c.ExpectLiteral("(")
		params := []ast.FunctionParameter{}

		for !c.Eof() {
			next, _ := c.Current()
			if next.Literal == ")" {
				break
			}

			ident, _ := c.Expect(token.Identifier)
			typ, _ := c.Expect(token.Type)

			params = append(params, ast.FunctionParameter{
				Identifier: ident.Literal,
				Kind:       typ.Kind,
			})

			next, _ = c.Current()
			if next == nil {
				return nil, c.Error("next character is nil")
			}

			if next.Literal == "," {
				_, _ = c.ExpectLiteral(",")
			}
		}

		_, _ = c.ExpectLiteral(")")
		return ast.FunctionValue{Parameters: params}, nil
	}

	return nil, c.Errorf("unknown value type: %v", t.Literal)
}

func ParseInstructionCall(c Context) (ast.Node, error) {
	var (
		isUser bool
		name   string
	)

	t, _ := c.Current()
	switch t.Kind {
	case token.Separator: // "["
		_, _ = c.ExpectLiteral("[")
		ident, _ := c.Expect(token.Identifier)
		_, _ = c.ExpectLiteral("]")
		isUser = true
		name = ident.Literal
	case token.BaseInstruction:
		ident, _ := c.Expect(token.BaseInstruction)
		name = ident.Literal
	}

	// parse arguments normally
	var args []ast.InstructionCallArgument
	for {
		if c.Eof() {
			break
		}
		nextToken, _ := c.Current()
		if nextToken.Literal == ";" {
			c.Advance(1)
			break
		}
		if nextToken.Literal == "," {
			c.Advance(1)
			continue
		}

		val, err := ParseValue(c, false) 
		if err != nil {
			return nil, err
		}
		args = append(args, ast.InstructionCallArgument{
			Value: val,
			Kind:  nextToken.Kind,
		})
	}

	return ast.InstructionCall{
		Name:      name,
		IsUser:    isUser,
		Arguments: args,
	}, nil
}
