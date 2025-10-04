package parser

import (
	"strconv"

	"github.com/dywoq/dywoqlang/ast"
	"github.com/dywoq/dywoqlang/meta"
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
//   - meta expression: `meta(<literal, strings>)`. Function values and identifiers are not allowed.
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

		if next.Literal != "{" && (!declared || !linked) {
			return nil, c.Errorf("non-declared or non-linked functions must have a body")
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

	case t.Literal == "meta":
		_, _ = c.ExpectLiteral("meta")
		_, _ = c.ExpectLiteral("(")
		expr, err := ParseValue(c, false, false)
		if err != nil {
			return nil, err
		}
		var m meta.Type[any]
		switch t := expr.(type) {
		case ast.Value:
			if t.Kind == token.Identifier {
				return nil, c.Errorf("identifiers are not allowed in meta(...) expression")
			}
			num, err := strconv.Atoi(t.Value)
			if err != nil {
				return nil, err
			}
			m, err = meta.Integral(num)
			if err != nil {
				return nil, err
			}
		case ast.FunctionValue:
			return nil, c.Errorf("function values are not allowed in meta(...) expression")
		case ast.BinaryExpression:
			return nil, c.Errorf("binary expressions are not allowed in meta(...) expression")
		default:
			return nil, c.Errorf("unknown ast node: %v", ast.ToString(expr))
		}
		_, _ = c.ExpectLiteral(")")
		return ast.MetaExpression{
			Type:  m.Name,
			Value: expr,
			Data:  m,
		}, nil
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
	case token.String:
		return ParseModuleDeclaration(c)
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

// ParseModuleDeclaration parses a module declaration.
//
// Allowed syntax:
//    "main": {
//       # top declarations (functions, variables, constants...)
//    }
//
// Returns ast.ModuleDeclaration.
func ParseModuleDeclaration(c Context) (ast.Node, error) {
	ident, err := c.Expect(token.String)
	if err != nil {
		return nil, err
	}
	_, _ = c.ExpectLiteral(":")
	_, _ = c.ExpectLiteral("{")
	var body []ast.Node
	for {
		if c.Eof() {
			return nil, c.Errorf("module %s must be closed", ident.Literal)
		}
		if brace, _ := c.Current(); brace.Literal == "}" {
			break
		}
		n, err := ParseTopStatement(c)
		if err != nil {
			return nil, err
		}
		body = append(body, n)
	}
	_, _ = c.ExpectLiteral("}")

	c.SetModule(ident.Literal)
	return ast.ModuleDeclaration{
		Name: c.Module(),
		Body: body,
	}, nil
}

// ParseTopStatement parses the top statements.
// It can be a function, variable, constant or module.
func ParseTopStatement(c Context) (ast.Node, error) {
	t, err := c.Current()
	if err != nil {
		return nil, err
	}

	switch t.Kind {
	case token.String:
		return ParseModuleDeclaration(c)
	case token.Identifier, token.Keyword:
		return ParseDeclaration(c)
	default:
		return nil, c.Errorf("unexpected token at top level: %v", t.Literal)
	}
}
