package ast

import (
	"encoding/json"

	"github.com/dywoq/dywoqlang/token"
)

type Node interface {
	Node()
}

type Declaration struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Exported    bool   `json:"exported"`
	Declared    bool   `json:"declared"`
	Linked      bool   `json:"linked"`
	LinkedFrom  string `json:"linked_from"`
	CanBeLinked bool   `json:"can_be_linked"`
	Source      string `json:"source"`
	Value       Node   `json:"value"`
}

type FunctionParameter struct {
	Identifier  string `json:"identifier"`
	CopyAllowed bool   `json:"copy_allowed"`
	Kind        string `json:"kind"`
}

type FunctionValue struct {
	Body       []Node              `json:"body"`
	Parameters []FunctionParameter `json:"parameters"`
}

type Value struct {
	Consteval bool       `json:"consteval"`
	Copied    bool       `json:"copied"`
	Kind      token.Kind `json:"kind"`
	Value     string     `json:"value,omitempty"`
	ValueNode Node       `json:"value_node,omitempty"`
}

type InstructionCall struct {
	Name      string                    `json:"name"`
	IsUser    bool                      `json:"is_user"`
	Arguments []InstructionCallArgument `json:"arguments"`
}

type InstructionCallArgument struct {
	Value     Node       `json:"value"`
	Consteval bool       `json:"consteval"`
	Kind      token.Kind `json:"kind"`
}

type BinaryExpression struct {
	Operator string `json:"operator"`
	Children []Node `json:"children"`
}

type ModuleDeclaration struct {
	Name string `json:"name"`
	Body []Node `json:"body"`
}

type ArrayValue struct {
	MaxSize  int            `json:"max_size"`
	Elements []ArrayElement `json:"elements"`
}

type ArrayElement struct {
	Value Node `json:"value"`
}

func ToString(n Node) string {
	if n == nil {
		return "<nil>"
	}
	json, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(json)
}

func (Declaration) Node()             {}
func (FunctionParameter) Node()       {}
func (FunctionValue) Node()           {}
func (Value) Node()                   {}
func (InstructionCall) Node()         {}
func (InstructionCallArgument) Node() {}
func (BinaryExpression) Node()        {}
func (ModuleDeclaration) Node()       {}
func (ArrayValue) Node()              {}
func (ArrayElement) Node()            {}
