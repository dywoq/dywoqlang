package main

import (
	"fmt"
	"os"

	"github.com/dywoq/dywoqlang/ast"
	"github.com/dywoq/dywoqlang/parser"
	"github.com/dywoq/dywoqlang/scanner"
	"github.com/dywoq/dywoqlang/token"
)

func main() {
	bytes, err := os.ReadFile("main.dl")
	if err != nil {
		panic(err)
	}

	s := scanner.New(true)

	tokens, err := s.Scan(string(bytes))
	if err != nil {
		panic(err)
	}

	for _, tok := range tokens {
		fmt.Println(token.ToString(tok))
	}

	p := parser.New(true)

	nodes, err := p.Parse(tokens)
	if err != nil {
		panic(err)
	}

	for _, node := range nodes {
		fmt.Println(ast.ToString(node))
	}
}
