package scanner

import "github.com/dywoq/dywoqlang/token"

type Reader interface {
	// Peek returns the future character.
	//
	// Returns an error if the scanner reached EOF,
	// the current position is out of the input.
	// or the current position+1 will make the scanner position out of the input.
	Peek() (rune, error)

	// Current returns the current character.
	//
	// Returns an error if the scanner reached EOF
	// or the current position is out of the input.
	Current() (rune, error)
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

// Context is an interface that allows you to work with scanners directly,
// such as reading character, getting position and create tokens more simple and comfortably.
type Context interface {
	Reader
	Tracker
	Creator
}

// Peek uses Reader.Peek method to peek the future character.
func Peek(r Reader) (rune, error) {
	return r.Peek()
}

// Current uses Reader.Current method to see the current character.
func Current(r Reader) (rune, error) {
	return r.Current()
}
