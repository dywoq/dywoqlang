package parser

import (
	"github.com/dywoq/dywoqlang/token"
)

type Reader interface {
	// Current returns the current token.
	//
	// Returns a error if the parser reached End of File (EOF) token.
	Current() (*token.Token, error)

	// Peek returns the future token.
	//
	// Returns a error if the parser reached End Of File (EOF) token,
	// or the current position+1 will make the position out of bounds.
	Peek() (*token.Token, error)
}

type Tracker interface {
	// Position returns the current position.
	Position() int
}

type EofChecker interface {
	// Eof reports whether parser meet the End Of File (EOF) token.
	Eof() bool
}

type Advancer interface {
	// Advance advances to the next position by n.
	//
	// Returns an error if the parser reached End Of File (EOF) token,
	// or the current position+n will make the position out of bounds of the tokens.
	Advance(n int) error
}

type ErrorCreator interface {
	// Error returns a new error.
	// The difference from errors.New is that Error automatically inserts the position where the error occurred.
	Error(v ...any) error

	// Errorf returns a new error,
	// but formatted.
	// The difference from fmt.Errorf is that Error automatically inserts the position where the error occurred.
	Errorf(format string, v ...any) error
}

type Expecter interface {
	Expect(kind token.Kind) (*token.Token, error)
	ExpectLiteral(lit string) (*token.Token, error)
	ExpectMultiple(kind ...token.Kind) (*token.Token, error)
	ExpectLiterals(lits ...string) (*token.Token, error)
}

type Context interface {
	Reader
	Tracker
	EofChecker
	Advancer
	ErrorCreator
	Expecter
}
