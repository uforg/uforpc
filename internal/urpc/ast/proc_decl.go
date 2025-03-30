package ast

// ProcDecl represents a procedure declaration in the URPC schema.
type ProcDecl struct {
	Pos      Position
	Doc      string
	Name     string
	Input    ProcInput
	Output   ProcOutput
	Metadata ProcMeta
}

func (p *ProcDecl) NodeType() NodeType    { return NodeTypeProcDecl }
func (p *ProcDecl) GetPosition() Position { return p.Pos }

// ProcInput represents the input of a procedure.
type ProcInput struct {
	Pos    Position
	Fields []Field
}

func (i *ProcInput) NodeType() NodeType    { return NodeTypeInput }
func (i *ProcInput) GetPosition() Position { return i.Pos }

// ProcOutput represents the output of a procedure.
type ProcOutput struct {
	Pos    Position
	Fields []Field
}

func (o *ProcOutput) NodeType() NodeType    { return NodeTypeOutput }
func (o *ProcOutput) GetPosition() Position { return o.Pos }

// ProcMeta represents the metadata of a procedure.
type ProcMeta struct {
	Pos     Position
	Entries []ProcMetaKV
}

func (m *ProcMeta) NodeType() NodeType    { return NodeTypeMetadata }
func (m *ProcMeta) GetPosition() Position { return m.Pos }

type ProcMetaKV struct {
	Pos   Position
	Type  PrimitiveType // Only primitive types are allowed
	Key   string
	Value string
}
