// Package schema provides the data structures for UFO RPC schema definition
package schema

import (
	"github.com/orsinium-labs/enum"
)

// PrimitiveType represents the allowed primitive data types
type PrimitiveType enum.Member[string]

func (primitiveType PrimitiveType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + primitiveType.Value + `"`), nil
}

func (primitiveType *PrimitiveType) UnmarshalJSON(data []byte) error {
	value := string(data[1 : len(data)-1])
	*primitiveType = PrimitiveType{value}
	return nil
}

// Allowed primitive types
var (
	PrimitiveTypeString  = PrimitiveType{"string"}
	PrimitiveTypeInt     = PrimitiveType{"int"}
	PrimitiveTypeFloat   = PrimitiveType{"float"}
	PrimitiveTypeBoolean = PrimitiveType{"boolean"}
)

// ProcedureType represents the allowed procedure types
type ProcedureType enum.Member[string]

// MarshalJSON implements the json.Marshaler interface
func (procedureType ProcedureType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + procedureType.Value + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (procedureType *ProcedureType) UnmarshalJSON(data []byte) error {
	value := string(data[1 : len(data)-1])
	*procedureType = ProcedureType{value}
	return nil
}

// Allowed procedure types
var (
	ProcedureTypeQuery    = ProcedureType{"query"}
	ProcedureTypeMutation = ProcedureType{"mutation"}
)

// Schema represents the complete UFO RPC schema
type Schema struct {
	Types      map[string]Field     `json:"types,omitzero"`
	Procedures map[string]Procedure `json:"procedures,omitzero"`
}

// Procedure represents an RPC procedure (query or mutation)
type Procedure struct {
	Type        ProcedureType  `json:"type"`
	Description string         `json:"description,omitzero"`
	Input       Field          `json:"input,omitzero"`
	Output      Field          `json:"output,omitzero"`
	Meta        map[string]any `json:"meta,omitzero"`
}
