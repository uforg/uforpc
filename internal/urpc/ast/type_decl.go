package ast

// TypeDecl represents a type declaration in the URPC schema.
type TypeDecl struct {
	Doc    string
	Name   string
	Fields []Field
}

func (t *TypeDecl) NodeType() NodeType { return NodeTypeTypeDecl }
