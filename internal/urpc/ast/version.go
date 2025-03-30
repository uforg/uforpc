package ast

// Version represents the version of the URPC schema.
type Version struct {
	Pos   Position
	IsSet bool
	Value int
}

func (v *Version) NodeType() NodeType    { return NodeTypeVersion }
func (v *Version) GetPosition() Position { return v.Pos }
