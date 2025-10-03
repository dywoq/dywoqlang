package parser

import (
	"errors"
	"fmt"
	"github.com/dywoq/dywoqlang/ast"
	"log"
	"slices"

	"github.com/dywoq/dywoqlang/token"
)

type Parser struct {
	pos     int
	parsers []MiniFunc
	tokens  []*token.Token

	debug   bool
	setupOn bool
}

func New(debug bool) *Parser {
	return &Parser{
		pos:     0,
		parsers: make([]MiniFunc, 0),
		tokens:  make([]*token.Token, 0),
		debug:   debug,
		setupOn: false,
	}
}

func (p *Parser) Current() (*token.Token, error) {
	if p.Eof() {
		return nil, ErrEof
	}
	p.outputf("getting current token: %v\n", p.tokens[p.pos])
	return p.tokens[p.pos], nil
}

func (p *Parser) Peek() (*token.Token, error) {
	switch {
	case p.Eof():
		return nil, ErrEof
	case p.pos+1 >= len(p.tokens):
		return nil, errors.New("the current position+1 will make the position out of tokens")
	}
	return p.tokens[p.pos+1], nil
}

func (p *Parser) Advance(n int) error {
	switch {
	case p.Eof():
		return nil
	case p.pos+n >= len(p.tokens):
		return fmt.Errorf("the current position+%d will make the position out of tokens", n)
	}
	p.pos += n
	return nil
}

func (p *Parser) Eof() bool {
	return p.tokens[p.pos].Kind == token.Eof || p.pos > len(p.tokens)
}

func (p *Parser) Position() int {
	return p.pos
}

func (p *Parser) Expect(kind token.Kind) (*token.Token, error) {
	if p.Eof() {
		return nil, ErrEof
	}
	tok := p.tokens[p.pos]
	if tok.Kind != kind {
		return nil, p.Errorf("expected token kind %v, got %v", kind, tok.Kind)
	}
	p.pos++
	return tok, nil
}

func (p *Parser) ExpectLiteral(lit string) (*token.Token, error) {
	if p.Eof() {
		return nil, ErrEof
	}
	tok := p.tokens[p.pos]
	if tok.Literal != lit {
		return nil, p.Errorf("expected literal '%s', got '%s'", lit, tok.Literal)
	}
	p.pos++
	return tok, nil
}

func (p *Parser) ExpectMultiple(kinds ...token.Kind) (*token.Token, error) {
	if p.Eof() {
		return nil, ErrEof
	}
	tok := p.tokens[p.pos]
	if !slices.Contains(kinds, tok.Kind) {
		return nil, p.Errorf("expected at least one token kind (%v), got %v", kinds, tok.Kind)
	}
	p.pos++
	return tok, nil
}

func (p *Parser) Error(v ...any) error {
	t, _ := p.Current()
	pos := &token.Position{}
	if t != nil {
		pos = t.Position
	}
	formatted := fmt.Sprintf("%v (source is %d, token position: %v)", v, p.pos, pos)
	return errors.New(formatted)
}

func (p *Parser) Errorf(format string, v ...any) error {
	t, _ := p.Current()
	pos := &token.Position{}
	if t != nil {
		pos = t.Position
	}
	custom := fmt.Sprintf(format, v...)
	formatted := fmt.Sprintf("%s (source is %d, token position: %v)", custom, p.pos, pos)
	return errors.New(formatted)
}

func (p *Parser) outputf(format string, v ...any) {
	if p.debug {
		log.Printf(format, v...)
	}
}

func (p *Parser) Parse(tokens []*token.Token) ([]ast.Node, error) {
	if len(tokens) == 0 {
		return nil, errors.New("tokens slice is empty")
	}
	p.setup()
	if len(p.parsers) == 0 {
		return nil, errors.New("there are no mini parsers")
	}
	p.reset(tokens)

	nodes := []ast.Node{}
	for !p.Eof() {
		node, err := p.parse()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
		p.outputf("parsed %s\n", ast.ToString(node))
	}

	return nodes, nil
}

func (p *Parser) parse() (ast.Node, error) {
	if len(p.parsers) == 0 {
		return nil, errors.New("there are no mini parsers")
	}
	for _, parser := range p.parsers {
		node, err := parser(p)
		if err != nil {
			if errors.Is(err, ErrEof) {
				continue
			}
			return nil, err
		}
		return node, nil
	}
	t, _ := p.Current()
	return nil, p.Errorf("met illegal token: %s", token.ToString(t))
}

func (p *Parser) setup() {
	if !p.setupOn {
		p.parsers = []MiniFunc{
			ParseDeclaration,
		}
		p.setupOn = true
	}
}

func (p *Parser) reset(tokens []*token.Token) {
	p.tokens = tokens
	p.pos = 0
}
