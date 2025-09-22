package main

import (
	"fmt"
	"os"

	// "github.com/dywoq/dywoqlang/parser"
	"github.com/dywoq/dywoqlang/scanner"
)

func main() {
	data, err := os.ReadFile("./main.dl")
	if err != nil {
		panic(err)
	}

	s := scanner.New(string(data))
	tokens, err := s.Scan()
	if err != nil {
		panic(err)
	}

	for _, token := range tokens {
		fmt.Println(*token)
	}

	// p := parser.New(tokens)
	// ast, err := p.Parse()
	// if err != nil {
	// 	panic(err)
	// }

	// json, err := parser.NodeToJson(ast)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(json))
}
