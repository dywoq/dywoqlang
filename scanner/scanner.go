package scanner

import "github.com/dywoq/dywoqlang/token"

// tokenizer is a function type, representing a tokenizer.
type tokenizer func() (*token.Token, error)

type Scanner struct {
	input      string
	tokenizers []tokenizer
}
