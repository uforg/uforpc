package ast

// Node is the interface that all AST nodes implement.
type Node interface {
	NodeType() NodeType
	GetPosition() Position
}

// Schema is the root node of the AST representing an entire URPC schema.
type Schema struct {
	Pos         Position
	Version     Version
	CustomRules []CustomRuleDecl
	Types       []TypeDecl
	Procedures  []ProcDecl
}

func (s *Schema) NodeType() NodeType    { return NodeTypeSchema }
func (s *Schema) GetPosition() Position { return s.Pos }
