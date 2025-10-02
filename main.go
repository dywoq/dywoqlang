package main

import (
	"fmt"

	"github.com/dywoq/dywoqlang/scanner"
	"github.com/dywoq/dywoqlang/token"
)

func main() {
	s := scanner.New()

	tokens, err := s.Scan("122")
	if err != nil {
		panic(err)
	}

	for _, tok := range tokens {
		fmt.Println(token.ToString(tok))
	}
}
