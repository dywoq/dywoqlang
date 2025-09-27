package main

import (
	"fmt"
	"os"

	"github.com/dywoq/dywoqlang/parser"
	"github.com/dywoq/dywoqlang/scanner"
	"github.com/dywoq/dywoqlang/token"
)

func main() {
	content, err := os.ReadFile("main.dl")
	if err != nil {
		panic(err)
	}

	s := scanner.NewScanner()
	tokens, err := s.Scan(string(content))
	if err != nil {
		panic(err)
	}

	for _, tok := range tokens {
		fmt.Printf("tok: %v\n", token.String(tok))
	}

	p := parser.NewParser()
	tree, err := p.Parse(tokens)
	if err != nil {
		panic(err)
	}

	for _, node := range tree {
		fmt.Printf("parser.NodeString(node): %v\n", parser.NodeString(node))
	}
}
