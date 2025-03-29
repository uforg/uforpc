package ast

// ProcDecl represents a procedure declaration in the URPC schema.
type ProcDecl struct {
	Doc      string
	Name     string
	Input    ProcInput
	Output   ProcOutput
	Metadata ProcMeta
}

func (p *ProcDecl) NodeType() NodeType { return NodeTypeProcDecl }

// ProcInput represents the input of a procedure.
type ProcInput struct {
	Fields []Field
}

func (i *ProcInput) NodeType() NodeType { return NodeTypeInput }

// ProcOutput represents the output of a procedure.
type ProcOutput struct {
	Fields []Field
}

func (o *ProcOutput) NodeType() NodeType { return NodeTypeOutput }

// ProcMeta represents the metadata of a procedure.
type ProcMeta struct {
	Entries []ProcMetaKV
}

func (m *ProcMeta) NodeType() NodeType { return NodeTypeMetadata }

type ProcMetaKV struct {
	Type  PrimitiveType // Only primitive types are allowed
	Key   string
	Value string
}
