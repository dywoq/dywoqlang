package token_test

import (
	"testing"

	"github.com/dywoq/dywoqlang/token"
)

func TestIsArithmeticOperator(t *testing.T) {
	tests := []struct {
		input rune
		want  bool
	}{
		{'+', true},
		{'-', true},
		{'/', true},
		{'*', true},
		{'#', false},
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			got := token.IsArithmeticOperator(test.input)
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestIsIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"y", true},
		{"x_y", true},
		{"x_y2", true},
		{"export", false},
		{"int32", false},
		{",", false},
		{"", false},
		{"mov", false},
		{"msd#", false},
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			got := token.IsIdentifier(test.input)
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		name   string
		m      token.Map
		inputs map[string]bool
	}{
		{"types", token.Types, map[string]bool{
			"int8":   true,
			"int16":  true,
			"int32":  true,
			"int64":  true,
			"uint8":  true,
			"uint16": true,
			"uint32": true,
			"uint64": true,
			"string": true,
			"bool":   true,
			"export": false,
			"mov":    false,
		}},
		{"keywords", token.Keywords, map[string]bool{
			"export": true,
			"int8":   false,
			"int32":  false,
			"int64":  false,
		}},
		{"separators", token.Separators, map[string]bool{
			",": true,
			"{": true,
			"}": true,
			"(": true,
			")": true,
			"#": false,
			"/": false,
		}},
		{"bool constants", token.BoolConstants, map[string]bool{
			"true":  true,
			"false": true,
			"1":     false,
			"0":     false,
			"TRUE":  false,
			"FALSE": false,
			"T":     false,
			"F":     false,
		}},
		{"instructions", token.Instructions, map[string]bool{
			"add":      true,
			"minus":    true,
			"divide":   true,
			"multiply": true,
			"mov":      true,
			"make":     true,
			"write":    false,
			"output":   false,
			"format":   false,
			"export":   false,
			"int8":     false,
			"string":   false,
			"bool":     false,
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for input, want := range test.inputs {
				got := test.m.Is(input)
				if got != want {
					t.Errorf("%s: got %v, want %v", input, got, want)
				}
			}
		})
	}
}
