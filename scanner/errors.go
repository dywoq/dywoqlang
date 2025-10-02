package scanner

import "errors"

var (
	// ErrNoMatch is returned by the tokenizers if the current character doesn't meet their requirements.
	// It's telling a scanner to try other tokenizer.
	ErrNoMatch = errors.New("no match")

	// ErrEof is returned by scanner if the scanner reached End Of File (EOF).
	ErrEof = errors.New("reached eof")
)
