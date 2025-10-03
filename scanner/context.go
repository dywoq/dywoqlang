package scanner

import "github.com/dywoq/dywoqlang/token"

type Reader interface {
	// Peek returns the future character.
	//
	// Returns an error if the scanner reached EOF,
	// or the current position+1 will make the scanner position out of the input.
	Peek() (rune, error)

	// Current returns the current character.
	//
	// Returns an error if the scanner reached EOF
	// or the current position is out of the input.
	Current() (rune, error)

	// Input returns the current input from the scanner.
	Input() string
}

type Tracker interface {
	// Position returns a current position.
	Position() *token.Position
}

type Creator interface {
	// New returns a new token.
	//
	// The only difference from token.NewToken is
	// that scanner automatically inserts the position.
	New(literal string, kind token.Kind) *token.Token
}

type EofChecker interface {
	// Eof returns true if scanner reached End Of File (EOF).
	Eof() bool
}

type Slicer interface {
	// Slice takes a substring from the input surrounded by start and end.
	//
	// Returns an error if start is negative, start is higher than end,
	// or end is higher than the input.
	Slice(start, end int) (string, error)
}

type Advancer interface {
	// Advance advances to the next position by n
	//
	// If newline character is met,
	// scanner increases line and column beside the current position.
	//
	// Returns an error if scanner reached End Of File (EOF),
	// or the current position+n will make the position out of the input.
	//
	// Does nothing if n is zero.
	Advance(n int) error
}

// Context is an interface that allows you to work with scanners directly,
// such as reading character, getting position and create tokens more simple and comfortably.
type Context interface {
	Reader
	Tracker
	Creator
	EofChecker
	Slicer
	Advancer
}

// Peek uses Reader.Peek method to peek the future character.
func Peek(r Reader) (rune, error) {
	return r.Peek()
}

// Current uses Reader.Current method to see the current character.
func Current(r Reader) (rune, error) {
	return r.Current()
}

// Eof uses EofChecker.Eof method to check if scanner reached EOF.
func Eof(e EofChecker) bool {
	return e.Eof()
}

// Slice uses Slicer.Slice method to slice the scanner input.
func Slice(s Slicer, start, end int) (string, error) {
	return s.Slice(start, end)
}

// Advance uses Advancer.Advance method to advance to the next position by n.
func Advance(a Advancer, n int) error {
	return a.Advance(n)
}
