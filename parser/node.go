package parser

type Node interface {
	Node()
}

type Program struct {
	Statements []Node
}

type FunctionDeclaration struct {
	Name        string
	ParamsTypes []string
	ReturnType  string
	Body        []Node
	Exported    bool
}

type InstructionCall struct {
	Name      string
	Arguments []Node
}

type VariableDeclaration struct {
	Name, Type, Value string
	Exported          bool
}

type ReturnStatement struct {
	Value string
}

func (Program) Node()             {}
func (FunctionDeclaration) Node() {}
func (InstructionCall) Node()     {}
func (VariableDeclaration) Node() {}
func (ReturnStatement) Node()     {}
