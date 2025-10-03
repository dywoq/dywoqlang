package parser

import (
	"github.com/dywoq/dywoqlang/ast"
	"github.com/dywoq/dywoqlang/token"
)

// MiniFunc is an alias for functions that represent mini parser.
type MiniFunc func(c Context) (ast.Node, error)

func ParseDeclaration(c Context) (ast.Node, error) {
	if c.Eof() {
		return nil, ErrEof
	}

	var (
		exported, declared, linked bool
		linkedFrom                 string
	)

loop:
	for {
		t, err := c.Current()
		if err != nil {
			return nil, err
		}

		switch t.Literal {
		case "link":
			c.Advance(1)
			if t, _ = c.Current(); t.Literal != "(" {
				return nil, c.Error("expected left paren")
			}

			c.Advance(1)
			if t, _ = c.Current(); t.Kind != token.String {
				return nil, c.Error("expected string")
			} else {
				linkedFrom = t.Literal
			}

			c.Advance(1)
			if t, _ = c.Current(); t.Literal != ")" {
				return nil, c.Error("expected right paren")
			}

			c.Advance(1)

			linked = true

		case "export":
			exported = true
			c.Advance(1)

		case "declare":
			declared = true
			c.Advance(1)

		default:
			break loop
		}
	}

	identifier, _ := c.Current()
	if identifier.Kind != token.Identifier {
		return nil, c.Error("expected an identifier")
	}
	c.Advance(1)

	tType, _ := c.Current()
	if tType.Kind != token.Type {
		return nil, c.Error("expected a type")
	}
	c.Advance(1)

	value, err := ParseValue(c, declared)
	if err != nil {
		return nil, err
	}

	return ast.Declaration{
		Name:       identifier.Literal,
		Kind:       tType.Kind,
		Exported:   exported,
		Declared:   declared,
		Linked:     linked,
		LinkedFrom: linkedFrom,
		Value:      value,
	}, nil
}

func ParseValue(c Context, declared bool) (ast.Node, error) {
	t, _ := c.Current()
	switch {
	case t.Kind == token.Integer,
		t.Kind == token.Float,
		t.Kind == token.String:
		return ast.Value{Value: t.Literal, Kind: t.Kind}, nil
	case t.Literal == "(":
		c.Advance(1)

		params := []ast.FunctionParameter{}
		for {
			if c.Eof() {
				break
			}
			if r, _ := c.Current(); r.Literal == ")" {
				break
			}
			identifier, _ := c.Current()
			if identifier.Kind != token.Identifier {
				return nil, c.Error("expected identifier in function parameter")
			}
			c.Advance(1)

			tType, _ := c.Current()
			if tType.Kind != token.Type {
				return nil, c.Error("expected type in function parameter")
			}
			c.Advance(1)

			params = append(params, ast.FunctionParameter{Identifier: identifier.Literal, Kind: tType.Kind})

			if r, _ := c.Current(); r.Literal == "," {
				continue
			}
		}

		return ast.FunctionValue{Parameters: params}, nil
	}

	return nil, c.Error("unknown value type")
}
