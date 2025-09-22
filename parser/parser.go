package parser

import (
	"errors"
	"fmt"

	"github.com/dywoq/dywoqlang/token"
)

type Parser struct {
	pos           int
	tokens        []*token.Token
	parsers       []parserFunc
	setupOn       bool
	currentModule string
}

func New(tokens []*token.Token) *Parser {
	return &Parser{0, tokens, []parserFunc{}, false, "main"}
}

type parserFunc func(t *token.Token) (Node, error)

var (
	errNoMatch = errors.New("no match")
)

func (p *Parser) advance(n int) {
	if p.pos+n >= len(p.tokens) {
		p.pos = len(p.tokens)
		return
	}
	p.pos += n
}

func (p *Parser) current() *token.Token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.pos]
}

func (p *Parser) reset() {
	p.pos = 0
}

func (p *Parser) setup() {
	p.setupOn = true
	p.parsers = append(p.parsers, p.parseInstructionCall)
	p.parsers = append(p.parsers, p.parseDeclaration)
	p.parsers = append(p.parsers, p.parseModule)
}

func (p *Parser) parse() (Node, error) {
	if len(p.parsers) == 0 {
		return nil, errors.New("there are no parsers")
	}
	for _, parser := range p.parsers {
		n, err := parser(p.current())
		if err != nil {
			if err == errNoMatch {
				continue
			}
			return nil, err
		}
		return n, nil
	}
	return Illegal{}, nil
}

func (p *Parser) isFunctionStart() bool {
	if p.current() == nil {
		return false
	}
	if p.current().Literal == "export" {
		return true
	}
	if p.pos+1 < len(p.tokens) &&
		p.current().Kind == token.KIND_IDENTIFIER &&
		p.tokens[p.pos+1].Kind == token.KIND_TYPE {
		return true
	}
	return false
}

func (p *Parser) parseValue() (Node, error) {
	t := p.current()
	if t == nil {
		return ValueNode{}, errors.New("unexpected EOF")
	}

	switch t.Kind {
	case token.KIND_INTEGER, token.KIND_IDENTIFIER, token.KIND_BOOL_CONSTANT:
		node := ValueNode{Kind: t.Kind, Value: t.Literal}
		p.advance(1)
		return node, nil
	case token.KIND_SEPARATOR:
		if t.Literal == "(" {
			p.advance(1)
			nodes := []Node{}
			for p.current() != nil && p.current().Literal != ")" {
				v, err := p.parseValue()
				if err != nil {
					return ValueNode{}, err
				}
				nodes = append(nodes, v)
				if p.current() != nil && p.current().Literal == "," {
					p.advance(1)
				}
			}
			if p.current() == nil || p.current().Literal != ")" {
				return ValueNode{}, errors.New("expected closing )")
			}
			p.advance(1)
			return ValueNode{Kind: token.KIND_SEPARATOR, Value: "composite"}, nil
		}
	default:
		return ValueNode{}, fmt.Errorf("unexpected token kind: %s", t.Literal)
	}
	return ValueNode{}, errNoMatch
}

func (p *Parser) parseInstructionCall(t *token.Token) (Node, error) {
	if t == nil || t.Kind != token.KIND_BASE_INSTRUCTION {
		return nil, errNoMatch
	}

	instructionName := t.Literal
	p.advance(1)

	args := []Node{}
	for p.current() != nil &&
		p.current().Kind != token.KIND_EOF &&
		p.current().Kind != token.KIND_BASE_INSTRUCTION &&
		!p.isFunctionStart() {

		if p.current().Kind == token.KIND_SEPARATOR && p.current().Literal == ";" {
			p.advance(1)
			continue
		}

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		args = append(args, value)

		if p.current() != nil && p.current().Kind == token.KIND_SEPARATOR && p.current().Literal == "," {
			p.advance(1)
		}
	}

	return InstructionCall{
		Name:      instructionName,
		Arguments: args,
	}, nil
}

func (p *Parser) parseDeclaration(t *token.Token) (Node, error) {
	startPos := p.pos
	exported := false
	if t != nil && t.Literal == "export" {
		exported = true
		p.advance(1)
	}

	identifier := p.current()
	if identifier == nil || identifier.Kind != token.KIND_IDENTIFIER {
		p.pos = startPos
		return nil, errNoMatch
	}
	p.advance(1)

	ttype := p.current()
	if ttype == nil || ttype.Kind != token.KIND_TYPE {
		p.pos = startPos
		return nil, errNoMatch
	}
	p.advance(1)

	if p.current() != nil && p.current().Literal == "(" {
		return p.parseFunctionDeclaration(identifier, ttype, exported)
	}

	valueNode, err := p.parseValue()
	if err != nil {
		p.pos = startPos
		return nil, err
	}
	return VariableDeclaration{
		Name:         identifier.Literal,
		Type:         ttype.Literal,
		Value:        valueNode,
		Exported:     exported,
		DeclaredIn:   p.currentModule,
	}, nil
}

func (p *Parser) parseFunctionDeclaration(identifier, ttype *token.Token, exported bool) (Node, error) {
	startPos := p.pos

	if p.current() == nil || p.current().Literal != "(" {
		p.pos = startPos
		return nil, errNoMatch
	}
	p.advance(1)

	params := []string{}
	for p.current() != nil && p.current().Literal != ")" {
		tok := p.current()
		if tok.Kind != token.KIND_TYPE {
			return nil, fmt.Errorf("expected type in function declaration of %s", identifier.Literal)
		}
		params = append(params, tok.Literal)
		p.advance(1)
		if p.current() != nil && p.current().Literal == "," {
			p.advance(1)
		}
	}

	if p.current() == nil || p.current().Literal != ")" {
		p.pos = startPos
		return nil, errNoMatch
	}
	p.advance(1)

	if p.current() == nil || p.current().Literal != ":" {
		p.pos = startPos
		return nil, errNoMatch
	}
	p.advance(1)

	body := []Node{}
	for p.current() != nil && p.current().Kind != token.KIND_EOF && !p.isFunctionStart() {
		if p.current().Kind == token.KIND_SEPARATOR && p.current().Literal == ";" {
			p.advance(1)
			continue
		}

		if p.current().Kind == token.KIND_BASE_INSTRUCTION {
			instr, err := p.parseInstructionCall(p.current())
			if err != nil {
				return nil, err
			}
			body = append(body, instr)
			continue
		}

		p.advance(1)
	}
	return FunctionDeclaration{
		Name:         identifier.Literal,
		ParamsTypes:  params,
		ReturnType:   ttype.Literal,
		Body:         body,
		Exported:     exported,
		DeclaredIn:   p.currentModule,
	}, nil
}

func (p *Parser) parseModule(t *token.Token) (Node, error) {
	if t.Literal == "module" {
		p.advance(1)
		moduleName := p.current()
		if moduleName.Kind != token.KIND_IDENTIFIER {
			return nil, fmt.Errorf("expected module name at line %d, column %d: %s", moduleName.Position.Line, moduleName.Position.Column, moduleName.Literal)
		}
		p.currentModule = moduleName.Literal
		p.advance(1)
		return Module{Name: moduleName.Literal, Type: ModuleDeclaration}, nil
	}

	if t.Literal == "import" {
		p.advance(1)
		moduleName := p.current()
		if moduleName.Kind != token.KIND_IDENTIFIER {
			return nil, fmt.Errorf("expected module name after import at line %d, column %d: %s", moduleName.Position.Line, moduleName.Position.Column, moduleName.Literal)
		}
		p.advance(1)
		return Module{Name: moduleName.Literal, Type: ModuleImporting}, nil
	}

	return nil, errNoMatch
}

func (p *Parser) Parse() (Node, error) {
	p.reset()
	if !p.setupOn {
		p.setup()
	}
	program := Program{}
	for p.pos < len(p.tokens) && p.current().Kind != token.KIND_EOF {
		node, err := p.parse()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, node)
	}
	return program, nil
}
