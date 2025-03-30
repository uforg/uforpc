package ast

// NodeType represents the type of a node in the URPC schema.
type NodeType string

const (
	NodeTypeSchema         NodeType = "schema"
	NodeTypeVersion        NodeType = "version"
	NodeTypeCustomRuleDecl NodeType = "rule_decl"
	NodeTypeTypeDecl       NodeType = "type_decl"
	NodeTypeProcDecl       NodeType = "proc_decl"
	NodeTypeField          NodeType = "field"
	NodeTypeValidationRule NodeType = "validation_rule"
	NodeTypeInput          NodeType = "input"
	NodeTypeOutput         NodeType = "output"
	NodeTypeMetadata       NodeType = "metadata"
	NodeTypeType           NodeType = "type"
	NodeTypeTypeReference  NodeType = "type_reference"
)

// PrimitiveType represents the primitive types in the URPC schema
// excluding composite types (arrays, and objects).
type PrimitiveType string

const (
	PrimitiveTypeString  PrimitiveType = "string"
	PrimitiveTypeInt     PrimitiveType = "int"
	PrimitiveTypeFloat   PrimitiveType = "float"
	PrimitiveTypeBoolean PrimitiveType = "boolean"
)

// Position represents the position range of a node in the source file.
// This information is used for error reporting and LSP integration.
type Position struct {
	Filename  string
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
}
