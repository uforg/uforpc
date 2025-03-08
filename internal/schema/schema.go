// Package schema provides the data structures for UFO RPC schema definition
package schema

import (
	"github.com/orsinium-labs/enum"
)

type (
	// PrimitiveType represents the allowed primitive data types
	PrimitiveType enum.Member[string]

	// ProcedureType represents the allowed procedure types
	ProcedureType enum.Member[string]
)

var (
	PrimitiveTypeString  = PrimitiveType{"string"}
	PrimitiveTypeInt     = PrimitiveType{"int"}
	PrimitiveTypeFloat   = PrimitiveType{"float"}
	PrimitiveTypeBoolean = PrimitiveType{"boolean"}

	ProcedureTypeQuery    = ProcedureType{"query"}
	ProcedureTypeMutation = ProcedureType{"mutation"}
)

// Schema represents the complete UFO RPC schema
type Schema struct {
	Types      map[string]Field `json:"types,omitzero"`
	Procedures []Procedure      `json:"procedures,omitzero"`
}

// Procedure represents an RPC procedure (query or mutation)
type Procedure struct {
	Name        string         `json:"name"`
	Type        ProcedureType  `json:"type"`
	Description string         `json:"description,omitzero"`
	Input       Field          `json:"input,omitzero"`
	Output      Field          `json:"output,omitzero"`
	Meta        map[string]any `json:"meta,omitzero"`
}
