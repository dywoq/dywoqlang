package main

import (
	"fmt"

	"github.com/dywoq/dywoqlang/scanner"
)

func main() {
	input := "sds23d_SD 232323\nstring"
	s := scanner.New(input)
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
