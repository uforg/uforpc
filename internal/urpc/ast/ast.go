package ast

// Node is the interface that all AST nodes implement.
type Node interface {
	NodeType() NodeType
}

// Schema is the root node of the AST representing an entire URPC schema.
type Schema struct {
	Version     Version
	CustomRules []CustomRuleDecl
	Types       []TypeDecl
	Procedures  []ProcDecl
}

func (s *Schema) NodeType() NodeType { return NodeTypeSchema }
