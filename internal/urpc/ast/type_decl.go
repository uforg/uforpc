package ast

// TypeDecl represents a type declaration in the URPC schema.
type TypeDecl struct {
	Pos     Position
	Doc     string
	Name    string
	Extends []TypeDecl
	Fields  []Field
}

func (t *TypeDecl) NodeType() NodeType    { return NodeTypeTypeDecl }
func (t *TypeDecl) GetPosition() Position { return t.Pos }
