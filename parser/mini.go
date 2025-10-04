package parser

import (
	"strconv"

	"github.com/dywoq/dywoqlang/ast"
	"github.com/dywoq/dywoqlang/token"
)

// MiniFunc is an alias for functions that represent mini parser.
type MiniFunc func(c Context) (ast.Node, error)

// ParseDeclaration parses a declaration statement.
//
// A declaration can include optional modifiers:
// - `link("module")` or `link(true|false)`
// - `export`
// - `declare`
//
// The parser expects the following structure:
//
//	[link(...)] [export] [declare] <identifier> <type> <value>
//
// Returns an *ast.Declaration node containing all metadata
// and parsed value (function, literal, or identifier).
func ParseDeclaration(c Context) (ast.Node, error) {
	var (
		exported, declared, linked bool
		canBeLinked                bool = true
		linkedFrom                 string
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
				return nil, c.Error("symbols that can't be linked can't be exported")
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
	value, err := ParseValue(c, declared, linked)
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

// ParseValue parses any value expression.
//
// A value can be:
//   - integer, float, or string literal
//   - identifier
//   - function value: `(x i32, y i32) { ... }`
//   - consteval expression: `consteval(<expr>)`
//
// Returns an *ast.Value or *ast.FunctionValue node.
//
// If the value is consteval, Consteval=true and
// the evaluated expression is stored in ValueNode.
func ParseValue(c Context, declared, linked bool) (ast.Node, error) {
	t, err := c.Current()
	if err != nil {
		return nil, err
	}

	switch {
	case t.Kind == token.Integer, t.Kind == token.Float, t.Kind == token.String:
		nodes, err := ParseBinaryExpression(c)
		if err != nil {
			return nil, err
		}
		if len(nodes) == 1 {
			return nodes[0], nil
		}
		return nodes[0], nil

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

			ident, err := c.Expect(token.Identifier)
			if err != nil {
				return nil, err
			}

			typ, err := c.Expect(token.Type)
			if err != nil {
				return nil, err
			}

			copyAllowed := true

			next, err = c.Current()
			if err != nil {
				return nil, err
			}

			if next.Literal == "copy" {
				_, _ = c.ExpectLiteral("copy")
				_, _ = c.ExpectLiteral("(")
				val, _ := c.Expect(token.BoolConstant)
				copyAllowed, _ = strconv.ParseBool(val.Literal)
				_, _ = c.ExpectLiteral(")")
			}

			params = append(params, ast.FunctionParameter{
				Identifier:  ident.Literal,
				Kind:        typ.Literal,
				CopyAllowed: copyAllowed,
			})

			next, _ = c.Current()
			if next.Literal == "," {
				_, _ = c.ExpectLiteral(",")
			}
		}
		_, _ = c.ExpectLiteral(")")
		next, err := c.Current()
		if err == nil && next.Literal == "{" {
			if declared || linked {
				return nil, c.Errorf("declared or linked functions cannot have a body")
			}
			body, err := ParseBody(c)
			if err != nil {
				return nil, err
			}
			return ast.FunctionValue{
				Parameters: params,
				Body:       body,
			}, nil
		}
		return ast.FunctionValue{
			Parameters: params,
			Body:       nil,
		}, nil

	case t.Literal == "consteval":
		_, _ = c.ExpectLiteral("consteval")
		_, _ = c.ExpectLiteral("(")
		expr, err := ParseValue(c, declared, linked)
		if err != nil {
			return nil, err
		}
		_, _ = c.ExpectLiteral(")")
		return ast.Value{
			Kind:      t.Kind,
			ValueNode: expr,
			Consteval: true,
		}, nil

	case t.Literal == "copy":
		_, _ = c.ExpectLiteral("copy")
		_, _ = c.ExpectLiteral("(")
		expr, err := ParseValue(c, false, false)
		if err != nil {
			return nil, err
		}
		_, _ = c.ExpectLiteral(")")
		return ast.Value{
			Kind:      t.Kind,
			ValueNode: expr,
			Consteval: false,
			Copied:    true,
		}, nil

	case token.BinaryOperatorsMap.Is(t.Literal):
		nodes, err := ParseBinaryExpression(c)
		if err != nil {
			return nil, err
		}
		if len(nodes) == 1 {
			return nodes[0], nil
		}
		return nodes[0], nil
	}

	return nil, c.Errorf("unknown value type: %v", t.Literal)
}

// ParseInstructionCall parses an instruction call.
//
// Instruction calls can be base or user-defined:
//   - Base: `mov x, 10;`
//   - User: `[ret] 10, 20;`
//
// Each argument inside the instruction is parsed via ParseValue.
// Returns an *ast.InstructionCall node containing arguments.
//
// If you want to provide a copied arguments:
//
//   - `[user_function] copy`
func ParseInstructionCall(c Context) (ast.Node, error) {
	var (
		isUser bool
		name   string
	)

	t, err := c.Current()
	if err != nil {
		return nil, err
	}
	switch t.Kind {
	case token.Separator:
		_, _ = c.ExpectLiteral("[")
		ident, _ := c.Expect(token.Identifier)
		_, _ = c.ExpectLiteral("]")
		isUser = true
		name = ident.Literal
	case token.BaseInstruction:
		ident, _ := c.Expect(token.BaseInstruction)
		name = ident.Literal
	}

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

		val, err := ParseValue(c, false, false)
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

// ParseStatement parses a single statement.
//
// A statement currently can only be:
//   - Base instruction (e.g. `mov ...;`)
//   - User instruction (e.g. `[ret] ...;`)
//
// Returns an *ast.InstructionCall node.
// Any unexpected token produces an error.
func ParseStatement(c Context) (ast.Node, error) {
	t, _ := c.Current()

	switch t.Kind {
	case token.BaseInstruction:
		return ParseInstructionCall(c)
	case token.Separator:
		return ParseInstructionCall(c)
	default:
		return nil, c.Errorf("unexpected token in statement: %v", t.Literal)
	}
}

// ParseBody parses a function or block body.
//
// The body must start with '{' and end with '}'.
// Inside the braces, statements are parsed by ParseStatement.
//
// Returns a slice of AST nodes representing statements.
func ParseBody(c Context) ([]ast.Node, error) {
	_, _ = c.ExpectLiteral("{")
	var statements []ast.Node

	for !c.Eof() {
		t, _ := c.Current()
		if t.Literal == "}" {
			break
		}

		stmt, err := ParseStatement(c)
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	_, _ = c.ExpectLiteral("}")
	return statements, nil
}

func ParseBinaryExpression(c Context) ([]ast.Node, error) {
	var nodes []ast.Node

	parseOperand := func() (ast.Node, error) {
		t, err := c.Current()
		if err != nil {
			return nil, err
		}
		switch t.Kind {
		case token.Integer, token.Float, token.Identifier:
			_ = c.Advance(1)
			return ast.Value{Value: t.Literal, Kind: t.Kind}, nil
		default:
			return nil, c.Errorf("expected number or identifier, got %q", t.Literal)
		}
	}

	parseMulDiv := func() (ast.Node, error) {
		left, err := parseOperand()
		if err != nil {
			return nil, err
		}

		for !c.Eof() {
			t, _ := c.Current()
			if t.Kind != token.BinaryOperator || (t.Literal != "*" && t.Literal != "/") {
				break
			}
			op := t.Literal[0]
			_ = c.Advance(1)

			right, err := parseOperand()
			if err != nil {
				return nil, err
			}

			left = ast.BinaryExpression{
				Operator: string(op),
				Children: []ast.Node{left, right},
			}
		}
		return left, nil
	}

	for !c.Eof() {
		left, err := parseMulDiv()
		if err != nil {
			return nil, err
		}

		for !c.Eof() {
			t, _ := c.Current()
			if t.Kind != token.BinaryOperator || (t.Literal != "+" && t.Literal != "-") {
				break
			}
			op := t.Literal[0]
			_ = c.Advance(1)

			right, err := parseMulDiv()
			if err != nil {
				return nil, err
			}

			left = ast.BinaryExpression{
				Operator: string(op),
				Children: []ast.Node{left, right},
			}
		}

		nodes = append(nodes, left)
		break
	}

	return nodes, nil
}
