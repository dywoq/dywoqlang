package parser

import (
	"encoding/json"

	"github.com/dywoq/dywoqlang/token"
)

type Node interface {
	Node()
}

type FunctionDeclaration struct {
	Name         string   `json:"name"`
	ParamsTypes  []string `json:"params_types"`
	ReturnType   string   `json:"return_type"`
	Body         []Node   `json:"body"`
	Exported     bool     `json:"exported"`
	ExportedFrom string   `json:"exported_from"`
}

type InstructionCall struct {
	Name      string `json:"name"`
	Arguments []Node `json:"arguments"`
}

type InstructionCallArgument struct {
	Value string `json:"value"`
}

type VariableDeclaration struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Value        Node   `json:"value"`
	Exported     bool   `json:"exported"`
	ExportedFrom string `json:"exported_from"`
}

type Illegal struct{}

type Program struct {
	Statements []Node `json:"statements"`
}

type ValueNode struct {
	Kind  token.Kind `json:"kind"`
	Value string     `json:"value"`
}

type ModuleDeclaration struct {
	Name string `json:"name"`
}

func (FunctionDeclaration) Node()     {}
func (InstructionCall) Node()         {}
func (VariableDeclaration) Node()     {}
func (Illegal) Node()                 {}
func (Program) Node()                 {}
func (InstructionCallArgument) Node() {}
func (ValueNode) Node()               {}
func (ModuleDeclaration) Node()       {}

// NodeToJson converts n to JSON presentation.
func NodeToJson(n Node) ([]byte, error) {
	return json.MarshalIndent(n, "ast expression: ", "  ")
}
