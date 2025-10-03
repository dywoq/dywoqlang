package parser

import "errors"

var (
	// ErrNoMatch error returned is by mini parsers,
	// meaning the error tells the parser it should try other mini parser.
	ErrNoMatch = errors.New("no match")

	// ErrEof error is usually returned by mini parsers,
	// meaning the parser reached EOF token.
	ErrEof = errors.New("reached eof token")
)
