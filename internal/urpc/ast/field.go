package ast

// Field represents a field in a type declaration or procedure input/output.
type Field struct {
	Name            string
	Type            Type
	Optional        bool
	ValidationRules []ValidationRule
}

func (f *Field) NodeType() NodeType { return NodeTypeField }
