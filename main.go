package main

import (
	"fmt"

	"github.com/dywoq/dywoqlang/scanner"
	"github.com/dywoq/dywoqlang/token"
)

func main() {
	s := scanner.New(true)

	tokens, err := s.Scan("mov e, 2+2;\nstdout eax;")
	if err != nil {
		panic(err)
	}

	for _, tok := range tokens {
		fmt.Println(token.ToString(tok))
	}
}
