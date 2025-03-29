package ast

// Version represents the version of the URPC schema.
type Version struct {
	IsSet bool
	Value int
}

func (v *Version) NodeType() NodeType { return NodeTypeVersion }
