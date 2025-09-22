package parser

import (
	"encoding/json"

	"github.com/dywoq/dywoqlang/token"
)

type ModuleType string

type Node interface {
	Node()
}

type FunctionDeclaration struct {
	Name        string   `json:"name"`
	ParamsTypes []string `json:"params_types"`
	ReturnType  string   `json:"return_type"`
	Body        []Node   `json:"body"`
	Exported    bool     `json:"exported"`
	DeclaredIn  string   `json:"declared_in"`
	Declared    bool     `json:"declared"`
}

type InstructionCall struct {
	Name      string `json:"name"`
	Arguments []Node `json:"arguments"`
}

type InstructionCallArgument struct {
	Value string `json:"value"`
}

type VariableDeclaration struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Value      Node   `json:"value"`
	Exported   bool   `json:"exported"`
	DeclaredIn string `json:"declared_in"`
	Declared   bool   `json:"declared"`
}

type Illegal struct{}

type Program struct {
	Statements []Node `json:"statements"`
}

type ValueNode struct {
	Kind  token.Kind `json:"kind"`
	Value string     `json:"value"`
}

type Module struct {
	Name string     `json:"name"`
	Type ModuleType `json:"type"`
}

const (
	ModuleDeclaration ModuleType = "declaration"
	ModuleImporting   ModuleType = "importing"
)

func (FunctionDeclaration) Node()     {}
func (InstructionCall) Node()         {}
func (VariableDeclaration) Node()     {}
func (Illegal) Node()                 {}
func (Program) Node()                 {}
func (InstructionCallArgument) Node() {}
func (ValueNode) Node()               {}
func (Module) Node()                  {}

// NodeToJson converts n to JSON presentation.
func NodeToJson(n Node) ([]byte, error) {
	return json.MarshalIndent(n, "", "  ")
}
