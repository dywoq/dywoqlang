package main

import (
	"fmt"
	"os"

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
		if token != nil {
			fmt.Println(*token)
		}
	}
}
